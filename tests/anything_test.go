package tests

import (
	"net/http"
	"testing"
)

func TestAnything(t *testing.T) {
	resp, err := get("/anything")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["method"] != "GET" {
		t.Errorf("method = %q, want GET", body["method"])
	}
}

func TestAnythingAllMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, ts.URL+"/anything", nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			assertStatus(t, resp, 200)
			var body map[string]any
			decodeJSON(t, resp, &body)
			if body["method"] != method {
				t.Errorf("method = %q, want %q", body["method"], method)
			}
		})
	}
}

func TestAnythingSubpath(t *testing.T) {
	resp, err := get("/anything/some/nested/path")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	_ = resp.Body.Close()
}
