// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the concurrency and throughput of batched gRPC requests.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which requests are dispatched.
type Limiter struct {
	tokens chan struct{}
	interval time.Duration
}

// Config holds the configuration for the rate limiter.
type Config struct {
	// MaxConcurrent is the maximum number of in-flight requests at any time.
	MaxConcurrent int
	// Interval is the minimum delay between acquiring successive tokens.
	// Zero means no delay beyond concurrency limits.
	Interval time.Duration
}

// New creates a new Limiter. MaxConcurrent must be >= 1.
func New(cfg Config) *Limiter {
	if cfg.MaxConcurrent < 1 {
		cfg.MaxConcurrent = 1
	}
	tokens := make(chan struct{}, cfg.MaxConcurrent)
	for i := 0; i < cfg.MaxConcurrent; i++ {
		tokens <- struct{}{}
	}
	return &Limiter{
		tokens:   tokens,
		interval: cfg.Interval,
	}
}

// Acquire blocks until a token is available or ctx is cancelled.
// Returns an error if the context is done before a token is obtained.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
		if l.interval > 0 {
			select {
			case <-time.After(l.interval):
			case <-ctx.Done():
				l.Release()
				return ctx.Err()
			}
		}
		return nil
	}
}

// Release returns a token to the pool, allowing another caller to proceed.
func (l *Limiter) Release() {
	l.tokens <- struct{}{}
}
