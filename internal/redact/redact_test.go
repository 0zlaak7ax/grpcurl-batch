package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/user/grpcurl-batch/internal/redact"
)

func TestNew_NoFields_PassThrough(t *testing.T) {
	r := redact.New(nil)
	input := `{"password":"secret123"}`
	if got := r.Scrub(input); got != input {
		t.Fatalf("expected passthrough, got %s", got)
	}
}

func TestScrub_SensitiveFieldReplaced(t *testing.T) {
	r := redact.New([]string{"password", "token"})
	input := `{"username":"alice","password":"hunter2","token":"abc123"}`
	out := r.Scrub(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["password"] != "[REDACTED]" {
		t.Errorf("password not redacted, got %v", m["password"])
	}
	if m["token"] != "[REDACTED]" {
		t.Errorf("token not redacted, got %v", m["token"])
	}
	if m["username"] != "alice" {
		t.Errorf("username should be unchanged, got %v", m["username"])
	}
}

func TestScrub_CaseInsensitive(t *testing.T) {
	r := redact.New([]string{"Secret"})
	input := `{"SECRET":"val"}`
	out := r.Scrub(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["SECRET"] != "[REDACTED]" {
		t.Errorf("expected redaction, got %v", m["SECRET"])
	}
}

func TestScrub_NestedObject(t *testing.T) {
	r := redact.New([]string{"apikey"})
	input := `{"meta":{"apikey":"xyz","env":"prod"}}`
	out := r.Scrub(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	meta := m["meta"].(map[string]interface{})
	if meta["apikey"] != "[REDACTED]" {
		t.Errorf("nested apikey not redacted, got %v", meta["apikey"])
	}
	if meta["env"] != "prod" {
		t.Errorf("env should be unchanged")
	}
}

func TestScrub_ArrayElements(t *testing.T) {
	r := redact.New([]string{"token"})
	input := `[{"token":"a"},{"token":"b","id":1}]`
	out := r.Scrub(input)

	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for i, item := range arr {
		if item["token"] != "[REDACTED]" {
			t.Errorf("arr[%d].token not redacted", i)
		}
	}
}

func TestScrub_NonJSON_PassThrough(t *testing.T) {
	r := redact.New([]string{"password"})
	input := "not json at all"
	if got := r.Scrub(input); got != input {
		t.Fatalf("expected passthrough for non-JSON, got %s", got)
	}
}
