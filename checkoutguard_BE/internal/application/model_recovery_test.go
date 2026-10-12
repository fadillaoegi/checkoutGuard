package application

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"checkoutguard/internal/domain"
)

func TestModelRecoveryHonorsPermanentErrorsAndRetryAfter(t *testing.T) {
	for _, c := range []struct {
		err      error
		failures int
		delay    time.Duration
		stop     bool
	}{
		{&domain.PlannerError{Kind: "authentication", Message: "model HTTP 401"}, 1, 0, true},
		{&domain.PlannerError{Kind: "timeout", Retryable: true}, 1, time.Second, false},
		{&domain.PlannerError{Kind: "rate_limited", Retryable: true, RetryAfter: 4 * time.Second}, 1, 4 * time.Second, false},
		{&domain.PlannerError{Kind: "rate_limited", Retryable: true, RetryAfter: 90 * time.Second}, 1, 0, true},
		{&domain.PlannerError{Kind: "timeout", Retryable: true}, 2, 0, true},
	} {
		delay, reason := modelRecovery(c.err, c.failures, c.failures)
		if delay != c.delay || (reason != "") != c.stop {
			t.Fatalf("recovery %v: %v %s", c.err, delay, reason)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(waitModelRetry(ctx, 10*time.Second), context.Canceled) {
		t.Fatal("cancellation did not interrupt delay")
	}
	transient := &domain.PlannerError{Kind: "timeout", Retryable: true}
	if _, reason := modelRecovery(transient, 1, 2); reason != "" {
		t.Fatal("two separated failures exhausted consecutive budget")
	}
	if _, reason := modelRecovery(transient, 1, 3); reason == "" {
		t.Fatal("total error budget not enforced")
	}
}

func TestFailedDecisionHistoryContainsNoAssistantAction(t *testing.T) {
	step := Step{Action: Action{Tool: "model_error"}, Result: map[string]any{"error": "model response timed out"}}
	messages := decisionMessages([]Message{{Role: "system", Content: "policy"}, {Role: "user", Content: "goal"}}, step, Message{Role: "user", Content: "current observation"})
	if len(messages) != 4 || messages[2].Role != "user" || !strings.Contains(messages[2].Content, "no browser action") {
		t.Fatal("failed planning supplied a fabricated assistant action")
	}
}

func TestDecisionContextReplacesOldPageButPreservesCheckoutSafetyHistory(t *testing.T) {
	policy := Message{Role: "system", Content: systemPrompt}
	// A user's goal must never be classified as a removable page snapshot.
	goal := Message{Role: "user", Content: observationPrefix + "user supplied goal"}
	messages := []Message{policy, goal, observationMessage(Observation{Text: "Old product page", Controls: []Control{{Selector: "#add"}}})}
	steps := []Step{
		{Action: Action{Tool: "verify_total"}, Result: map[string]any{"expected_minor": 9800, "displayed_minor": 11050, "matches": false}},
		{Action: Action{Tool: "click", Selector: "#place"}, Result: map[string]any{"error": "payment timed out"}},
		{Action: Action{Tool: "model_error"}, Result: modelFailureResult(&domain.PlannerError{Kind: "timeout", Message: "model response timed out", Retryable: true})},
		{Action: Action{Tool: "order_status"}, Result: map[string]any{"count": 1, "status": "confirmed"}},
	}
	for i, step := range steps {
		latest := observationMessage(Observation{Text: "Current checkout page", Controls: []Control{{Selector: "#status"}}})
		messages = decisionMessages(messages, step, latest)
		if messages[0] != policy || messages[1] != goal || messages[len(messages)-1] != latest {
			t.Fatal("policy, user goal or latest observation changed")
		}
		observations := 0
		for _, message := range messages[2:] {
			if strings.HasPrefix(message.Content, observationPrefix) {
				observations++
			}
			if strings.Contains(message.Content, "Old product page") {
				t.Fatal("superseded product page retained")
			}
		}
		if observations != 1 {
			t.Fatal("exactly one current observation required")
		}
		for _, previous := range steps[:i+1] {
			found := false
			for _, message := range messages[2:] {
				if strings.Contains(message.Content, resultJSON(previous.Result)) {
					found = true
				}
			}
			if !found {
				t.Fatal("pricing, payment or confirmed-order result was dropped")
			}
		}
	}
	if len(messages) != 10 { // policy + goal + three action/result pairs + error + observation
		t.Fatal("unexpected context history length", len(messages))
	}
}

func TestPermanentProviderFailureStopsAfterOneDecisionAndPersistsEvidence(t *testing.T) {
	repository := newRepository()
	calls := 0
	fixture := &marketFixture{state: "product"}
	s := newService(repository, testPlanner{next: func(context.Context, []Message) (Action, Usage, error) {
		calls++
		return Action{}, Usage{}, &domain.PlannerError{Kind: "authentication", Message: "model HTTP 401"}
	}}, &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true}, browser: fixture})
	r, err := s.Start(context.Background(), StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/1/2", Goal: "Prepare"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s.wg.Wait()
	final, _ := s.Get(context.Background(), r.ID)
	s.Close()
	if calls != 1 || final.Status != "failed" || !strings.Contains(final.Reason, "401") || len(final.Steps) != 2 || final.Steps[1].Action.Tool != "model_error" {
		t.Fatalf("permanent failure repeated or not stored: %+v calls=%d", final, calls)
	}
}

func TestTransientProviderFailureRecoversWithoutReplayingBrowserActions(t *testing.T) {
	repository := newRepository()
	calls := 0
	fixture := &marketFixture{state: "review"}
	estimate := 0.000004
	s := newService(repository, testPlanner{next: func(context.Context, []Message) (Action, Usage, error) {
		calls++
		if calls == 1 {
			return Action{}, Usage{CostIncomplete: true}, &domain.PlannerError{Kind: "unavailable", Retryable: true, Message: "model HTTP 503"}
		}
		return Action{Tool: "finish", Reason: "Review visible"}, Usage{Input: 10, Output: 2, CostUSD: &estimate}, nil
	}}, &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true}, browser: fixture})
	r, err := s.Start(context.Background(), StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/1/2", Goal: "Prepare"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s.wg.Wait()
	final, _ := s.Get(context.Background(), r.ID)
	s.Close()
	if calls != 2 || final.Status != "ready_for_review" || final.CheckoutReview.OrderSubmitted || final.Usage.Input != 10 || final.Steps[1].Error != "model HTTP 503" {
		t.Fatalf("transient recovery lost history: %+v", final)
	}
	if final.Usage.CostUSD == nil || *final.Usage.CostUSD != estimate || !final.Usage.CostIncomplete {
		t.Fatal("unreported usage lost its partial-cost label after recovery")
	}
}
