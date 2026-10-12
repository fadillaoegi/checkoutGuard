package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"checkoutguard/internal/domain"
)

// Each failed decision stays a persisted step. Recovery only repeats planning;
// browser/order actions are never replayed by this policy.
func modelRecovery(err error, consecutive, total int) (time.Duration, string) {
	var failure *domain.PlannerError
	if errors.As(err, &failure) && !failure.Retryable {
		return 0, "Model request cannot be retried: " + failure.Message
	}
	if consecutive >= 2 || total >= 3 {
		return 0, "Model recovery budget exhausted; provider or response remains unavailable"
	}
	delay := time.Second
	if failure != nil && failure.RetryAfter > delay {
		delay = failure.RetryAfter
	}
	if delay > 15*time.Second {
		return 0, "Provider requires a longer wait; start another run later"
	}
	return delay, ""
}

func waitModelRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func modelFailureResult(err error) map[string]any {
	result := map[string]any{"error": err.Error(), "browser_action_executed": false}
	var failure *domain.PlannerError
	if errors.As(err, &failure) {
		result["error_kind"] = failure.Kind
		result["retryable"] = failure.Retryable
		result["retry_after_ms"] = failure.RetryAfter.Milliseconds()
	}
	return result
}
func decisionMessages(messages []Message, step Step, observation Message) []Message {
	// Keep the policy, goal and complete action/result history. Superseded page
	// snapshots repeat large untrusted DOM text and stale actionable controls;
	// only the latest observation belongs in the next decision request. The
	// persisted screenshots and steps remain untouched.
	current := make([]Message, 0, len(messages)+3)
	for i, message := range messages {
		if i >= 2 && message.Role == "user" && (strings.HasPrefix(message.Content, observationPrefix) || strings.HasPrefix(message.Content, verificationStatePrefix)) {
			continue
		}
		current = append(current, message)
	}
	if step.Action.Tool == "model_error" {
		return append(current, Message{Role: "user", Content: "Model decision failed; no browser action was executed. " + resultJSON(step.Result)}, observation)
	}
	return append(current, Message{Role: "assistant", Content: actionJSON(step.Action)}, Message{Role: "user", Content: "Tool result: " + resultJSON(step.Result)}, observation)
}
