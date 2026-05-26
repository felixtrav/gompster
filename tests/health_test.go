package tests

import (
	"testing"
)

func TestHealth(t *testing.T) {
	resp, err := get("/health")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/json")

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["status"] != "ok" {
		t.Errorf("status = %q, want \"ok\"", body["status"])
	}
	if _, ok := body["uptime_seconds"]; !ok {
		t.Error("missing uptime_seconds field")
	}
}
