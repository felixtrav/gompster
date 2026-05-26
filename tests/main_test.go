// Package tests contains integration tests for the gompster HTTP server.
// Tests run against a real httptest.Server with all route groups registered,
// so chi URL params, middleware, and the full request lifecycle are exercised.
package tests

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/felixtrav/gompster/internal/config"
	"github.com/felixtrav/gompster/internal/handler"
	"github.com/felixtrav/gompster/internal/server"
)

var (
	ts     *httptest.Server
	client *http.Client // does not follow redirects
	redir  *http.Client // follows redirects, shares a cookie jar
)

func TestMain(m *testing.M) {
	cfg := config.Load()
	groups := []handler.Group{
		handler.HealthGroup(),
		handler.EchoGroup(),
		handler.AnythingGroup(),
		handler.InspectGroup(),
		handler.ResponseGroup(),
		handler.FormatsGroup(),
		handler.DynamicGroup(),
		handler.UUIDGroup(),
		handler.CookiesGroup(),
		handler.RedirectGroup(),
		handler.AuthGroup(),
		handler.StreamGroup(),
		handler.EncodingGroup(),
		handler.ImageGroup(),
	}
	srv := server.New(cfg, groups)
	ts = httptest.NewServer(srv.Handler)

	jar, _ := cookiejar.New(nil)
	// DisableCompression prevents the transport from transparently decompressing
	// gzip responses, so compressed-encoding tests can read the raw bytes.
	client = &http.Client{
		Transport: &http.Transport{DisableCompression: true},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	redir = &http.Client{
		Transport: &http.Transport{DisableCompression: true},
		Jar:       jar,
	}

	code := m.Run()
	ts.Close()
	os.Exit(code)
}

// get performs a GET request against the test server and returns the response.
func get(path string) (*http.Response, error) {
	return client.Get(ts.URL + path)
}

// getWith performs a GET with custom request headers.
func getWith(path string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

// decodeJSON reads and JSON-decodes the response body into v.
func decodeJSON(t *testing.T, r *http.Response, v any) {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
}

// decodeGzip returns the decompressed body of a gzip-encoded response.
func decodeGzip(t *testing.T, r *http.Response) []byte {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	gr, err := gzip.NewReader(r.Body)
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer func() { _ = gr.Close() }()
	b, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("gzip read: %v", err)
	}
	return b
}

// decodeDeflate returns the decompressed body of a deflate-encoded response.
func decodeDeflate(t *testing.T, r *http.Response) []byte {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	fr := flate.NewReader(r.Body)
	defer func() { _ = fr.Close() }()
	b, err := io.ReadAll(fr)
	if err != nil {
		t.Fatalf("deflate read: %v", err)
	}
	return b
}

// decodeBrotli returns the decompressed body of a brotli-encoded response.
func decodeBrotli(t *testing.T, r *http.Response) []byte {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	b, err := io.ReadAll(brotli.NewReader(r.Body))
	if err != nil {
		t.Fatalf("brotli read: %v", err)
	}
	return b
}

// bodyString reads and returns the full response body as a string.
func bodyString(t *testing.T, r *http.Response) string {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

// assertStatus fails the test if resp.StatusCode != want.
func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Errorf("status = %d, want %d", resp.StatusCode, want)
	}
}

// assertHeader fails if the named response header does not equal want.
func assertHeader(t *testing.T, resp *http.Response, name, want string) {
	t.Helper()
	if got := resp.Header.Get(name); got != want {
		t.Errorf("header %q = %q, want %q", name, got, want)
	}
}

// assertHeaderPresent fails if the named response header is empty.
func assertHeaderPresent(t *testing.T, resp *http.Response, name string) {
	t.Helper()
	if resp.Header.Get(name) == "" {
		t.Errorf("header %q is missing", name)
	}
}
