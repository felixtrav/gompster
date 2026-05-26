// healthcheck is a minimal liveness probe binary designed to be embedded
// in Docker images (including distroless) where curl/wget are unavailable.
//
// Exit codes:
//
//	0  — /health returned HTTP 200
//	1  — connection error or non-200 response
//
// Environment variables (mirror the server's own vars):
//
//	PORT  — port the server is listening on (default: 80)
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := "80"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	url := fmt.Sprintf("http://127.0.0.1:%s/health", port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "healthcheck: unexpected status %d\n", resp.StatusCode)
		os.Exit(1)
	}
}
