// Package circuitbreaker provides a simple circuit breaker that opens after
// a configurable number of consecutive failures and resets after a cooldown.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Breaker is a circuit breaker that tracks consecutive failures.
type Breaker struct {
	mu          sync.Mutex
	maxFailures int
	cooldown    time.Duration
	failures    int
	state       State
	openedAt    time.Time
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts recovery after the cooldown period.
func New(maxFailures int, cooldown time.Duration) *Breaker {
	if maxFailures <= 0 {
		maxFailures = 3
	}
	if cooldown <= 0 {
		cooldown = 10 * time.Second
	}
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
		state:       StateClosed,
	}
}

// Allow reports whether a call should be allowed through.
// It transitions an open breaker to half-open once the cooldown has elapsed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful call, resetting the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call. After maxFailures the breaker opens.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
