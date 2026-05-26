package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// RedirectGroup returns routes for testing HTTP redirect behavior.
func RedirectGroup() Group {
	return Group{
		Tag:         "Redirects",
		Description: "Endpoints for testing redirect chains and redirect targets.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/redirect/{n}",
				Handler:     redirectHandler,
				Summary:     "Redirect n times",
				Description: "Performs n relative redirects before landing at /get.",
				Tags:        []string{"Redirects"},
				Params: []Param{
					{Name: "n", In: "path", Description: "Number of redirects", Required: true, Schema: Schema{Type: "integer", Minimum: f64p(0)}},
				},
				Responses: map[int]Response{
					302: {Description: "Redirect"},
					200: {Description: "Final destination at /get"},
				},
			},
			{
				Method:      "",
				Pattern:     "/redirect-to",
				Handler:     redirectToHandler,
				Summary:     "Redirect to a URL",
				Description: "Redirects to the URL specified in the `url` query parameter. Accepts any HTTP method.",
				Tags:        []string{"Redirects"},
				Params: []Param{
					{Name: "url", In: "query", Description: "Target URL", Required: true, Schema: Schema{Type: "string"}},
					{Name: "status_code", In: "query", Description: "Redirect status code (default 302)", Required: false, Schema: Schema{Type: "integer", Default: 302}},
				},
				Responses: map[int]Response{302: {Description: "Redirect to target"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/absolute-redirect/{n}",
				Handler:     absoluteRedirectHandler,
				Summary:     "Absolute redirect n times",
				Description: "Performs n absolute URL redirects before landing at /get.",
				Tags:        []string{"Redirects"},
				Params: []Param{
					{Name: "n", In: "path", Description: "Number of redirects", Required: true, Schema: Schema{Type: "integer", Minimum: f64p(0)}},
				},
				Responses: map[int]Response{
					302: {Description: "Absolute redirect"},
					200: {Description: "Final destination at /get"},
				},
			},
		},
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(chi.URLParam(r, "n"))
	if err != nil || n < 0 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "n must be a non-negative integer"})
		return
	}
	hops := currentHops(r) + 1
	if n == 0 {
		http.Redirect(w, r, fmt.Sprintf("/get?hops=%d", hops-1), http.StatusFound)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/redirect/%d?hops=%d", n-1, hops), http.StatusFound)
}

func redirectToHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	if target == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "url query parameter is required"})
		return
	}
	code := http.StatusFound
	if sc := r.URL.Query().Get("status_code"); sc != "" {
		if c, err := strconv.Atoi(sc); err == nil && c >= 300 && c < 400 {
			code = c
		}
	}
	http.Redirect(w, r, target, code)
}

func absoluteRedirectHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(chi.URLParam(r, "n"))
	if err != nil || n < 0 {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "n must be a non-negative integer"})
		return
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	hops := currentHops(r) + 1
	var target string
	if n == 0 {
		target = fmt.Sprintf("%s://%s/get?hops=%d", scheme, host, hops-1)
	} else {
		target = fmt.Sprintf("%s://%s/absolute-redirect/%d?hops=%d", scheme, host, n-1, hops)
	}
	http.Redirect(w, r, target, http.StatusFound)
}

// currentHops reads the hops counter threaded through the redirect chain.
func currentHops(r *http.Request) int {
	h, _ := strconv.Atoi(r.URL.Query().Get("hops"))
	return h
}
