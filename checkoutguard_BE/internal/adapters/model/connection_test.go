package model

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"checkoutguard/internal/domain"
)

func TestProviderHeadersAndFailureClassification(t *testing.T) {
	for _, c := range []struct {
		status int
		kind   string
		retry  bool
	}{{401, "authentication", false}, {403, "authentication", false}, {400, "request_rejected", false}, {429, "rate_limited", true}, {503, "unavailable", true}} {
		t.Run(http.StatusText(c.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("User-Agent") != "CheckoutGuard/0.1" || r.Header.Get("Accept") != "application/json" {
					t.Error("application headers absent")
				}
				w.Header().Set("Retry-After", "3")
				w.WriteHeader(c.status)
				w.Write([]byte("private provider details must not escape"))
			}))
			defer server.Close()
			m := Model{Base: server.URL, Key: "private-test-key", Name: "test", Client: server.Client()}
			_, _, err := m.Next(context.Background(), nil)
			var failure *domain.PlannerError
			if !errors.As(err, &failure) || failure.Kind != c.kind || failure.Retryable != c.retry || failure.RetryAfter != 3*time.Second {
				t.Fatalf("failure classification: %+v", err)
			}
			if strings.Contains(err.Error(), "private") {
				t.Fatal("provider body or credentials escaped")
			}
		})
	}
}

func TestTransportTimeoutAndCancellationAreDifferent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		select {
		case <-r.Context().Done():
		case <-time.After(200 * time.Millisecond):
		}
	}))
	defer server.Close()
	m := Model{Base: server.URL, Key: "key", Name: "test", Client: &http.Client{Timeout: 25 * time.Millisecond}}
	_, _, err := m.Next(context.Background(), nil)
	var failure *domain.PlannerError
	if !errors.As(err, &failure) || failure.Kind != "timeout" || !failure.Retryable {
		t.Fatalf("timeout: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = m.Next(ctx, nil)
	if !errors.As(err, &failure) || failure.Kind != "canceled" || failure.Retryable {
		t.Fatalf("cancel: %v", err)
	}
}

func TestRetryAfterDateAndConfiguredTimeoutBounds(t *testing.T) {
	now := time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC)
	if got := retryAfter(now.Add(8*time.Second).Format(http.TimeFormat), now); got != 8*time.Second {
		t.Fatal(got)
	}
	for _, value := range []string{"-5", "oops", now.Add(-time.Second).Format(http.TimeFormat)} {
		if retryAfter(value, now) != 0 {
			t.Fatal("invalid retry delay accepted")
		}
	}
	for _, c := range []struct {
		value   string
		seconds int
	}{{"60", 60}, {"5", 45}, {"120", 45}, {"bad", 45}} {
		t.Setenv("MODEL_REQUEST_TIMEOUT_SECONDS", c.value)
		client := newClient()
		defer client.CloseIdleConnections()
		if client.Timeout != time.Duration(c.seconds)*time.Second {
			t.Fatal("timeout outside bounded configuration")
		}
	}
}

func TestTimeoutWhileReadingProviderBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		select {
		case <-r.Context().Done():
		case <-time.After(200 * time.Millisecond):
		}
	}))
	defer server.Close()
	m := Model{Base: server.URL, Key: "key", Name: "test", Client: &http.Client{Timeout: 25 * time.Millisecond}}
	_, usage, err := m.Next(context.Background(), nil)
	var failure *domain.PlannerError
	if !errors.As(err, &failure) || failure.Kind != "timeout" || !failure.Retryable || usage.CostUSD != nil || !usage.CostIncomplete {
		t.Fatalf("body timeout mislabeled: %v", err)
	}
}

func TestProviderUsageRetainedOnRetryableMalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"model": "returned", "choices": []any{map[string]any{"message": map[string]string{"content": "not JSON"}}}, "usage": map[string]int{"prompt_tokens": 80, "completion_tokens": 4, "total_tokens": 84}})
	}))
	defer server.Close()
	m := Model{Base: server.URL, Key: "key", Name: "requested", Client: server.Client()}
	_, usage, err := m.Next(context.Background(), nil)
	var failure *domain.PlannerError
	if !errors.As(err, &failure) || !failure.Retryable || usage.Input != 80 || usage.Output != 4 || usage.CostUSD != nil {
		t.Fatalf("actual usage on invalid action lost: %+v %v", usage, err)
	}
}
