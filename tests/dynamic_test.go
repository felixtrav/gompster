package tests

import (
	"testing"
)

func TestRandomBytes(t *testing.T) {
	resp, err := get("/bytes/128")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/octet-stream")
}

func TestRandomBytesZero(t *testing.T) {
	resp, err := get("/bytes/0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)
}

func TestDelayZero(t *testing.T) {
	resp, err := get("/delay/0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)
	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["method"] != "GET" {
		t.Errorf("method = %q, want GET", body["method"])
	}
}

func TestDelayInvalid(t *testing.T) {
	resp, err := get("/delay/-1")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 400)
	_ = resp.Body.Close()
}
