package tests

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestEchoGet(t *testing.T) {
	resp, err := getWith("/get?foo=bar", map[string]string{"X-Test": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)

	if body["method"] != "GET" {
		t.Errorf("method = %q, want GET", body["method"])
	}
	args, _ := body["args"].(map[string]any)
	if args["foo"] != "bar" {
		t.Errorf("args.foo = %q, want \"bar\"", args["foo"])
	}
	headers, _ := body["headers"].(map[string]any)
	if headers["X-Test"] != "hello" {
		t.Errorf("headers.X-Test = %q, want \"hello\"", headers["X-Test"])
	}
}

func TestEchoMethods(t *testing.T) {
	cases := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/post"},
		{http.MethodPut, "/put"},
		{http.MethodPatch, "/patch"},
		{http.MethodDelete, "/delete"},
	}
	for _, tc := range cases {
		t.Run(tc.method, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, ts.URL+tc.path, bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			assertStatus(t, resp, 200)
			var body map[string]any
			decodeJSON(t, resp, &body)
			if body["method"] != tc.method {
				t.Errorf("method = %q, want %q", body["method"], tc.method)
			}
		})
	}
}

func TestEchoPostJSONBody(t *testing.T) {
	payload := `{"hello":"world"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/post", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["data"] != payload {
		t.Errorf("data = %q, want %q", body["data"], payload)
	}
	jsonField, _ := body["json"].(map[string]any)
	if jsonField["hello"] != "world" {
		t.Errorf("json.hello = %q, want \"world\"", jsonField["hello"])
	}
}

func TestEchoWrongMethod(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/get", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 405)
	_ = resp.Body.Close()
}
