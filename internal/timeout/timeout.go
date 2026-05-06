// Package timeout provides per-request deadline management for gRPC batch calls.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrDeadlineExceeded is returned when a request exceeds its configured deadline.
var ErrDeadlineExceeded = errors.New("timeout: deadline exceeded")

// Config holds timeout settings for a single request or a batch.
type Config struct {
	// PerRequest is the maximum duration allowed for a single gRPC call.
	PerRequest time.Duration
	// Total is the maximum duration allowed for the entire batch run.
	Total time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PerRequest: 30 * time.Second,
		Total:      5 * time.Minute,
	}
}

// Limiter wraps a parent context with deadline enforcement.
type Limiter struct {
	cfg Config
}

// New creates a Limiter from the given Config.
// If PerRequest or Total are zero they are replaced by DefaultConfig values.
func New(cfg Config) *Limiter {
	def := DefaultConfig()
	if cfg.PerRequest <= 0 {
		cfg.PerRequest = def.PerRequest
	}
	if cfg.Total <= 0 {
		cfg.Total = def.Total
	}
	return &Limiter{cfg: cfg}
}

// WithTotal wraps ctx with the configured total deadline.
// The caller is responsible for calling the returned cancel function.
func (l *Limiter) WithTotal(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, l.cfg.Total)
}

// WithRequest wraps ctx with the per-request deadline.
// The caller is responsible for calling the returned cancel function.
func (l *Limiter) WithRequest(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, l.cfg.PerRequest)
}

// WrapError converts a context deadline error into ErrDeadlineExceeded with
// additional context about which operation timed out.
func WrapError(operation string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %s", ErrDeadlineExceeded, operation)
	}
	return err
}
