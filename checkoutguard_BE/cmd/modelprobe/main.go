// modelprobe measures one real provider request without browser or order actions.
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"

	"checkoutguard/internal/adapters/model"
	"checkoutguard/internal/config"
	"checkoutguard/internal/domain"
)

func main() {
	attempts := flag.Int("attempts", 1, "number of diagnostic requests (1..3)")
	http1 := flag.Bool("http1", false, "diagnose with HTTP/1 only")
	catalog := flag.Bool("catalog", false, "read model IDs instead of making a completion")
	modelName := flag.String("model", "", "diagnostic model override; does not change server configuration")
	messagesFile := flag.String("messages-file", "", "private diagnostic messages JSON; only runs inference, never executes returned actions")
	envFile := flag.String("env-file", "", "provider configuration file; defaults to MODEL_ENV_FILE or the existing Public AI configuration")
	pricingFile := flag.String("pricing-file", "", "optional model pricing file; does not change server configuration")
	flag.Parse()
	if *attempts < 1 || *attempts > 3 {
		fmt.Fprintln(os.Stderr, "attempts must be 1..3")
		os.Exit(1)
	}
	configuration := *envFile
	if configuration == "" {
		defaultFile := "../../.apertus-api.env"
		if _, err := os.Stat("../.env.model.local"); err == nil {
			defaultFile = "../.env.model.local"
		}
		configuration = config.Env("MODEL_ENV_FILE", defaultFile)
	}
	prices := *pricingFile
	if prices == "" {
		prices = os.Getenv("MODEL_PRICING_ENV_FILE")
	}
	if prices != "" {
		if err := config.LoadEnv(prices); err != nil {
			fmt.Fprintln(os.Stderr, "cannot load pricing configuration")
			os.Exit(1)
		}
	}
	if err := config.LoadEnv(configuration); err != nil {
		fmt.Fprintln(os.Stderr, "cannot load model configuration")
		os.Exit(1)
	}
	planner := model.NewModel()
	if *modelName != "" {
		planner.Name = *modelName
	}
	if !planner.Ready() {
		fmt.Fprintln(os.Stderr, "model configuration is incomplete; set endpoint, key and model for the same provider")
		os.Exit(1)
	}
	messages := []domain.Message{{Role: "system", Content: `Return exactly one JSON object: {"tool":"observe","reason":"Inspect checkout"}. No other keys or text.`}, {Role: "user", Content: "Inspect the synthetic checkout."}}
	if *messagesFile != "" {
		file, err := os.Open(*messagesFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot read diagnostic messages")
			os.Exit(1)
		}
		decoder := json.NewDecoder(io.LimitReader(file, 1<<20))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&messages)
		file.Close()
		if err != nil || len(messages) == 0 {
			fmt.Fprintln(os.Stderr, "invalid diagnostic messages")
			os.Exit(1)
		}
	}
	if *http1 {
		protocols := new(http.Protocols)
		protocols.SetHTTP1(true)
		planner.Client.Transport.(*http.Transport).Protocols = protocols
	}
	if *catalog {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		r, err := http.NewRequestWithContext(ctx, "GET", planner.Base+"/models", nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid provider endpoint")
			os.Exit(1)
		}
		r.Header.Set("Authorization", "Bearer "+planner.Key)
		r.Header.Set("User-Agent", "CheckoutGuard/0.1")
		response, err := planner.Client.Do(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, "model catalog request failed")
			os.Exit(1)
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "catalog HTTP %d\n", response.StatusCode)
			os.Exit(1)
		}
		var result struct {
			Data []struct {
				ID      string `json:"id"`
				Context *int   `json:"context_length,omitempty"`
			} `json:"data"`
		}
		if err = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil {
			fmt.Fprintln(os.Stderr, "invalid catalog response")
			os.Exit(1)
		}
		found := false
		for _, m := range result.Data {
			if m.ID == planner.Name {
				found = true
			}
		}
		json.NewEncoder(os.Stdout).Encode(map[string]any{"requested_model_available": found, "models": result.Data, "protocol": response.Proto})
		return
	}
	succeeded := false
	for attempt := 1; attempt <= *attempts; attempt++ {
		started := time.Now()
		var mu sync.Mutex
		phase := func(name string, ok bool) {
			mu.Lock()
			defer mu.Unlock()
			fmt.Printf("attempt=%d phase=%s elapsed_ms=%d ok=%t\n", attempt, name, time.Since(started).Milliseconds(), ok)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
			DNSDone:          func(info httptrace.DNSDoneInfo) { phase("dns", info.Err == nil) },
			ConnectDone:      func(_, _ string, err error) { phase("connect", err == nil) },
			TLSHandshakeDone: func(_ tls.ConnectionState, err error) { phase("tls", err == nil) },
			GotConn: func(info httptrace.GotConnInfo) {
				if info.Reused {
					phase("connection_reused", true)
				} else {
					phase("connection_ready", true)
				}
			},
			WroteRequest:         func(info httptrace.WroteRequestInfo) { phase("request_written", info.Err == nil) },
			GotFirstResponseByte: func() { phase("first_response_byte", true) },
		})
		action, usage, err := planner.Next(ctx, messages)
		cancel()
		result := map[string]any{"attempt": attempt, "elapsed_ms": time.Since(started).Milliseconds(), "provider": planner.Provider, "requested_model": planner.Name, "usage": usage, "diagnostic_only": true}
		if err != nil {
			result["error"] = err.Error()
		} else {
			result["action"] = action
			succeeded = true
		}
		json.NewEncoder(os.Stdout).Encode(result)
	}
	if !succeeded {
		os.Exit(1)
	}
}
