package model

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"checkoutguard/internal/domain"
)

func clearProviderEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{"MODEL_BASE_URL", "MODEL_API_KEY", "MODEL_NAME", "MODEL_PROVIDER", "APERTUS_API_BASE_URL", "APERTUS_API_KEY", "APERTUS_MODEL"} {
		t.Setenv(name, "")
		os.Unsetenv(name)
	}
}

func TestCandidateProfileNeverInheritsPrimaryCredentials(t *testing.T) {
	clearProviderEnvironment(t)
	t.Setenv("APERTUS_API_BASE_URL", "https://api.publicai.co/v1")
	t.Setenv("APERTUS_API_KEY", "private-primary-test-key")
	t.Setenv("APERTUS_MODEL", "primary-model")
	t.Setenv("MODEL_BASE_URL", "https://candidate.example/v1")
	t.Setenv("MODEL_NAME", "candidate-model")
	planner := NewModel()
	defer planner.Client.CloseIdleConnections()
	if planner.Ready() || planner.Key != "" || planner.Provider == "publicai" {
		t.Fatal("incomplete candidate inherited primary credentials or provider")
	}
	called := false
	planner.Client = &http.Client{Transport: requestTransport(func(*http.Request) (*http.Response, error) {
		called = true
		return nil, errors.New("unexpected request")
	})}
	_, _, err := planner.Next(context.Background(), nil)
	var failure *domain.PlannerError
	if !errors.As(err, &failure) || failure.Kind != "configuration" || failure.Retryable || called {
		t.Fatal("incomplete configuration sent a request")
	}
}

type requestTransport func(*http.Request) (*http.Response, error)

func (r requestTransport) RoundTrip(request *http.Request) (*http.Response, error) { return r(request) }

func TestProviderProfileIsSelectedAsOneCompleteSet(t *testing.T) {
	clearProviderEnvironment(t)
	t.Setenv("APERTUS_API_BASE_URL", "https://api.publicai.co/v1/")
	t.Setenv("APERTUS_API_KEY", "private-primary-test-key")
	t.Setenv("APERTUS_MODEL", "primary-model")
	legacy := NewModel()
	defer legacy.Client.CloseIdleConnections()
	if !legacy.Ready() || legacy.Base != "https://api.publicai.co/v1" || legacy.Provider != "publicai" {
		t.Fatal("existing Public AI profile no longer works")
	}
	t.Setenv("MODEL_BASE_URL", "https://candidate.example/v1/")
	t.Setenv("MODEL_API_KEY", "private-candidate-test-key")
	t.Setenv("MODEL_NAME", "candidate-model")
	t.Setenv("MODEL_PROVIDER", "litellm")
	candidate := NewModel()
	defer candidate.Client.CloseIdleConnections()
	if !candidate.Ready() || candidate.Base != "https://candidate.example/v1" || candidate.Key != "private-candidate-test-key" || candidate.Name != "candidate-model" || candidate.Provider != "litellm" {
		t.Fatal("candidate profile fields mixed with primary")
	}
}
