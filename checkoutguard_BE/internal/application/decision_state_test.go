package application

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPlannerStateRequiresPostSubmitOrderVerificationAndReplacesStaleState(t *testing.T) {
	run := Run{TotalVerified: true, OrderCount: 1}
	for _, c := range []struct{ placed, checked, allowed bool }{
		{false, true, false}, // A pre-submit order check never permits finish.
		{true, false, false}, // Page confirmation is not independent verification.
		{true, true, true},
	} {
		message := sandboxDecisionState(run, c.placed, c.checked, Observation{})
		var state struct {
			Allowed bool `json:"finish_allowed"`
			Count   *int `json:"confirmed_order_count_after_submit"`
		}
		if err := json.Unmarshal([]byte(strings.TrimPrefix(message.Content, verificationStatePrefix)), &state); err != nil {
			t.Fatal(err)
		}
		if state.Allowed != c.allowed || (state.Count != nil) != c.allowed {
			t.Fatal("planner was given unverified order state")
		}
		messages := []Message{{Role: "system", Content: systemPrompt}, {Role: "user", Content: "goal"}, message}
		updated := decisionMessages(messages, Step{Action: Action{Tool: "order_status"}, Result: map[string]any{"count": 1}}, observationMessage(Observation{Text: "latest page"}))
		for _, previous := range updated[2:] {
			if strings.HasPrefix(previous.Content, verificationStatePrefix) {
				t.Fatal("stale finish eligibility retained")
			}
		}
	}
	run.OrderCount = 2
	if strings.Contains(sandboxDecisionState(run, true, true, Observation{}).Content, `"finish_allowed":true`) {
		t.Fatal("duplicate order state permitted finish")
	}
	product := sandboxDecisionState(Run{}, false, false, Observation{Controls: []Control{{Type: "button", Selector: "#add"}, {Type: "input", Selector: "#quantity"}}})
	if strings.Contains(product.Content, "verify_total") || strings.Contains(product.Content, "order_status\"") {
		t.Fatal("planner was offered verification tools before checkout or submission")
	}
}
