package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// EchoResponse mirrors the httpbin-style request echo payload.
type EchoResponse struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data,omitempty"`
	Files   map[string]string `json:"files"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	JSON    any               `json:"json,omitempty"`
	Method  string            `json:"method"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

// JSON writes v as indented JSON with the given HTTP status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// CaptureRequest builds an EchoResponse from an incoming request.
func CaptureRequest(r *http.Request) EchoResponse {
	resp := EchoResponse{
		Args:    queryToMap(r.URL.Query()),
		Headers: headersToMap(r.Header),
		Method:  r.Method,
		Origin:  ClientIP(r),
		URL:     RequestURL(r),
		Files:   make(map[string]string),
		Form:    make(map[string]string),
	}

	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		body, _ := io.ReadAll(r.Body)
		resp.Data = string(body)
		var parsed any
		if json.Unmarshal(body, &parsed) == nil {
			resp.JSON = parsed
		}

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		_ = r.ParseForm()
		for k, vs := range r.Form {
			resp.Form[k] = vs[0]
		}

	case strings.HasPrefix(ct, "multipart/form-data"):
		_ = r.ParseMultipartForm(32 << 20)
		if r.MultipartForm != nil {
			for k, vs := range r.MultipartForm.Value {
				resp.Form[k] = vs[0]
			}
			for k, fhs := range r.MultipartForm.File {
				if len(fhs) > 0 {
					resp.Files[k] = fhs[0].Filename
				}
			}
		}

	default:
		body, _ := io.ReadAll(r.Body)
		resp.Data = string(body)
	}

	return resp
}

// ClientIP extracts the real client IP, respecting proxy headers.
func ClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		return addr[:i]
	}
	return addr
}

// RequestURL reconstructs the full request URL.
func RequestURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = r.Header.Get("X-Forwarded-Host")
	}
	return scheme + "://" + host + r.URL.RequestURI()
}

func queryToMap(q url.Values) map[string]string {
	m := make(map[string]string, len(q))
	for k, vs := range q {
		m[k] = vs[0]
	}
	return m
}

func headersToMap(h http.Header) map[string]string {
	m := make(map[string]string, len(h))
	for k, vs := range h {
		m[k] = strings.Join(vs, ", ")
	}
	return m
}
