package model

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderUsageAndMalformedAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("provider credential absent")
		}
		json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"role": "assistant", "content": `{"tool":"finish","reason":"done","unexpected":true}`}}}, "usage": map[string]int{"prompt_tokens": 123, "completion_tokens": 12, "total_tokens": 135}})
	}))
	defer server.Close()
	model := Model{Base: server.URL, Key: "test-key", Name: "test-only", Client: server.Client()}
	_, usage, err := model.Next(context.Background(), []Message{{Role: "user", Content: "synthetic test"}})
	if err == nil {
		t.Fatal("unknown action property was accepted")
	}
	if usage.Input != 123 || usage.Output != 12 {
		t.Fatal("actual provider usage lost on invalid action")
	}
	if usage.CostUSD != nil {
		t.Fatal("unknown pricing should not appear free")
	}
}

func TestReturnedModelAndCostAreAccurate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"model": "returned-model", "choices": []any{map[string]any{"message": map[string]string{"role": "assistant", "content": `{"tool":"observe","reason":"Inspect checkout"}`}}}, "usage": map[string]int{"prompt_tokens": 100, "completion_tokens": 10, "total_tokens": 110}})
	}))
	defer server.Close()
	in, out := 1.0, 2.0
	model := Model{Base: server.URL, Key: "test-key", Name: "requested-model", Client: server.Client(), InputRate: &in, OutputRate: &out}
	_, usage, err := model.Next(context.Background(), []Message{{Role: "user", Content: "synthetic"}})
	if err != nil {
		t.Fatal(err)
	}
	if usage.Model != "returned-model" {
		t.Fatal("returned model was lost")
	}
	if usage.CostUSD != nil {
		t.Fatal("rates for a different requested model were applied")
	}
	t.Setenv("MODEL_PRICED_MODEL", "returned-model")
	_, usage, err = model.Next(context.Background(), []Message{{Role: "user", Content: "synthetic"}})
	if err != nil || usage.CostUSD == nil || *usage.CostUSD != 0.00012 {
		t.Fatalf("usage calculation: %+v, %v", usage, err)
	}
}

func TestRatesMustMatchBothProviderAndReturnedModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"model": "returned-model", "choices": []any{map[string]any{"message": map[string]string{"content": `{"tool":"observe","reason":"Inspect checkout"}`}}}, "usage": map[string]int{"prompt_tokens": 100, "completion_tokens": 10, "total_tokens": 110}})
	}))
	defer server.Close()
	in, out := 0.25, 0.50
	m := Model{Base: server.URL, Key: "test-key", Name: "requested", Provider: "candidate", Client: server.Client(), InputRate: &in, OutputRate: &out}
	t.Setenv("MODEL_PRICED_MODEL", "returned-model")
	t.Setenv("MODEL_PRICED_PROVIDER", "publicai")
	_, usage, err := m.Next(context.Background(), nil)
	if err != nil || usage.CostUSD != nil || !usage.CostIncomplete {
		t.Fatal("another provider's rates were applied")
	}
	t.Setenv("MODEL_PRICED_PROVIDER", "candidate")
	_, usage, err = m.Next(context.Background(), nil)
	if err != nil || usage.CostUSD == nil || usage.CostIncomplete || *usage.CostUSD != 0.000030 {
		t.Fatal("matching rate estimate not recorded accurately")
	}
}
