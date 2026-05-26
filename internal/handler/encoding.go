package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// EncodingGroup returns routes for encoding and decoding utilities.
func EncodingGroup() Group {
	return Group{
		Tag:         "Encoding",
		Description: "Endpoints for encoding and decoding data in various formats.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/base64/{value}",
				Handler:     base64DecodeHandler,
				Summary:     "Decode a base64 value",
				Description: "Decodes a standard or URL-safe base64 path segment and returns the raw content as text/plain.",
				Tags:        []string{"Encoding"},
				Params: []Param{
					{Name: "value", In: "path", Description: "Base64-encoded string", Required: true, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{
					200: {Description: "Decoded content", ContentType: "text/plain"},
					400: {Description: "Invalid base64 input"},
				},
			},
		},
	}
}

func base64DecodeHandler(w http.ResponseWriter, r *http.Request) {
	encoded := chi.URLParam(r, "value")

	// Try standard encoding first, then URL-safe.
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(encoded)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid base64 input"})
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(decoded)
}
