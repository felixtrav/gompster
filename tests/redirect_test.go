package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRelativeRedirect(t *testing.T) {
	resp, err := get("/redirect/3")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 302)

	loc := resp.Header.Get("Location")
	if !strings.Contains(loc, "/redirect/2") {
		t.Errorf("Location = %q, expected to contain /redirect/2", loc)
	}
}

// TestRelativeRedirectFollowed uses redir (not client) because redir is
// configured to follow redirects automatically; client stops at the first
// 3xx response. After following the full chain the final destination is /get,
// which returns 200.
func TestRelativeRedirectFollowed(t *testing.T) {
	resp, err := redir.Get(ts.URL + "/redirect/3")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)
}

func TestRelativeRedirectZero(t *testing.T) {
	// /redirect/0 goes straight to /get.
	resp, err := redir.Get(ts.URL + "/redirect/0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)
}

func TestAbsoluteRedirect(t *testing.T) {
	resp, err := get("/absolute-redirect/2")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 302)

	loc := resp.Header.Get("Location")
	parsed, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("Location is not a valid URL: %v", err)
	}
	if !parsed.IsAbs() {
		t.Errorf("Location %q is not an absolute URL", loc)
	}
}

func TestRedirectTo(t *testing.T) {
	target := "http://example.com/target"
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

	for _, method := range methods {
		method := method
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, ts.URL+"/redirect-to?url="+url.QueryEscape(target), nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther {
				t.Errorf("status = %d, want 302 or 303", resp.StatusCode)
			}
			if loc := resp.Header.Get("Location"); loc != target {
				t.Errorf("Location = %q, want %q", loc, target)
			}
		})
	}
}

func TestRedirectToCustomStatus(t *testing.T) {
	target := "http://example.com"
	resp, err := get(fmt.Sprintf("/redirect-to?url=%s&status_code=307", url.QueryEscape(target)))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 307)
	assertHeader(t, resp, "Location", target)
}
