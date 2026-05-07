// Package sampling provides request sampling for grpcurl-batch,
// allowing a fraction of requests to be selected for detailed tracing
// or logging without processing every request.
package sampling

import (
	"math/rand"
	"sync"
)

// Sampler decides whether a given request should be sampled.
type Sampler interface {
	Sample() bool
}

// Config holds configuration for a rate-based sampler.
type Config struct {
	// Rate is the fraction of requests to sample, in the range [0.0, 1.0].
	// 0.0 means never sample; 1.0 means always sample.
	Rate float64
}

type sampler struct {
	mu   sync.Mutex
	rng  *rand.Rand
	rate float64
}

// New returns a Sampler that samples requests at the given rate.
// Rate is clamped to [0.0, 1.0].
func New(cfg Config) Sampler {
	rate := cfg.Rate
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &sampler{
		// Use a local source so tests can be deterministic via seed injection.
		rng:  rand.New(rand.NewSource(rand.Int63())), //nolint:gosec
		rate: rate,
	}
}

// Sample returns true if the current request should be sampled.
func (s *sampler) Sample() bool {
	if s.rate == 0 {
		return false
	}
	if s.rate == 1 {
		return true
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}
