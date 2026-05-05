package formatter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewTemplateFormatter_InvalidTemplate(t *testing.T) {
	_, err := NewTemplateFormatter(&bytes.Buffer{}, "{{ .Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestTemplateFormatter_WriteSuccess(t *testing.T) {
	var buf bytes.Buffer
	tf, err := NewTemplateFormatter(&buf, "[{{ checkmark .Success }}] {{ .Method }} ({{ durationMs .Duration }}ms)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := Result{
		Method:   "helloworld.Greeter/SayHello",
		Address:  "localhost:50051",
		Success:  true,
		Output:   `{"message": "Hello"}`,
		Attempts: 1,
		Duration: 42 * time.Millisecond,
	}
	if err := tf.Write(r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "✓") {
		t.Errorf("expected checkmark in output, got: %q", got)
	}
	if !strings.Contains(got, "SayHello") {
		t.Errorf("expected method name in output, got: %q", got)
	}
	if !strings.Contains(got, "42ms") {
		t.Errorf("expected duration in output, got: %q", got)
	}
}

func TestTemplateFormatter_WriteFailure(t *testing.T) {
	var buf bytes.Buffer
	tf, err := NewTemplateFormatter(&buf, "{{ checkmark .Success }} {{ .Method }}: {{ .Error }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := Result{
		Method:  "pkg.Svc/Method",
		Address: "localhost:9090",
		Success: false,
		Error:   "connection refused",
		Attempts: 3,
	}
	if err := tf.Write(r); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "✗") {
		t.Errorf("expected failure mark in output, got: %q", got)
	}
	if !strings.Contains(got, "connection refused") {
		t.Errorf("expected error message in output, got: %q", got)
	}
}

func TestTemplateFormatter_CustomFields(t *testing.T) {
	var buf bytes.Buffer
	tf, err := NewTemplateFormatter(&buf, "attempts={{ .Attempts }} addr={{ .Address }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := Result{
		Method:   "svc/Op",
		Address:  "remote:443",
		Success:  true,
		Attempts: 2,
	}
	_ = tf.Write(r)

	got := buf.String()
	if !strings.Contains(got, "attempts=2") {
		t.Errorf("expected attempts=2 in output, got: %q", got)
	}
	if !strings.Contains(got, "addr=remote:443") {
		t.Errorf("expected address in output, got: %q", got)
	}
}
