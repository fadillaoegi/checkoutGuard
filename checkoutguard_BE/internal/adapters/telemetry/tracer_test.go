package telemetry

import (
	"context"
	"encoding/json"
	"testing"

	neatlogs "github.com/neatlogs/neatlogs-go"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestWorkflowAndToolOutputsKeepTheirScope(t *testing.T) {
	ctx := context.Background()
	sink := tracetest.NewInMemoryExporter()
	client, err := neatlogs.NewClient(ctx, neatlogs.Config{WorkflowName: "workflow"}, neatlogs.WithExporter(sink))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Shutdown(ctx)
	ctx, root := (Tracer{}).Workflow(client.Context(ctx), "workflow")
	_, tool := (Tracer{}).Tool(ctx, "tool")
	tool.Input(map[string]any{"tool": "order_status"})
	tool.Output(map[string]any{"order_count": 1})
	tool.End()
	root.Input(map[string]any{"synthetic": true})
	root.Output(map[string]any{"status": "passed", "model_calls": []string{"synthetic-diagnostic"}})
	root.End()
	if err := client.Flush(ctx); err != nil {
		t.Fatal(err)
	}
	// The SDK emits a transport completion marker after the workflow ends.
	spans := []tracetest.SpanStub{}
	for _, exported := range sink.GetSpans() {
		if exported.Name != "neatlogs.trace.complete" {
			spans = append(spans, exported)
		}
	}
	if len(spans) != 2 {
		t.Fatalf("unexpected application span count: %d", len(spans))
	}
	for _, exported := range spans {
		attrs := map[string]string{}
		for _, attr := range exported.Attributes {
			attrs[string(attr.Key)] = attr.Value.AsString()
		}
		var output map[string]any
		if err := json.Unmarshal([]byte(attrs["neatlogs.output.value"]), &output); err != nil {
			t.Fatal(err)
		}
		switch exported.Name {
		case "workflow":
			if output["status"] != "passed" || output["model_calls"] == nil || attrs["neatlogs.trace.output"] != attrs["neatlogs.workflow.output"] || attrs["neatlogs.workflow.input"] == "" {
				t.Fatal("workflow result or diagnostics lost")
			}
		case "tool":
			if output["order_count"] != float64(1) || attrs["neatlogs.trace.output"] != "" || attrs["neatlogs.tool.output"] != attrs["neatlogs.output.value"] || attrs["neatlogs.tool.input"] == "" {
				t.Fatal("tool output masquerades as workflow result")
			}
		default:
			t.Fatalf("unexpected span: %s", exported.Name)
		}
	}
}
