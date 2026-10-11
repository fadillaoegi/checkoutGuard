package domain

import (
	"testing"
)

func TestActionSchemaRejectsUnsafeTools(t *testing.T) {
	for _, a := range []Action{{Tool: "evaluate", Value: "fetch('https://production')", Reason: "escape"}, {Tool: "finish", Selector: "#place", Reason: "unexpected arguments"}, {Tool: "fill", Selector: "#email", Reason: "missing value"}, {Tool: "click", Selector: "#place", Value: "anything", Reason: "extra value"}} {
		if ValidateAction(a) == nil {
			t.Errorf("accepted invalid action %+v", a)
		}
	}
}
func TestFinishStatusIncludesSubtotalDefect(t *testing.T) {
	r := Run{Expected: 12400, Displayed: 12400, Findings: []Finding{{Kind: "incorrect_total", Verified: true}}}
	if !r.HasDefect() {
		t.Fatal("matching final total concealed a breakdown defect")
	}
}
