package tests

import (
	"strings"
	"testing"
)

func TestSpecJSON(t *testing.T) {
	resp, err := get("/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/json")

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["openapi"] != "3.0.3" {
		t.Errorf("openapi version = %q, want \"3.0.3\"", body["openapi"])
	}
	if _, ok := body["paths"]; !ok {
		t.Error("spec is missing \"paths\"")
	}
}

func TestSpecUI(t *testing.T) {
	resp, err := get("/")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	body := bodyString(t, resp)
	if !strings.Contains(body, "SwaggerUIBundle") {
		t.Error("Swagger UI HTML does not include SwaggerUIBundle")
	}
}
