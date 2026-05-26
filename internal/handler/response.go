package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ResponseGroup returns routes for controlling the HTTP response.
func ResponseGroup() Group {
	return Group{
		Tag:         "Responses",
		Description: "Endpoints for generating responses with specific status codes or headers.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/status/{code}",
				Handler:     statusHandler,
				Summary:     "Return a given status code",
				Description: "Returns a response with the specified HTTP status code (100–599).",
				Tags:        []string{"Responses"},
				Params: []Param{
					{
						Name:        "code",
						In:          "path",
						Description: "HTTP status code to return",
						Required:    true,
						Schema:      Schema{Type: "integer", Minimum: f64p(100), Maximum: f64p(599)},
					},
				},
				Responses: map[int]Response{
					200: {Description: "Response with the requested status code"},
				},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/response-headers",
				Handler:     responseHeadersHandler,
				Summary:     "Return given response headers",
				Description: "Sets query parameters as response headers and echoes them as JSON.",
				Tags:        []string{"Responses"},
				Params: []Param{
					{Name: "key", In: "query", Description: "Any key=value pair to echo as a response header", Required: false, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{
					200: {Description: "Response headers set and echoed", ContentType: "application/json"},
				},
			},
		},
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.Atoi(chi.URLParam(r, "code"))
	if err != nil || code < 100 || code > 599 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid status code, must be 100–599"})
		return
	}
	JSON(w, code, map[string]any{
		"code":        code,
		"description": http.StatusText(code),
	})
}

func responseHeadersHandler(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]string)
	for k, vs := range r.URL.Query() {
		w.Header().Set(k, vs[0])
		result[k] = vs[0]
	}
	JSON(w, http.StatusOK, result)
}
