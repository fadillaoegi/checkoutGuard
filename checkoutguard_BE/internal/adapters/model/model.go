package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"checkoutguard/internal/config"
	"checkoutguard/internal/domain"
	neatlogs "github.com/neatlogs/neatlogs-go"
)

type Model struct {
	Base, Key, Name, Provider string
	Client                    *http.Client
	InputRate, OutputRate     *float64
}

func NewModel() *Model {
	rate := func(k string) *float64 {
		v, err := strconv.ParseFloat(os.Getenv(k), 64)
		if err != nil || v < 0 || math.IsNaN(v) || math.IsInf(v, 0) {
			return nil
		}
		return &v
	}
	base, key, name := os.Getenv("APERTUS_API_BASE_URL"), os.Getenv("APERTUS_API_KEY"), os.Getenv("APERTUS_MODEL")
	// A provider profile is one credential set. Never fill a missing candidate
	// key/model/endpoint using credentials from the legacy Public AI profile.
	for _, variable := range []string{"MODEL_BASE_URL", "MODEL_API_KEY", "MODEL_NAME"} {
		if _, explicit := os.LookupEnv(variable); explicit {
			base, key, name = os.Getenv("MODEL_BASE_URL"), os.Getenv("MODEL_API_KEY"), os.Getenv("MODEL_NAME")
			break
		}
	}
	provider := "openai-compatible"
	if endpoint, err := url.Parse(base); err == nil && endpoint.Hostname() == "api.publicai.co" {
		provider = "publicai"
	}
	return &Model{Base: strings.TrimRight(base, "/"), Key: key, Name: name, Provider: config.Env("MODEL_PROVIDER", provider), Client: newClient(), InputRate: rate("MODEL_INPUT_USD_PER_MILLION"), OutputRate: rate("MODEL_OUTPUT_USD_PER_MILLION")}
}
func (m *Model) Ready() bool { return m.Key != "" && m.Base != "" && m.Name != "" }
func (m *Model) Next(ctx context.Context, messages []Message) (action Action, usage Usage, err error) {
	nlMessages := make([]neatlogs.LLMMessage, len(messages))
	for i, msg := range messages {
		nlMessages[i] = neatlogs.LLMMessage{Role: msg.Role, Content: msg.Content}
	}
	ctx, span := neatlogs.StartLLMSpan(ctx, neatlogs.LLMCallOptions{Provider: m.Provider, Model: m.Name, Messages: nlMessages, MaxTokens: 450})
	defer func() {
		if err != nil {
			span.SetError(err)
			if usage.Input == 0 && usage.Output == 0 {
				usage.CostIncomplete = true
			}
		}
		span.End()
	}()
	if !m.Ready() {
		return action, usage, &domain.PlannerError{Kind: "configuration", Message: "model configuration is incomplete; set endpoint, key and model for the same provider"}
	}
	body, _ := json.Marshal(map[string]any{"model": m.Name, "messages": messages, "max_tokens": 450, "temperature": 0.1, "response_format": map[string]string{"type": "json_object"}})
	request, err := http.NewRequestWithContext(ctx, "POST", m.Base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return action, usage, &domain.PlannerError{Kind: "configuration", Message: "invalid model endpoint"}
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "CheckoutGuard/0.1")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+m.Key)
	response, err := m.Client.Do(request)
	if err != nil {
		m.Client.CloseIdleConnections()
		return action, usage, requestError(ctx, err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return action, usage, statusError(response.StatusCode, response.Header)
	}
	var result struct {
		Choices []struct {
			Message Message `json:"message"`
			Finish  string  `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			Input  int `json:"prompt_tokens"`
			Output int `json:"completion_tokens"`
			Total  int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
		ID    string `json:"id"`
	}
	if err = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil {
		var network net.Error
		if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &network) && network.Timeout()) {
			return action, usage, requestError(ctx, err)
		}
		return action, usage, invalidResponse("invalid model response", true)
	}
	usage.Input = result.Usage.Input
	usage.Output = result.Usage.Output
	usage.Model = result.Model
	pricedModel := config.Env("MODEL_PRICED_MODEL", m.Name)
	pricedProvider := config.Env("MODEL_PRICED_PROVIDER", m.Provider)
	if m.InputRate != nil && m.OutputRate != nil && usage.Model == pricedModel && m.Provider == pricedProvider {
		v := (float64(usage.Input)*(*m.InputRate) + float64(usage.Output)*(*m.OutputRate)) / 1e6
		usage.CostUSD = &v
	}
	usage.CostIncomplete = usage.CostUSD == nil
	span.SetUsage(usage.Input, usage.Output, result.Usage.Total)
	span.SetModel(result.Model)
	span.SetResponseID(result.ID)
	if len(result.Choices) != 1 {
		return action, usage, invalidResponse("model must return one action", true)
	}
	content := strings.TrimSpace(result.Choices[0].Message.Content)
	span.SetOutputMessage("assistant", content)
	span.SetFinishReason(result.Choices[0].Finish)
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&action); err != nil {
		return action, usage, invalidResponse("action must match JSON schema", true)
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return action, usage, invalidResponse("action must contain one JSON object", true)
	}
	if err = domain.ValidateAction(action); err != nil {
		err = invalidResponse(err.Error(), true)
	}
	return action, usage, err
}

func (m *Model) Info() ModelInfo {
	return ModelInfo{Name: m.Name, Provider: m.Provider, Configured: m.Ready()}
}
