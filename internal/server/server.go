// Package server wires the HTTP server, registers all route groups, and
// exposes the generated OpenAPI spec at /openapi.json and Swagger UI at /.
package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/felixtrav/gompster/internal/config"
	"github.com/felixtrav/gompster/internal/handler"
)

// New creates a fully configured *http.Server from cfg and the provided route groups.
// To add new endpoints, append another handler.Group to the groups slice in main.go.
func New(cfg *config.Config, groups []handler.Group) *http.Server {
	r := chi.NewRouter()

	// middleware.RealIP mutates r.RemoteAddr to the leftmost X-Forwarded-For value.
	// IP spoofing is an accepted trade-off for a debugging tool — the /ip endpoints
	// are explicitly designed to reflect whatever the client or proxy sends.
	r.Use(middleware.RealIP) //nolint:staticcheck
	r.Use(middleware.Recoverer)
	if cfg.LogRequests {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	for _, g := range groups {
		for _, route := range g.Routes {
			if route.Method == "" {
				r.HandleFunc(route.Pattern, route.Handler)
			} else {
				r.Method(route.Method, route.Pattern, route.Handler)
			}
		}
	}

	spec := generateSpec(groups)
	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_, _ = w.Write(spec)
	})
	r.Get("/", swaggerUIHandler)

	return &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
