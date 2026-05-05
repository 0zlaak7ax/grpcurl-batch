package reporter

import (
	"errors"
	"os"
	"strings"
	"testing"

	"grpcurl-batch/internal/formatter"
)

func TestWriteHTML_AllSuccess(t *testing.T) {
	results := []formatter.Result{
		{Method: "pkg.Svc/MethodA", Output: `{"id":1}`, Attempts: 1},
		{Method: "pkg.Svc/MethodB", Output: `{"id":2}`, Attempts: 1},
	}

	tmp := t.TempDir() + "/report.html"
	if err := WriteHTML(tmp, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := readFile(t, tmp)
	assertContains(t, body, "pkg.Svc/MethodA")
	assertContains(t, body, "pkg.Svc/MethodB")
	assertContains(t, body, "PASS")
	if strings.Contains(body, "FAIL") {
		t.Error("expected no FAIL in all-success report")
	}
	assertContains(t, body, "Total:</strong> 2")
	assertContains(t, body, "Passed:</strong> 2")
}

func TestWriteHTML_WithFailure(t *testing.T) {
	results := []formatter.Result{
		{Method: "pkg.Svc/Good", Output: `{}`, Attempts: 1},
		{Method: "pkg.Svc/Bad", Err: errors.New("deadline exceeded"), Attempts: 3},
	}

	tmp := t.TempDir() + "/report.html"
	if err := WriteHTML(tmp, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := readFile(t, tmp)
	assertContains(t, body, "FAIL")
	assertContains(t, body, "deadline exceeded")
	assertContains(t, body, "Failed:</strong> 1")
	assertContains(t, body, "Passed:</strong> 1")
}

func TestWriteHTML_Empty(t *testing.T) {
	tmp := t.TempDir() + "/report.html"
	if err := WriteHTML(tmp, []formatter.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := readFile(t, tmp)
	assertContains(t, body, "Total:</strong> 0")
}

func TestWriteHTML_InvalidPath(t *testing.T) {
	err := WriteHTML("/nonexistent/dir/report.html", nil)
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

// helpers

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file %s: %v", path, err)
	}
	return string(b)
}

func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if !strings.Contains(body, substr) {
		t.Errorf("expected HTML to contain %q", substr)
	}
}
