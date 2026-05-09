// Package audit provides a structured audit log for recording gRPC batch
// request events, outcomes, and metadata for compliance and debugging.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventRequest  EventKind = "request"
	EventResponse EventKind = "response"
	EventRetry    EventKind = "retry"
	EventSkip     EventKind = "skip"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp     time.Time         `json:"timestamp"`
	Kind          EventKind         `json:"kind"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Method        string            `json:"method"`
	Attempt       int               `json:"attempt,omitempty"`
	Success       bool              `json:"success"`
	LatencyMS     int64             `json:"latency_ms,omitempty"`
	Error         string            `json:"error,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
}

// Logger writes structured audit events to an io.Writer.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New creates a new audit Logger writing to w.
func New(w io.Writer) *Logger {
	return &Logger{out: w}
}

// Record encodes e as a JSON line and writes it to the underlying writer.
// It is safe for concurrent use.
func (l *Logger) Record(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	if err != nil {
		return fmt.Errorf("audit: write event: %w", err)
	}
	return nil
}

// RecordRequest is a convenience wrapper for EventRequest events.
func (l *Logger) RecordRequest(correlationID, method string, labels map[string]string) error {
	return l.Record(Event{
		Kind:          EventRequest,
		CorrelationID: correlationID,
		Method:        method,
		Success:       true,
		Labels:        labels,
	})
}

// RecordResponse is a convenience wrapper for EventResponse events.
func (l *Logger) RecordResponse(correlationID, method string, attempt int, latencyMS int64, err error) error {
	e := Event{
		Kind:          EventResponse,
		CorrelationID: correlationID,
		Method:        method,
		Attempt:       attempt,
		LatencyMS:     latencyMS,
		Success:       err == nil,
	}
	if err != nil {
		e.Error = err.Error()
	}
	return l.Record(e)
}
