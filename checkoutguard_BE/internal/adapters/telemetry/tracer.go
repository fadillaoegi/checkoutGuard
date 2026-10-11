package telemetry

import (
	"checkoutguard/internal/application"
	"context"
	neatlogs "github.com/neatlogs/neatlogs-go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{}
type span struct {
	trace.Span
	end func()
}

func (s span) ID() string       { return s.SpanContext().TraceID().String() }
func (s span) Input(v any)      { _ = neatlogs.SetTraceInput(s.Span, v) }
func (s span) Output(v any)     { _ = neatlogs.SetTraceOutput(s.Span, v) }
func (s span) Error(msg string) { s.SetStatus(codes.Error, msg) }
func (s span) End()             { s.end() }
func (Tracer) Workflow(ctx context.Context, name string) (context.Context, application.Span) {
	ctx, s, end := neatlogs.Trace(ctx, name)
	return ctx, span{Span: s, end: end}
}
func (Tracer) Tool(ctx context.Context, name string) (context.Context, application.Span) {
	ctx, s, end := neatlogs.StartSpan(ctx, name, "tool")
	return ctx, span{Span: s, end: end}
}
