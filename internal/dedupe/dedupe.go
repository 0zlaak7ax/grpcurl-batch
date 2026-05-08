// Package dedupe provides request deduplication for gRPC batch runs.
// It tracks which requests have already been executed (by key) within
// a single batch session and skips duplicates.
package dedupe

import (
	"sync"
)

// Filter tracks seen request keys and reports whether a given key
// has already been processed.
type Filter struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// New returns an initialised, empty Filter.
func New() *Filter {
	return &Filter{
		seen: make(map[string]struct{}),
	}
}

// IsDuplicate returns true if key has been seen before and records it
// as seen. The first call for a given key always returns false.
func (f *Filter) IsDuplicate(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.seen[key]; ok {
		return true
	}
	f.seen[key] = struct{}{}
	return false
}

// Reset clears all recorded keys, allowing all requests to be
// processed again.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]struct{})
}

// Len returns the number of unique keys recorded so far.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}
