package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"checkoutguard/internal/domain"
	neatlogs "github.com/neatlogs/neatlogs-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Use the SDK's public provider-span API to attach request/response diagnostics
// to the actual LLM span. The SDK deliberately keeps its active context private;
// trace.SpanFromContext would not retrieve the span opened by StartLLMSpan.
type modelCallTrace struct {
	span             trace.Span
	end              func()
	requestedModel   string
	returnedModel    string
	timeoutMS        int64
	responseReceived bool
	responseDecoded  bool
	httpStatus       int
	output           string
}

func (m *Model) startCallTrace(ctx context.Context, messages []Message) (context.Context, *modelCallTrace) {
	ctx, span, end := neatlogs.StartProviderSpan(ctx, "checkoutguard.model.request", "llm")
	call := &modelCallTrace{span: span, end: end, requestedModel: m.Name}
	if m.Client != nil {
		call.timeoutMS = m.Client.Timeout.Milliseconds()
	}
	span.SetAttributes(
		attribute.String("neatlogs.llm.provider", m.Provider),
		attribute.String("neatlogs.llm.system", m.Provider),
		attribute.String("neatlogs.llm.model_name", m.Name),
		attribute.Bool("neatlogs.llm.is_streaming", false),
		attribute.Int("neatlogs.llm.max_tokens", 450),
		attribute.Float64("neatlogs.llm.temperature", 0.1),
		attribute.String("neatlogs.llm.invocation_parameters", `{"max_tokens":450,"temperature":0.1}`),
	)
	for i, message := range messages {
		span.SetAttributes(
			attribute.String(fmt.Sprintf("neatlogs.llm.input_messages.%d.role", i), message.Role),
			attribute.String(fmt.Sprintf("neatlogs.llm.input_messages.%d.content", i), message.Content),
		)
	}
	_ = neatlogs.SetTraceInput(span, messages)
	return ctx, call
}

func (c *modelCallTrace) recordResponse(model, id string, usage Usage, total int) {
	c.responseDecoded = true
	c.returnedModel = model
	if id != "" {
		c.span.SetAttributes(attribute.String("neatlogs.llm.response_id", id))
	}
	// Preserve unknown usage instead of recording zero as a measured value.
	if usage.Input != 0 {
		c.span.SetAttributes(attribute.Int("neatlogs.llm.token_count.prompt", usage.Input))
	}
	if usage.Output != 0 {
		c.span.SetAttributes(attribute.Int("neatlogs.llm.token_count.completion", usage.Output))
	}
	if total != 0 {
		c.span.SetAttributes(attribute.Int("neatlogs.llm.token_count.total", total))
	}
}

func (c *modelCallTrace) recordOutput(content, finish string) {
	c.output = content
	if content != "" {
		c.span.SetAttributes(
			attribute.String("neatlogs.llm.output_messages.0.role", "assistant"),
			attribute.String("neatlogs.llm.output_messages.0.content", content),
		)
	}
	if finish != "" {
		c.span.SetAttributes(attribute.String("neatlogs.llm.finish_reason", finish))
	}
}

func (c *modelCallTrace) finish(err error) domain.ModelCall {
	kind := "none"
	retryable := false
	if err != nil {
		kind = "unknown"
		var failure *domain.PlannerError
		if errors.As(err, &failure) {
			kind, retryable = failure.Kind, failure.Retryable
		}
		c.span.RecordError(err)
		c.span.SetStatus(codes.Error, err.Error())
	}
	diagnostics := domain.ModelCall{
		SpanID:         c.span.SpanContext().SpanID().String(),
		RequestedModel: c.requestedModel, TimeoutMS: c.timeoutMS,
		ErrorKind: kind, Retryable: retryable,
		HTTPResponseReceived: c.responseReceived, ResponseDecoded: c.responseDecoded,
		ReturnedModelKnown: c.returnedModel != "",
	}
	if c.responseReceived {
		diagnostics.HTTPStatusCode = c.httpStatus
	}
	if c.returnedModel != "" {
		mismatch := c.returnedModel != c.requestedModel
		diagnostics.ReturnedModel = c.returnedModel
		diagnostics.ModelLabelMismatch = &mismatch
	}
	// Keep SDK metadata; the workflow output also carries these observations
	// correlated by span ID, since hosted v3 omits custom LLM metadata.
	metadata, _ := json.Marshal(map[string]any{"checkoutguard": diagnostics})
	c.span.SetAttributes(attribute.String("neatlogs.metadata", string(metadata)))
	// Keep provider output intact. SetTraceOutput is reserved for workflow roots.
	encoded, _ := json.Marshal(c.output)
	c.span.SetAttributes(attribute.String("neatlogs.output.value", string(encoded)))
	c.end()
	return diagnostics
}
