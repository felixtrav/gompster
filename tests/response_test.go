package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestStatusCodes(t *testing.T) {
	for _, code := range []int{200, 201, 400, 401, 403, 404, 500, 503} {
		code := code
		t.Run(http.StatusText(code), func(t *testing.T) {
			resp, err := get(fmt.Sprintf("/status/%d", code))
			if err != nil {
				t.Fatal(err)
			}
			assertStatus(t, resp, code)
			_ = resp.Body.Close()
		})
	}
}

func TestStatusInvalid(t *testing.T) {
	resp, err := get("/status/999")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 400)
	_ = resp.Body.Close()
}

func TestResponseHeaders(t *testing.T) {
	resp, err := get("/response-headers?X-My-Header=hello&X-Other=world")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)
	assertHeader(t, resp, "X-My-Header", "hello")
	assertHeader(t, resp, "X-Other", "world")

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["X-My-Header"] != "hello" {
		t.Errorf("body X-My-Header = %q, want \"hello\"", body["X-My-Header"])
	}
}
