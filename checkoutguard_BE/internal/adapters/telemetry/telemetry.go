package telemetry

import (
	"context"
	"log"
	"os"
	"time"

	neatlogs "github.com/neatlogs/neatlogs-go"
)

func Init(ctx context.Context) func() {
	shutdown, err := neatlogs.Init(ctx, neatlogs.Config{
		APIKey: os.Getenv("NEATLOGS_API_KEY"), Endpoint: os.Getenv("NEATLOGS_ENDPOINT"),
		WorkflowName: "checkoutguard", Tags: []string{"neatHack", "checkoutguard"},
	})
	if err != nil {
		log.Print("telemetry unavailable; application continues")
		return func() {}
	}
	return func() {
		flushCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := shutdown(flushCtx); err != nil {
			log.Print("telemetry shutdown failed")
		}
	}
}
