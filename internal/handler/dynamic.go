package handler

import (
	"crypto/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	maxBytes        = 100 * 1024 * 1024 // 100 MB
	maxDelaySeconds = 300.0
)

// DynamicGroup returns routes for generating random or time-based content.
func DynamicGroup() Group {
	return Group{
		Tag:         "Dynamic Data",
		Description: "Endpoints that generate random or dynamically computed content.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/bytes/{n}",
				Handler:     bytesHandler,
				Summary:     "Return n random bytes",
				Description: "Returns n random bytes as application/octet-stream. Capped at 100 MB.",
				Tags:        []string{"Dynamic Data"},
				Params: []Param{
					{Name: "n", In: "path", Description: "Number of bytes", Required: true, Schema: Schema{Type: "integer", Minimum: f64p(0)}},
				},
				Responses: map[int]Response{200: {Description: "Random bytes", ContentType: "application/octet-stream"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/delay/{seconds}",
				Handler:     delayHandler,
				Summary:     "Delay response",
				Description: "Waits the given number of seconds (float) before responding. Capped at 300 s.",
				Tags:        []string{"Dynamic Data"},
				Params: []Param{
					{Name: "seconds", In: "path", Description: "Delay in seconds", Required: true, Schema: Schema{Type: "number", Minimum: f64p(0), Maximum: f64p(maxDelaySeconds)}},
				},
				Responses: map[int]Response{200: {Description: "Echoed request after delay", ContentType: "application/json"}},
			},
		},
	}
}

func bytesHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(chi.URLParam(r, "n"))
	if err != nil || n < 0 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "n must be a non-negative integer"})
		return
	}
	if n > maxBytes {
		n = maxBytes
	}
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(n))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	sec, err := strconv.ParseFloat(chi.URLParam(r, "seconds"), 64)
	if err != nil || sec < 0 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "seconds must be a non-negative number"})
		return
	}
	if sec > maxDelaySeconds {
		sec = maxDelaySeconds
	}
	// Two cases: the timer fires and we fall through to send the response, or
	// the client disconnects (context cancelled) and we return early so we
	// don't hold the goroutine open for the remainder of the delay.
	select {
	case <-time.After(time.Duration(sec * float64(time.Second))):
	case <-r.Context().Done():
		return
	}
	JSON(w, http.StatusOK, CaptureRequest(r))
}
