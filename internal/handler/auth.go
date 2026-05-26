package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// AuthGroup returns routes for testing HTTP authentication schemes.
func AuthGroup() Group {
	return Group{
		Tag:         "Auth",
		Description: "Endpoints for testing HTTP authentication and authorization flows.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/basic-auth/{user}/{passwd}",
				Handler:     basicAuthHandler,
				Summary:     "HTTP Basic Auth",
				Description: "Challenges with HTTP Basic Auth. Returns 200 if credentials match, 401 otherwise.",
				Tags:        []string{"Auth"},
				Params: []Param{
					{Name: "user", In: "path", Description: "Expected username", Required: true, Schema: Schema{Type: "string"}},
					{Name: "passwd", In: "path", Description: "Expected password", Required: true, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{
					200: {Description: "Authenticated", ContentType: "application/json"},
					401: {Description: "Unauthorized"},
				},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/bearer",
				Handler:     bearerHandler,
				Summary:     "Bearer token auth",
				Description: "Returns 200 and echoes the token if a Bearer token is present in the Authorization header.",
				Tags:        []string{"Auth"},
				Params: []Param{
					{Name: "Authorization", In: "header", Description: "Bearer <token>", Required: true, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{
					200: {Description: "Token accepted", ContentType: "application/json"},
					401: {Description: "Missing or invalid Bearer token"},
				},
			},
		},
	}
}

func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	expectedUser := chi.URLParam(r, "user")
	expectedPass := chi.URLParam(r, "passwd")

	user, pass, ok := r.BasicAuth()
	if !ok || user != expectedUser || pass != expectedPass {
		w.Header().Set("WWW-Authenticate", `Basic realm="gompster"`)
		JSON(w, http.StatusUnauthorized, map[string]any{
			"authenticated": false,
			"user":          user,
		})
		return
	}
	JSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"user":          user,
	})
}

func bearerHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		w.Header().Set("WWW-Authenticate", `Bearer realm="gompster"`)
		JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid Bearer token"})
		return
	}
	JSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"token":         strings.TrimPrefix(auth, "Bearer "),
	})
}
