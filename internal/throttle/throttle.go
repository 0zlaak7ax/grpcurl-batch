// Package throttle provides a token-bucket based request throttler
// that limits the rate of outgoing gRPC calls over a sliding window.
package throttle

import (
	"context"
	"sync"
	"time"
)

// Throttler controls the rate at which requests are dispatched.
type Throttler struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// Config holds configuration for a Throttler.
type Config struct {
	// Rate is the number of requests allowed per second.
	Rate float64
	// Burst is the maximum burst size (bucket capacity).
	Burst float64
}

// New creates a Throttler from cfg. If Rate or Burst are zero the
// throttler is effectively unlimited (1 000 000 req/s).
func New(cfg Config) *Throttler {
	if cfg.Rate <= 0 {
		cfg.Rate = 1_000_000
	}
	if cfg.Burst <= 0 {
		cfg.Burst = cfg.Rate
	}
	return &Throttler{
		tokens:   cfg.Burst,
		max:      cfg.Burst,
		rate:     cfg.Rate,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or ctx is cancelled.
func (t *Throttler) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		t.mu.Lock()
		now := t.clock()
		elapsed := now.Sub(t.lastTick).Seconds()
		t.lastTick = now
		t.tokens += elapsed * t.rate
		if t.tokens > t.max {
			t.tokens = t.max
		}
		if t.tokens >= 1 {
			t.tokens--
			t.mu.Unlock()
			return nil
		}
		wait := time.Duration((1-t.tokens)/t.rate*1e9) * time.Nanosecond
		t.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// Available returns the current token count (snapshot).
func (t *Throttler) Available() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tokens
}
