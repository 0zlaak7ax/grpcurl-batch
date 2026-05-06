// Package backoff provides configurable backoff strategies for retry delays.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines how delay is computed between retry attempts.
type Strategy interface {
	Delay(attempt int) time.Duration
}

// Exponential implements exponential backoff with optional jitter.
type Exponential struct {
	Base    time.Duration
	Max     time.Duration
	Jitter  bool
	Factor  float64
}

// New returns an Exponential strategy with sensible defaults.
func New(base, max time.Duration, jitter bool) *Exponential {
	return &Exponential{
		Base:   base,
		Max:    max,
		Jitter: jitter,
		Factor: 2.0,
	}
}

// Delay returns the wait duration for the given attempt (0-indexed).
func (e *Exponential) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	mult := math.Pow(e.Factor, float64(attempt))
	d := time.Duration(float64(e.Base) * mult)
	if d > e.Max {
		d = e.Max
	}
	if e.Jitter && d > 0 {
		// add up to 25% random jitter
		jitter := time.Duration(rand.Int63n(int64(d) / 4))
		d += jitter
	}
	return d
}

// Fixed returns a constant delay regardless of attempt number.
type Fixed struct {
	Interval time.Duration
}

// Delay always returns the fixed interval.
func (f *Fixed) Delay(_ int) time.Duration {
	return f.Interval
}

// Linear increases delay linearly: base * (attempt + 1), capped at max.
type Linear struct {
	Base time.Duration
	Max  time.Duration
}

// Delay returns base*(attempt+1) capped at Max.
func (l *Linear) Delay(attempt int) time.Duration {
	d := l.Base * time.Duration(attempt+1)
	if l.Max > 0 && d > l.Max {
		return l.Max
	}
	return d
}
