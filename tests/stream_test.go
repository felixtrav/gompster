package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"testing"
)

func TestStreamLines(t *testing.T) {
	const n = 5
	resp, err := get(fmt.Sprintf("/stream/%d", n))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	assertStatus(t, resp, 200)

	// NDJSON is one JSON object per line. Scanning line-by-line and
	// unmarshalling each line independently is the correct way to consume
	// this format — the response is not a single top-level JSON array.
	scanner := bufio.NewScanner(resp.Body)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Errorf("line %d is not valid JSON: %v — %q", count, err, line)
		}
		count++
	}
	if count != n {
		t.Errorf("stream returned %d lines, want %d", count, n)
	}
}

func TestStreamInvalidN(t *testing.T) {
	resp, err := get("/stream/-1")
	if err != nil {
		t.Fatal(err)
	}
	assertStatus(t, resp, 400)
	_ = resp.Body.Close()
}
