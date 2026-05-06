// Package retry provides configurable retry logic with backoff strategies
// for use in grpcurl-batch request execution.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Policy defines the retry behaviour for a single operation.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// Multiplier scales the delay after each failure (exponential backoff).
	// A value of 1 produces constant delays.
	Multiplier float64
	// MaxDelay caps the computed delay.
	MaxDelay time.Duration
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		Multiplier:   2.0,
		MaxDelay:     5 * time.Second,
	}
}

// ErrExhausted is returned when all attempts have been consumed.
var ErrExhausted = errors.New("retry: all attempts exhausted")

// Do runs fn up to p.MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
// The last non-nil error from fn is wrapped with ErrExhausted.
func Do(ctx context.Context, p Policy, fn func(attempt int) error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.Multiplier < 1 {
		p.Multiplier = 1
	}

	var lastErr error
	for i := 0; i < p.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(i + 1); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < p.MaxAttempts-1 {
			delay := delay(p, i)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}
	return errors.Join(ErrExhausted, lastErr)
}

// delay computes the back-off duration for attempt index i (0-based).
func delay(p Policy, i int) time.Duration {
	d := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(i))
	if p.MaxDelay > 0 && time.Duration(d) > p.MaxDelay {
		return p.MaxDelay
	}
	return time.Duration(d)
}
