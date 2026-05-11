package observe_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/grpcurl-batch/internal/observe"
)

func TestRecord_SuccessSpan(t *testing.T) {
	var buf bytes.Buffer
	obs := observe.New(&buf)

	err := obs.Record(context.Background(), "ping", nil, func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	spans := obs.Spans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if !spans[0].OK() {
		t.Error("expected span to be OK")
	}
	if spans[0].Name != "ping" {
		t.Errorf("expected name 'ping', got %q", spans[0].Name)
	}
	if !strings.Contains(buf.String(), "ok") {
		t.Errorf("expected 'ok' in output, got %q", buf.String())
	}
}

func TestRecord_FailureSpan(t *testing.T) {
	var buf bytes.Buffer
	obs := observe.New(&buf)
	sentinel := errors.New("boom")

	err := obs.Record(context.Background(), "call", map[string]string{"svc": "foo"}, func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	spans := obs.Spans()
	if spans[0].OK() {
		t.Error("expected span to be not OK")
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected error message in output, got %q", buf.String())
	}
}

func TestRecord_DurationPositive(t *testing.T) {
	obs := observe.New(new(bytes.Buffer))
	_ = obs.Record(context.Background(), "op", nil, func(_ context.Context) error { return nil })
	if d := obs.Spans()[0].Duration(); d < 0 {
		t.Errorf("expected non-negative duration, got %v", d)
	}
}

func TestReset_ClearsSpans(t *testing.T) {
	obs := observe.New(new(bytes.Buffer))
	_ = obs.Record(context.Background(), "a", nil, func(_ context.Context) error { return nil })
	obs.Reset()
	if n := len(obs.Spans()); n != 0 {
		t.Errorf("expected 0 spans after reset, got %d", n)
	}
}

func TestRecord_Concurrent(t *testing.T) {
	obs := observe.New(new(bytes.Buffer))
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = obs.Record(context.Background(), "concurrent", nil, func(_ context.Context) error { return nil })
		}()
	}
	wg.Wait()
	if n := len(obs.Spans()); n != workers {
		t.Errorf("expected %d spans, got %d", workers, n)
	}
}

func TestSpans_ReturnsCopy(t *testing.T) {
	obs := observe.New(new(bytes.Buffer))
	_ = obs.Record(context.Background(), "x", nil, func(_ context.Context) error { return nil })
	s1 := obs.Spans()
	s1[0].Name = "mutated"
	s2 := obs.Spans()
	if s2[0].Name == "mutated" {
		t.Error("Spans() should return a copy, not a reference")
	}
}
