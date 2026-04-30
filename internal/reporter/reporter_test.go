package reporter_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/nicholasgasior/grpcurl-batch/internal/formatter"
	"github.com/nicholasgasior/grpcurl-batch/internal/reporter"
)

func makeResult(method, address string, success bool, attempts int, output string) formatter.Result {
	return formatter.Result{
		Method:   method,
		Address:  address,
		Success:  success,
		Attempts: attempts,
		Output:   output,
	}
}

func TestWriteJUnit_AllSuccess(t *testing.T) {
	results := []formatter.Result{
		makeResult("pkg.Service/MethodA", "localhost:50051", true, 1, `{"ok":true}`),
		makeResult("pkg.Service/MethodB", "localhost:50051", true, 1, `{"ok":true}`),
	}

	var buf bytes.Buffer
	if err := reporter.WriteJUnit(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var suites reporter.JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	if len(suites.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(suites.Suites))
	}
	s := suites.Suites[0]
	if s.Tests != 2 {
		t.Errorf("expected Tests=2, got %d", s.Tests)
	}
	if s.Failures != 0 {
		t.Errorf("expected Failures=0, got %d", s.Failures)
	}
}

func TestWriteJUnit_WithFailure(t *testing.T) {
	results := []formatter.Result{
		makeResult("pkg.Service/MethodA", "localhost:50051", true, 1, `{"ok":true}`),
		makeResult("pkg.Service/MethodB", "localhost:50051", false, 3, "connection refused"),
	}

	var buf bytes.Buffer
	if err := reporter.WriteJUnit(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var suites reporter.JUnitTestSuites
	if err := xml.Unmarshal(buf.Bytes(), &suites); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	s := suites.Suites[0]
	if s.Failures != 1 {
		t.Errorf("expected Failures=1, got %d", s.Failures)
	}
	if s.TestCases[1].Failure == nil {
		t.Fatal("expected failure element on second test case")
	}
	if !strings.Contains(s.TestCases[1].Failure.Text, "connection refused") {
		t.Errorf("failure text missing output: %q", s.TestCases[1].Failure.Text)
	}
}

func TestWriteJUnit_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.WriteJUnit(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "testsuites") {
		t.Errorf("expected testsuites root element in output")
	}
}
