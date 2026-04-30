package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/formatter"
)

func makeResult(success bool) formatter.Result {
	return formatter.Result{
		Method:   "pkg.Service/Method",
		Success:  success,
		Attempts: 2,
		Duration: 150 * time.Millisecond,
		Output:   `{"id":"1"}`,
		Error:    "",
	}
}

func TestWrite_TextSuccess(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	if err := f.Write(makeResult(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", out)
	}
	if !strings.Contains(out, "pkg.Service/Method") {
		t.Errorf("expected method name in output, got: %s", out)
	}
}

func TestWrite_TextFailure(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	r := makeResult(false)
	r.Error = "deadline exceeded"
	if err := f.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[FAIL]") {
		t.Errorf("expected [FAIL] in output, got: %s", out)
	}
	if !strings.Contains(out, "deadline exceeded") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatJSON, &buf)
	if err := f.Write(makeResult(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["method"] != "pkg.Service/Method" {
		t.Errorf("unexpected method: %v", got["method"])
	}
	if got["success"] != true {
		t.Errorf("expected success=true, got: %v", got["success"])
	}
}

func TestWrite_Summary(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatSummary, &buf)
	if err := f.Write(makeResult(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "✓") {
		t.Errorf("expected ✓ in summary output, got: %s", out)
	}
}

func TestWrite_SummaryFail(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatSummary, &buf)
	if err := f.Write(makeResult(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "✗") {
		t.Errorf("expected ✗ in summary output, got: %s", out)
	}
}
