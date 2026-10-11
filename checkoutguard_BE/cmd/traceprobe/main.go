package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"checkoutguard/internal/adapters/telemetry"
	"checkoutguard/internal/config"
	neatlogs "github.com/neatlogs/neatlogs-go"
)

func main() {
	if err := config.LoadEnv(os.Getenv("NEATLOGS_ENV_FILE")); err != nil {
		log.Fatal("cannot load telemetry configuration")
	}
	if len(os.Args) == 3 && os.Args[1] == "--doctor" {
		if os.Args[2] != "local" && os.Args[2] != "probe" {
			log.Fatal("doctor mode must be local or probe")
		}
		command := exec.Command("bin/neatlogs", "doctor", "--"+os.Args[2], "--json")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			os.Exit(1)
		}
		return
	}
	ctx := context.Background()
	shutdown := telemetry.Init(ctx)
	defer shutdown()
	ctx, span, end := neatlogs.Trace(ctx, "checkoutguard.bootstrap")
	log.Printf("bootstrap trace ID: %s", span.SpanContext().TraceID().String())
	_ = neatlogs.SetTraceInput(span, map[string]any{"purpose": "verify CheckoutGuard integration", "synthetic": true})
	_, child, childEnd := neatlogs.StartSpan(ctx, "sandbox.readiness", "tool")
	_ = neatlogs.SetTraceInput(child, map[string]any{"operation": "bootstrap"})
	_ = neatlogs.SetTraceOutput(child, map[string]any{"ready": true})
	childEnd()
	_ = neatlogs.SetTraceOutput(span, map[string]any{"status": "bootstrap_complete", "application_workflow": false})
	end()
	if err := neatlogs.Flush(ctx); err != nil {
		log.Fatal("trace flush failed")
	}
	log.Print("bootstrap spans flushed; hosted readback must be verified separately")
}
