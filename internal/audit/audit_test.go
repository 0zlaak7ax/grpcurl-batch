package audit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/your-org/grpcurl-batch/internal/audit"
)

func decodeLines(t *testing.T, buf *bytes.Buffer) []audit.Event {
	t.Helper()
	var events []audit.Event
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		var e audit.Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("unmarshal line %q: %v", line, err)
		}
		events = append(events, e)
	}
	return events
}

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	err := l.Record(audit.Event{
		Kind:    audit.EventRequest,
		Method:  "/svc.Foo/Bar",
		Success: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := decodeLines(t, &buf)
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	if events[0].Kind != audit.EventRequest {
		t.Errorf("want kind %q, got %q", audit.EventRequest, events[0].Kind)
	}
	if events[0].Method != "/svc.Foo/Bar" {
		t.Errorf("unexpected method: %s", events[0].Method)
	}
}

func TestRecordRequest_SetsFields(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	labels := map[string]string{"env": "prod"}
	if err := l.RecordRequest("cid-1", "/svc.X/Y", labels); err != nil {
		t.Fatal(err)
	}
	events := decodeLines(t, &buf)
	e := events[0]
	if e.CorrelationID != "cid-1" {
		t.Errorf("correlation_id: got %q", e.CorrelationID)
	}
	if e.Labels["env"] != "prod" {
		t.Errorf("label env: got %q", e.Labels["env"])
	}
	if !e.Success {
		t.Error("expected success=true")
	}
}

func TestRecordResponse_WithError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.RecordResponse("cid-2", "/svc.A/B", 2, 150, errors.New("deadline exceeded")); err != nil {
		t.Fatal(err)
	}
	events := decodeLines(t, &buf)
	e := events[0]
	if e.Success {
		t.Error("expected success=false")
	}
	if e.Error != "deadline exceeded" {
		t.Errorf("error field: got %q", e.Error)
	}
	if e.LatencyMS != 150 {
		t.Errorf("latency_ms: got %d", e.LatencyMS)
	}
	if e.Attempt != 2 {
		t.Errorf("attempt: got %d", e.Attempt)
	}
}

func TestRecord_Concurrent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_ = l.Record(audit.Event{Kind: audit.EventRetry, Method: "/x/y", Success: false})
		}()
	}
	wg.Wait()
	events := decodeLines(t, &buf)
	if len(events) != n {
		t.Errorf("want %d events, got %d", n, len(events))
	}
}
