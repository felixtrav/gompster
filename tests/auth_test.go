package tests

import (
	"net/http"
	"testing"
)

func TestBasicAuthSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/basic-auth/alice/s3cr3t", nil)
	req.SetBasicAuth("alice", "s3cr3t")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["authenticated"] != true {
		t.Error("authenticated should be true")
	}
	if body["user"] != "alice" {
		t.Errorf("user = %q, want \"alice\"", body["user"])
	}
}

func TestBasicAuthWrongPassword(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/basic-auth/alice/s3cr3t", nil)
	req.SetBasicAuth("alice", "wrong")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 401)
	_ = resp.Body.Close()
}

func TestBasicAuthMissing(t *testing.T) {
	resp, err := get("/basic-auth/user/pass")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 401)
	if resp.Header.Get("WWW-Authenticate") == "" {
		t.Error("missing WWW-Authenticate header")
	}
	_ = resp.Body.Close()
}

func TestBearerAuthSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/bearer", nil)
	req.Header.Set("Authorization", "Bearer mytoken123")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)
	if body["authenticated"] != true {
		t.Error("authenticated should be true")
	}
	if body["token"] != "mytoken123" {
		t.Errorf("token = %q, want \"mytoken123\"", body["token"])
	}
}

func TestBearerAuthMissing(t *testing.T) {
	resp, err := get("/bearer")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 401)
	_ = resp.Body.Close()
}
