package browser

import (
	"net/url"
	"testing"
)

func TestScopeRejectsExternalAndOtherRun(t *testing.T) {
	u, _ := url.Parse("http://127.0.0.1:8080/sandbox/run-one")
	b := Browser{allowed: u}
	for _, raw := range []string{"https://shop.example/", "http://127.0.0.1:8080/api/runs", "http://127.0.0.1:8080/sandbox/run-two", "http://127.0.0.1:8080/sandbox/run-one-evil", "http://user@127.0.0.1:8080/sandbox/run-one", "http://127.0.0.1:8081/sandbox/run-one"} {
		if b.Allowed(raw) {
			t.Errorf("scope accepted %s", raw)
		}
	}
	if !b.Allowed("http://127.0.0.1:8080/sandbox/run-one/quote?quantity=2") {
		t.Fatal("own sandbox resource rejected")
	}
}
