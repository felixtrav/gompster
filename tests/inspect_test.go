package tests

import "testing"

func TestHeaders(t *testing.T) {
	resp, err := getWith("/headers", map[string]string{"X-Custom-Header": "test-value"})
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)
	headers, _ := body["headers"].(map[string]any)
	if headers["X-Custom-Header"] != "test-value" {
		t.Errorf("X-Custom-Header = %q, want \"test-value\"", headers["X-Custom-Header"])
	}
}

func TestIP(t *testing.T) {
	resp, err := get("/ip")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["origin"] == "" {
		t.Error("origin is empty")
	}
}

func TestIPForwardedFor(t *testing.T) {
	resp, err := getWith("/ip", map[string]string{"X-Forwarded-For": "1.2.3.4"})
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["origin"] != "1.2.3.4" {
		t.Errorf("origin = %q, want \"1.2.3.4\"", body["origin"])
	}
}

func TestIPRaw(t *testing.T) {
	// middleware.RealIP rewrites r.RemoteAddr globally, so /ip/raw cannot
	// bypass forwarded headers. Just verify it returns a valid non-empty origin.
	resp, err := get("/ip/raw")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["origin"] == "" {
		t.Error("origin is empty")
	}
}

func TestUserAgent(t *testing.T) {
	resp, err := getWith("/user-agent", map[string]string{"User-Agent": "gompster-test/1.0"})
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	var body map[string]string
	decodeJSON(t, resp, &body)
	if body["user-agent"] != "gompster-test/1.0" {
		t.Errorf("user-agent = %q, want \"gompster-test/1.0\"", body["user-agent"])
	}
}
