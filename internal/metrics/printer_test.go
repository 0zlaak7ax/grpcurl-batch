package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/metrics"
)

func TestPrint_ContainsExpectedFields(t *testing.T) {
	c := metrics.New()
	c.Record(true, 100*time.Millisecond)
	c.Record(true, 200*time.Millisecond)
	c.Record(false, 50*time.Millisecond)

	var buf bytes.Buffer
	if err := metrics.Print(&buf, c.Snapshot()); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}
	out := buf.String()

	checks := []string{
		"Total requests:",
		"3",
		"Successful:",
		"2",
		"Failed:",
		"1",
		"Total duration:",
		"Avg latency:",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, out)
		}
	}
}

func TestPrint_ZeroSnapshot(t *testing.T) {
	var buf bytes.Buffer
	s := metrics.Snapshot{}
	if err := metrics.Print(&buf, s); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "0") {
		t.Errorf("expected zero values in output: %s", buf.String())
	}
}
