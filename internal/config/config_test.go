package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/grpcurl-batch/internal/config"
)

const validYAML = `
requests:
  - name: ping
    address: localhost:50051
    service: example.EchoService
    method: Echo
    data: '{"message": "hello"}'
    headers:
      authorization: Bearer token123
retry:
  max_attempts: 3
  delay: 500ms
output:
  format: json
  verbose: true
`

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "batch.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempFile(t, validYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(cfg.Requests))
	}
	if cfg.Requests[0].Name != "ping" {
		t.Errorf("expected name 'ping', got %q", cfg.Requests[0].Name)
	}
	if cfg.Retry.MaxAttempts != 3 {
		t.Errorf("expected max_attempts 3, got %d", cfg.Retry.MaxAttempts)
	}
	if cfg.Retry.Delay != 500*time.Millisecond {
		t.Errorf("expected delay 500ms, got %v", cfg.Retry.Delay)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Output.Format)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTempFile(t, `
requests:
  - name: bad
    service: svc.Foo
    method: Bar
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing address")
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTempFile(t, `
requests:
  - address: localhost:9090
    service: svc.Foo
    method: Bar
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Retry.MaxAttempts != 1 {
		t.Errorf("expected default max_attempts 1, got %d", cfg.Retry.MaxAttempts)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected default format 'json', got %q", cfg.Output.Format)
	}
}
