package model

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	neatlogs "github.com/neatlogs/neatlogs-go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestModelTraceSeparatesRequestFromResponse(t *testing.T) {
	const requested = "requested-model"
	for _, tc := range []struct {
		name, returned, content, errorKind string
		status                             int
		timeout, bodyTimeout, canceled     bool
		received, decoded                  bool
	}{
		{name: "matching response", returned: requested, content: `{"tool":"observe","reason":"Inspect checkout"}`, status: 200, received: true, decoded: true, errorKind: "none"},
		{name: "different response label", returned: "returned-model", content: `{"tool":"observe","reason":"Inspect checkout"}`, status: 200, received: true, decoded: true, errorKind: "none"},
		{name: "missing response label", content: `{"tool":"observe","reason":"Inspect checkout"}`, status: 200, received: true, decoded: true, errorKind: "none"},
		{name: "timeout before headers", timeout: true, errorKind: "timeout"},
		{name: "timeout after headers", timeout: true, bodyTimeout: true, status: 200, received: true, errorKind: "timeout"},
		{name: "HTTP unavailable", status: 503, received: true, errorKind: "unavailable"},
		{name: "invalid action with known model", returned: "returned-model", content: "not JSON", status: 200, received: true, decoded: true, errorKind: "invalid_response"},
		{name: "canceled request", canceled: true, errorKind: "canceled"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body struct {
					Model string `json:"model"`
				}
				var requestFields map[string]json.RawMessage
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &requestFields)
				if _, exists := requestFields["checkoutguard_diagnostics"]; exists {
					t.Error("diagnostics sent to provider")
				}
				if err := json.Unmarshal(raw, &body); err != nil || body.Model != requested {
					t.Error("request identity changed or request could not be decoded")
				}
				if tc.timeout {
					if tc.bodyTimeout {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(200)
						w.(http.Flusher).Flush()
					}
					<-r.Context().Done()
					return
				}
				if tc.status != 200 {
					w.WriteHeader(tc.status)
					io.WriteString(w, "private provider response")
					return
				}
				json.NewEncoder(w).Encode(map[string]any{
					"model": tc.returned, "id": "synthetic-response",
					"choices": []any{map[string]any{
						"message": map[string]string{"role": "assistant", "content": tc.content}, "finish_reason": "stop",
					}},
					"usage": map[string]int{"prompt_tokens": 80, "completion_tokens": 4, "total_tokens": 84},
				})
			}))
			defer server.Close()
			sink := tracetest.NewInMemoryExporter()
			client, err := neatlogs.NewClient(context.Background(), neatlogs.Config{WorkflowName: "checkoutguard.trace-test"}, neatlogs.WithExporter(sink))
			if err != nil {
				t.Fatal(err)
			}
			defer client.Shutdown(context.Background())
			ctx, cancel := context.WithCancel(client.Context(context.Background()))
			defer cancel()
			if tc.canceled {
				cancel()
			}
			httpClient := server.Client()
			httpClient.Timeout = 300 * time.Millisecond
			m := Model{Base: server.URL, Key: "private-test-key", Name: requested, Provider: "test-provider", Client: httpClient}
			_, usage, callErr := m.Next(ctx, []Message{{Role: "user", Content: "synthetic checkout"}})
			if (callErr != nil) != (tc.errorKind != "none") {
				t.Fatalf("unexpected call outcome: %v", callErr)
			}
			if usage.Model != tc.returned {
				t.Fatal("returned usage identity changed")
			}
			if err := client.Flush(context.Background()); err != nil {
				t.Fatal(err)
			}
			var llm *tracetest.SpanStub
			roots := 0
			for _, span := range sink.GetSpans() {
				if span.Name == "checkoutguard.trace-test" {
					roots++
				}
				if span.Name == "checkoutguard.model.request" {
					if llm != nil {
						t.Fatal("model call produced duplicate LLM spans")
					}
					copy := span
					llm = &copy
				}
			}
			if llm == nil || roots != 1 || !llm.Parent.IsValid() {
				t.Fatal("provider span or auto-root hierarchy lost")
			}
			attrs := map[string]any{}
			for _, attr := range llm.Attributes {
				attrs[string(attr.Key)] = attr.Value.AsInterface()
			}
			for key, want := range map[string]any{
				"neatlogs.span.kind": "llm", "neatlogs.llm.model_name": requested,
				"neatlogs.llm.provider": "test-provider",
			} {
				if attrs[key] != want {
					t.Errorf("%s = %v, want %v", key, attrs[key], want)
				}
			}
			var metadata struct {
				CheckoutGuard map[string]any `json:"checkoutguard"`
			}
			if err := json.Unmarshal([]byte(attrs["neatlogs.metadata"].(string)), &metadata); err != nil {
				t.Fatal(err)
			}
			output := metadata.CheckoutGuard
			if usage.ModelCall == nil || usage.ModelCall.SpanID != llm.SpanContext.SpanID().String() || usage.ModelCall.RequestedModel != requested || usage.ModelCall.ErrorKind != tc.errorKind {
				t.Fatal("usage lost span-correlated diagnostics")
			}
			if output["requested_model"] != requested || output["timeout_ms"] != float64(300) || output["error_kind"] != tc.errorKind || output["http_response_received"] != tc.received || output["response_decoded"] != tc.decoded {
				t.Fatal("canonical metadata lost diagnostic state")
			}
			var content string
			if err := json.Unmarshal([]byte(attrs["neatlogs.output.value"].(string)), &content); err != nil || content != tc.content {
				t.Fatal("provider output changed")
			}
			if tc.returned == "" {
				if _, exists := output["returned_model"]; exists || output["returned_model_known"] != false {
					t.Fatal("missing response label was invented")
				}
				if _, exists := output["model_label_mismatch"]; exists {
					t.Fatal("unknown response treated as a label mismatch")
				}
			} else {
				if output["returned_model"] != tc.returned || output["model_label_mismatch"] != (tc.returned != requested) {
					t.Fatal("requested/returned model distinction lost")
				}
				if attrs["neatlogs.llm.token_count.prompt"] != int64(80) || attrs["neatlogs.llm.token_count.completion"] != int64(4) {
					t.Fatal("provider usage lost")
				}
			}
			if tc.errorKind != "none" && (llm.Status.Code != codes.Error || len(llm.Events) == 0) {
				t.Fatal("error status or exception event lost")
			}
			encoded, _ := json.Marshal(llm)
			for _, private := range []string{"private-test-key", "private provider response", server.URL} {
				if strings.Contains(string(encoded), private) {
					t.Fatal("private credential, endpoint or rejected response leaked")
				}
			}
		})
	}
}

func TestUnconfiguredModelTraceDoesNotNeedHTTPClient(t *testing.T) {
	m := Model{}
	if _, _, err := m.Next(context.Background(), nil); err == nil {
		t.Fatal("unconfigured model accepted")
	}
}
