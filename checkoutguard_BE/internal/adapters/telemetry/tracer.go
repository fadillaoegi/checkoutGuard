package telemetry

import (
	"checkoutguard/internal/application"
	"context"
	"encoding/json"
	neatlogs "github.com/neatlogs/neatlogs-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{}
type span struct {
	trace.Span
	end      func()
	workflow bool
}

func (s span) ID() string { return s.SpanContext().TraceID().String() }
func (s span) Input(v any) {
	_ = neatlogs.SetTraceInput(s.Span, v)
	s.setKindIO("input", v)
}
func (s span) Output(v any) {
	if s.workflow {
		_ = neatlogs.SetTraceOutput(s.Span, v)
	}
	s.setKindIO("output", v)
}

func (s span) setKindIO(direction string, v any) {
	if encoded, err := json.Marshal(v); err == nil {
		kind := "tool"
		if s.workflow {
			kind = "workflow"
		}
		s.SetAttributes(attribute.String("neatlogs."+direction+".value", string(encoded)),
			attribute.String("neatlogs."+kind+"."+direction, string(encoded)))
	}
}
func (s span) Error(msg string) { s.SetStatus(codes.Error, msg) }
func (s span) End()             { s.end() }
func (Tracer) Workflow(ctx context.Context, name string) (context.Context, application.Span) {
	ctx, s, end := neatlogs.Trace(ctx, name)
	return ctx, span{Span: s, end: end, workflow: true}
}
func (Tracer) Tool(ctx context.Context, name string) (context.Context, application.Span) {
	ctx, s, end := neatlogs.StartSpan(ctx, name, "tool")
	return ctx, span{Span: s, end: end}
}
