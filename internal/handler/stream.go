package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const maxStreamLines = 1000

// StreamGroup returns routes for testing streaming HTTP responses.
func StreamGroup() Group {
	return Group{
		Tag:         "Streaming",
		Description: "Endpoints that stream data progressively using chunked transfer encoding.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/stream/{n}",
				Handler:     streamHandler,
				Summary:     "Stream n JSON lines",
				Description: "Streams n newline-delimited JSON objects (NDJSON). Capped at 1000 lines.",
				Tags:        []string{"Streaming"},
				Params: []Param{
					{Name: "n", In: "path", Description: "Number of JSON lines to stream", Required: true, Schema: Schema{Type: "integer", Minimum: f64p(1), Maximum: f64p(maxStreamLines)}},
				},
				Responses: map[int]Response{200: {Description: "Streaming NDJSON", ContentType: "application/x-ndjson"}},
			},
		},
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(chi.URLParam(r, "n"))
	if err != nil || n < 0 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "n must be a non-negative integer"})
		return
	}
	if n > maxStreamLines {
		n = maxStreamLines
	}

	// Not every ResponseWriter supports flushing (e.g. some test wrappers).
	// The Flush call after each line is what causes data to be sent to the
	// client incrementally rather than buffered until the handler returns.
	flusher, ok := w.(http.Flusher)
	if !ok {
		JSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming not supported by this server"})
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	for i := 0; i < n; i++ {
		// Check for client disconnect before writing each line so we stop
		// streaming immediately if the connection is dropped mid-response.
		select {
		case <-r.Context().Done():
			return
		default:
		}
		_ = enc.Encode(map[string]any{
			"id":        i,
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			"origin":    ClientIP(r),
			"url":       RequestURL(r),
		})
		flusher.Flush()
	}
}
