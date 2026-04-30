package formatter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/formatter"
)

func TestCollector_Empty(t *testing.T) {
	var c formatter.Collector
	rep := c.Report()
	if rep.Total != 0 || rep.Passed != 0 || rep.Failed != 0 {
		t.Errorf("expected zero counts, got %+v", rep)
	}
}

func TestCollector_MixedResults(t *testing.T) {
	var c formatter.Collector
	c.Add(formatter.Result{Method: "A", Success: true, Attempts: 1, Duration: 10 * time.Millisecond})
	c.Add(formatter.Result{Method: "B", Success: false, Attempts: 3, Duration: 300 * time.Millisecond, Error: "timeout"})
	c.Add(formatter.Result{Method: "C", Success: true, Attempts: 1, Duration: 20 * time.Millisecond})

	rep := c.Report()
	if rep.Total != 3 {
		t.Errorf("expected Total=3, got %d", rep.Total)
	}
	if rep.Passed != 2 {
		t.Errorf("expected Passed=2, got %d", rep.Passed)
	}
	if rep.Failed != 1 {
		t.Errorf("expected Failed=1, got %d", rep.Failed)
	}
}

func TestPrintSummary(t *testing.T) {
	var c formatter.Collector
	c.Add(formatter.Result{Method: "X", Success: true, Attempts: 1})
	c.Add(formatter.Result{Method: "Y", Success: false, Attempts: 2, Error: "err"})

	var buf bytes.Buffer
	if err := formatter.PrintSummary(&buf, c.Report()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Total: 2") {
		t.Errorf("expected 'Total: 2' in output, got: %s", out)
	}
	if !strings.Contains(out, "Passed: 1") {
		t.Errorf("expected 'Passed: 1' in output, got: %s", out)
	}
	if !strings.Contains(out, "Failed: 1") {
		t.Errorf("expected 'Failed: 1' in output, got: %s", out)
	}
}
