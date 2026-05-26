package tests

import (
	"regexp"
	"testing"
)

// uuidRegex validates a UUID v4 string. The literal "4" in the third group
// asserts version 4; the [89ab] nibble in the fourth group asserts the RFC 4122
// variant (the two most-significant bits must be 10xx).
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestUUID(t *testing.T) {
	resp, err := get("/uuid")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 200)

	var body map[string]any
	decodeJSON(t, resp, &body)
	id, _ := body["uuid"].(string)
	if !uuidRegex.MatchString(id) {
		t.Errorf("uuid = %q is not a valid v4 UUID", id)
	}
}

func TestUUIDUnique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 10; i++ {
		resp, err := get("/uuid")
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]any
		decodeJSON(t, resp, &body)
		id, _ := body["uuid"].(string)
		if seen[id] {
			t.Errorf("duplicate UUID: %s", id)
		}
		seen[id] = true
	}
}
