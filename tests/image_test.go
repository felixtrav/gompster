package tests

import (
	"strings"
	"testing"
)

func TestImageContentNegotiation(t *testing.T) {
	cases := []struct {
		accept      string
		wantCT      string
	}{
		{"image/jpeg", "image/jpeg"},
		{"image/png", "image/png"},
		{"image/svg+xml", "image/svg+xml"},
		{"image/webp", "image/webp"},
		{"image/gif", "image/gif"},
		// */* defaults to PNG
		{"*/*", "image/png"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.accept, func(t *testing.T) {
			resp, err := getWith("/image", map[string]string{"Accept": tc.accept})
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			assertStatus(t, resp, 200)
			if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, tc.wantCT) {
				t.Errorf("Content-Type = %q, want prefix %q", ct, tc.wantCT)
			}
			assertHeaderPresent(t, resp, "X-Attribution-Title")
			assertHeaderPresent(t, resp, "X-Attribution-Author")
			assertHeaderPresent(t, resp, "X-Attribution-License")
		})
	}
}

func TestImageUnsupportedAccept(t *testing.T) {
	resp, err := getWith("/image", map[string]string{"Accept": "image/tiff"})
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 406)
	_ = resp.Body.Close()
}

func TestImageByType(t *testing.T) {
	types := []struct {
		slug string
		ct   string
	}{
		{"jpeg", "image/jpeg"},
		{"png", "image/png"},
		{"svg", "image/svg+xml"},
		{"webp", "image/webp"},
		{"gif", "image/gif"},
		{"apng", "image/apng"},
	}

	for _, tc := range types {
		tc := tc
		t.Run(tc.slug, func(t *testing.T) {
			resp, err := get("/image/" + tc.slug)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			assertStatus(t, resp, 200)
			if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, tc.ct) {
				t.Errorf("Content-Type = %q, want prefix %q", ct, tc.ct)
			}
			assertHeaderPresent(t, resp, "X-Attribution-Title")
		})
	}
}

func TestImageOriginalWorkHeaders(t *testing.T) {
	// gif and apng are CC0 derivatives; they should carry X-Attribution-Original-* headers.
	for _, slug := range []string{"gif", "apng"} {
		slug := slug
		t.Run(slug, func(t *testing.T) {
			resp, err := get("/image/" + slug)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			assertStatus(t, resp, 200)
			assertHeaderPresent(t, resp, "X-Attribution-Original-Title")
			assertHeaderPresent(t, resp, "X-Attribution-Original-Author")
		})
	}
}

func TestImageUnsupportedType(t *testing.T) {
	resp, err := get("/image/tiff")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 404)
	_ = resp.Body.Close()
}

func TestImageAttribution(t *testing.T) {
	resp, err := get("/image/attribution")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "Content-Type", "application/json")

	var body map[string]any
	decodeJSON(t, resp, &body)
	for _, key := range []string{"jpeg", "png", "svg", "webp", "gif", "apng"} {
		if _, ok := body[key]; !ok {
			t.Errorf("attribution map missing key %q", key)
		}
	}
}
