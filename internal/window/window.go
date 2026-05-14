// Package window provides a sliding-window counter for tracking event
// rates over a rolling time interval. It is safe for concurrent use.
package window

import (
	"sync"
	"time"
)

// Window tracks how many events have occurred within a rolling duration.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  int
	counts   []int64
	times    []time.Time
	cursor   int
	now      func() time.Time
}

// New creates a Window that divides size into buckets equal-width slots.
// Panics if size <= 0 or buckets < 1.
func New(size time.Duration, buckets int) *Window {
	if size <= 0 {
		panic("window: size must be positive")
	}
	if buckets < 1 {
		panic("window: buckets must be >= 1")
	}
	return &Window{
		size:    size,
		buckets: buckets,
		counts:  make([]int64, buckets),
		times:   make([]time.Time, buckets),
		now:     time.Now,
	}
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(w.now())
	w.counts[w.cursor] += n
}

// Count returns the total number of events recorded within the window.
func (w *Window) Count() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.now()
	w.advance(now)
	cutoff := now.Add(-w.size)
	var total int64
	for i, t := range w.times {
		if !t.IsZero() && t.After(cutoff) {
			total += w.counts[i]
		}
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.counts = make([]int64, w.buckets)
	w.times = make([]time.Time, w.buckets)
	w.cursor = 0
}

// advance moves the cursor to the current bucket, evicting stale slots.
func (w *Window) advance(now time.Time) {
	bucketSize := w.size / time.Duration(w.buckets)
	// Determine which bucket slot now falls into.
	slot := int(now.UnixNano()/int64(bucketSize)) % w.buckets
	if w.cursor != slot || w.times[slot].IsZero() {
		// Clear slots that are older than one full window.
		cutoff := now.Add(-w.size)
		for i := range w.counts {
			if w.times[i].Before(cutoff) {
				w.counts[i] = 0
				w.times[i] = time.Time{}
			}
		}
		w.cursor = slot
	}
	if w.times[w.cursor].IsZero() {
		w.times[w.cursor] = now
	}
}
