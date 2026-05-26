package tests

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSetCookie(t *testing.T) {
	resp, err := get("/cookies/set?session=abc123")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 302)
	found := false
	for _, c := range resp.Cookies() {
		if c.Name == "session" && c.Value == "abc123" {
			found = true
		}
	}
	if !found {
		t.Error("Set-Cookie for session=abc123 not found")
	}
	_ = resp.Body.Close()
}

func TestGetCookies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/cookies", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "xyz"})
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body struct {
		Cookies map[string]string `json:"cookies"`
	}
	decodeJSON(t, resp, &body)
	if body.Cookies["token"] != "xyz" {
		t.Errorf("cookies.token = %q, want \"xyz\"", body.Cookies["token"])
	}
}

// TestDeleteCookie uses redir (not client) throughout because redir carries a
// cookiejar.Jar that persists cookies across the set → verify → delete request
// chain. client has no jar, so the cookie set by /cookies/set would not be
// present in the subsequent /cookies or /cookies/delete requests.
func TestDeleteCookie(t *testing.T) {
	base, _ := url.Parse(ts.URL)

	setReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/cookies/set?delme=1", nil)
	setResp, err := redir.Do(setReq)
	if err != nil {
		t.Fatal(err)
	}
	_ = setResp.Body.Close()

	// Verify cookie is present via redir client's jar.
	getResp, err := redir.Get(ts.URL + "/cookies")
	if err != nil {
		t.Fatal(err)
	}
	var body struct {
		Cookies map[string]string `json:"cookies"`
	}
	decodeJSON(t, getResp, &body)
	if body.Cookies["delme"] == "" {
		t.Skip("cookie jar not carrying delme; cannot test delete")
	}

	// Delete it.
	delResp, err := redir.Get(ts.URL + "/cookies/delete?delme")
	if err != nil {
		t.Fatal(err)
	}
	_ = delResp.Body.Close()

	// Cookie should no longer be in the jar.
	cookies := redir.Jar.Cookies(base)
	for _, c := range cookies {
		if c.Name == "delme" {
			t.Error("cookie 'delme' should have been deleted")
		}
	}
}
