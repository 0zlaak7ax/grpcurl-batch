package redact_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/grpcurl-batch/internal/redact"
)

// TestScrub_DeepNesting verifies redaction works at arbitrary depth.
func TestScrub_DeepNesting(t *testing.T) {
	r := redact.New([]string{"password"})
	input := `{"a":{"b":{"c":{"password":"deep"}}}}`
	out := r.Scrub(input)

	if strings.Contains(out, "deep") {
		t.Errorf("deeply nested password should be redacted, got: %s", out)
	}
}

// TestScrub_LargePayload ensures performance is acceptable for bigger objects.
func TestScrub_LargePayload(t *testing.T) {
	r := redact.New([]string{"secret"})

	records := make([]map[string]interface{}, 500)
	for i := range records {
		records[i] = map[string]interface{}{
			"id":     i,
			"name":   "user",
			"secret": "sensitive",
		}
	}
	b, _ := json.Marshal(records)
	out := r.Scrub(string(b))

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for i, item := range result {
		if item["secret"] != "[REDACTED]" {
			t.Errorf("record %d: secret not redacted", i)
		}
	}
}

// TestScrub_MultipleCallsIdempotent ensures repeated scrubbing is safe.
func TestScrub_MultipleCallsIdempotent(t *testing.T) {
	r := redact.New([]string{"token"})
	input := `{"token":"abc"}`
	first := r.Scrub(input)
	second := r.Scrub(first)

	if first != second {
		t.Errorf("scrub should be idempotent: first=%s second=%s", first, second)
	}
	if !strings.Contains(first, "[REDACTED]") {
		t.Errorf("expected redaction marker in output")
	}
}
