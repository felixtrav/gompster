package tests

import (
	"encoding/base64"
	"testing"
)

func TestBase64Decode(t *testing.T) {
	input := "Hello, gompster!"
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	resp, err := get("/base64/" + encoded)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	got := bodyString(t, resp)
	if got != input {
		t.Errorf("decoded = %q, want %q", got, input)
	}
}

func TestBase64URLSafeDecode(t *testing.T) {
	input := "url safe test?"
	encoded := base64.URLEncoding.EncodeToString([]byte(input))

	resp, err := get("/base64/" + encoded)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	got := bodyString(t, resp)
	if got != input {
		t.Errorf("decoded = %q, want %q", got, input)
	}
}

func TestBase64Invalid(t *testing.T) {
	resp, err := get("/base64/!!!invalid!!!")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 400)
	_ = resp.Body.Close()
}
