package tests

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGzip(t *testing.T) {
	resp, err := get("/gzip")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Encoding", "gzip")
	assertHeader(t, resp, "Content-Type", "application/json")

	body := decodeGzip(t, resp)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("gzip body is not valid JSON: %v", err)
	}
	if payload["project"] != "gompster" {
		t.Errorf("project = %q, want \"gompster\"", payload["project"])
	}
	if payload["encoding"] != "gzip" {
		t.Errorf("encoding = %q, want \"gzip\"", payload["encoding"])
	}
}

func TestDeflate(t *testing.T) {
	resp, err := get("/deflate")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Encoding", "deflate")
	assertHeader(t, resp, "Content-Type", "application/json")

	body := decodeDeflate(t, resp)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("deflate body is not valid JSON: %v", err)
	}
	if payload["encoding"] != "deflate" {
		t.Errorf("encoding = %q, want \"deflate\"", payload["encoding"])
	}
}

func TestBrotli(t *testing.T) {
	resp, err := get("/brotli")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Encoding", "br")
	assertHeader(t, resp, "Content-Type", "application/json")

	body := decodeBrotli(t, resp)
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("brotli body is not valid JSON: %v", err)
	}
	if payload["encoding"] != "brotli" {
		t.Errorf("encoding = %q, want \"brotli\"", payload["encoding"])
	}
}

func TestUTF8(t *testing.T) {
	resp, err := get("/utf8")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q, want text/plain", ct)
	}
	_ = resp.Body.Close()
}

func TestHTML(t *testing.T) {
	resp, err := get("/html")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	body := bodyString(t, resp)
	if !strings.Contains(body, "<html") {
		t.Error("response does not look like HTML")
	}
}

func TestJSONFormat(t *testing.T) {
	resp, err := get("/json")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/json")

	var payload map[string]any
	decodeJSON(t, resp, &payload)
	if payload["project"] != "gompster" {
		t.Errorf("project = %q, want \"gompster\"", payload["project"])
	}
}

func TestXMLFormat(t *testing.T) {
	resp, err := get("/xml")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/xml")
	body := bodyString(t, resp)
	if !strings.Contains(body, "gompster") {
		t.Error("XML body does not contain \"gompster\"")
	}
}

func TestRobotsTxt(t *testing.T) {
	resp, err := get("/robots.txt")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	body := bodyString(t, resp)
	if !strings.Contains(body, "Disallow: /deny") {
		t.Error("robots.txt does not contain \"Disallow: /deny\"")
	}
}

func TestDeny(t *testing.T) {
	resp, err := get("/deny")
	if err != nil {
		t.Fatal(err)
	}
	// /deny returns 200 with plain text content, disallowed via robots.txt
	assertStatus(t, resp, 200)
	body := bodyString(t, resp)
	if !strings.Contains(body, "YOU SHALL NOT PASS") {
		t.Error("deny page does not contain expected text")
	}
}
