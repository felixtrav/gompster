// gompster is a lightweight HTTP debugging server in the spirit of httpbin.
// It echoes requests, generates configurable responses, and exposes an
// OpenAPI spec + Swagger UI at the root.
//
// To extend the API, implement a handler.Group and append it to the groups
// slice in main — routes and spec entries are wired automatically.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/felixtrav/gompster/internal/config"
	"github.com/felixtrav/gompster/internal/handler"
	"github.com/felixtrav/gompster/internal/server"
)

func main() {
	cfg := config.Load()

	// Register route groups here. Adding a new handler.Group is all it takes
	// to extend the API — routes and OpenAPI spec are wired automatically.
	groups := []handler.Group{
		handler.HealthGroup(),
		handler.EchoGroup(),
		handler.AnythingGroup(),
		handler.InspectGroup(),
		handler.ResponseGroup(),
		handler.FormatsGroup(),
		handler.DynamicGroup(),
		handler.UUIDGroup(),
		handler.CookiesGroup(),
		handler.RedirectGroup(),
		handler.AuthGroup(),
		handler.StreamGroup(),
		handler.EncodingGroup(),
		handler.ImageGroup(),
	}

	srv := server.New(cfg, groups)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("gompster listening on http://%s", addr)
	log.Printf("API spec:   http://%s/openapi.json", addr)
	log.Printf("Swagger UI: http://%s/", addr)

	// NotifyContext cancels ctx on SIGINT/SIGTERM. Receiving the signal does
	// not stop the server immediately — it just unblocks <-ctx.Done(). The
	// actual drain of in-flight requests happens in srv.Shutdown below, which
	// waits up to 10 s for active connections to finish before returning.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gracefully…")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
