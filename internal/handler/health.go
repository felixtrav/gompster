package handler

import (
	"net/http"
	"time"
)

// serverStart is captured once at package initialisation (program start) and
// used by healthHandler to compute the server's uptime on each request.
var serverStart = time.Now()

func HealthGroup() Group {
	return Group{
		Tag:         "Health",
		Description: "Liveness and readiness probes.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/health",
				Summary:     "Health check",
				Description: "Returns `200 OK` with server status and uptime. Designed for use with Docker / Kubernetes liveness probes.",
				Responses: map[int]Response{
					200: {Description: "Server is healthy", ContentType: "application/json"},
				},
				Handler: healthHandler,
			},
		},
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]any{
		"status":         "ok",
		"uptime_seconds": time.Since(serverStart).Seconds(),
	})
}
