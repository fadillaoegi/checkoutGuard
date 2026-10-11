package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"checkoutguard/internal/adapters/browser"
	"checkoutguard/internal/adapters/model"
	"checkoutguard/internal/adapters/postgres"
	"checkoutguard/internal/adapters/telemetry"
	"checkoutguard/internal/application"
	"checkoutguard/internal/config"
	"checkoutguard/internal/delivery/httpapi"
	"github.com/google/uuid"
)

func main() {
	for _, key := range []string{"NEATLOGS_ENV_FILE", "MODEL_ENV_FILE", "MODEL_PRICING_ENV_FILE"} {
		if err := config.LoadEnv(os.Getenv(key)); err != nil {
			log.Fatalf("cannot load %s", key)
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	store, err := postgres.OpenStore(ctx, config.Env("DATABASE_URL", "postgres://checkoutguard:checkoutguard_local@127.0.0.1:55433/checkoutguard?sslmode=disable"))
	if err != nil {
		log.Fatal("cannot connect to dedicated PostgreSQL: ", err)
	}
	defer store.Close()
	shutdown := telemetry.Init(ctx)
	defer shutdown()
	addr := config.Env("HTTP_ADDR", "127.0.0.1:8080")
	host, _, err := net.SplitHostPort(addr)
	if err != nil || (host != "127.0.0.1" && host != "localhost") {
		log.Fatal("HTTP_ADDR must bind the local loopback interface")
	}
	planner := model.NewModel()
	service := &application.Service{Store: store, Planner: planner, Browsers: browser.Factory{}, Markets: &browser.Gateway{}, Tracer: telemetry.Tracer{}, NewID: uuid.NewString, BaseURL: "http://" + addr, EvidenceDir: config.Env("EVIDENCE_DIR", "evidence")}
	handler := &httpapi.Handler{Service: service, BaseURL: service.BaseURL, EvidenceDir: service.EvidenceDir}
	server := &http.Server{Addr: addr, Handler: handler.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(stopCtx)
	}()
	log.Printf("CheckoutGuard backend http://%s; sandbox and marketplace preparation; model configured=%t", addr, planner.Ready())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print("server error: ", err)
	}
	service.Close()
}
