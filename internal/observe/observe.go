// Package observe provides a lightweight span-style observation layer
// that records the start, end, and outcome of discrete operations for
// structured diagnostic output without a full tracing dependency.
package observe

import (
	"context"
	"io"
	"sync"
	"time"
)

// Span represents a single observed operation.
type Span struct {
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Err       error
	Attrs     map[string]string
}

// Duration returns the elapsed time of the span.
func (s Span) Duration() time.Duration {
	return s.EndedAt.Sub(s.StartedAt)
}

// OK reports whether the span completed without error.
func (s Span) OK() bool { return s.Err == nil }

// Observer records spans produced by operations.
type Observer struct {
	mu    sync.Mutex
	spans []Span
	out   io.Writer
}

// New returns an Observer that writes structured lines to out.
// Pass io.Discard to suppress output and only accumulate spans.
func New(out io.Writer) *Observer {
	return &Observer{out: out}
}

// Record executes fn, records a Span for the operation named name,
// and returns the error fn produced (if any).
func (o *Observer) Record(ctx context.Context, name string, attrs map[string]string, fn func(ctx context.Context) error) error {
	start := time.Now()
	err := fn(ctx)
	s := Span{
		Name:      name,
		StartedAt: start,
		EndedAt:   time.Now(),
		Err:       err,
		Attrs:     attrs,
	}
	o.mu.Lock()
	o.spans = append(o.spans, s)
	o.mu.Unlock()
	o.write(s)
	return err
}

// Spans returns a snapshot of all recorded spans.
func (o *Observer) Spans() []Span {
	o.mu.Lock()
	defer o.mu.Unlock()
	out := make([]Span, len(o.spans))
	copy(out, o.spans)
	return out
}

// Reset clears all recorded spans.
func (o *Observer) Reset() {
	o.mu.Lock()
	o.spans = o.spans[:0]
	o.mu.Unlock()
}

func (o *Observer) write(s Span) {
	status := "ok"
	if s.Err != nil {
		status = "error: " + s.Err.Error()
	}
	// Best-effort; ignore write errors.
	_, _ = io.WriteString(o.out, s.StartedAt.Format(time.RFC3339)+" "+s.Name+" "+s.Duration().String()+" "+status+"\n")
}
