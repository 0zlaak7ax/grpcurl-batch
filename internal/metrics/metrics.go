// Package metrics provides lightweight in-memory counters for tracking
// batch execution statistics such as total requests, successes, failures,
// and cumulative latency.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Snapshot holds a point-in-time copy of collected metrics.
type Snapshot struct {
	Total    int64
	Success  int64
	Failure  int64
	TotalDur time.Duration
}

// AvgLatency returns the mean latency per request, or zero if no requests
// have been recorded.
func (s Snapshot) AvgLatency() time.Duration {
	if s.Total == 0 {
		return 0
	}
	return time.Duration(int64(s.TotalDur) / s.Total)
}

// Collector accumulates metrics from concurrent goroutines.
type Collector struct {
	total   atomic.Int64
	success atomic.Int64
	failure atomic.Int64

	durMu   sync.Mutex
	totalNs int64
}

// New returns an initialised *Collector.
func New() *Collector {
	return &Collector{}
}

// Record registers a single request outcome.
func (c *Collector) Record(success bool, dur time.Duration) {
	c.total.Add(1)
	if success {
		c.success.Add(1)
	} else {
		c.failure.Add(1)
	}
	c.durMu.Lock()
	c.totalNs += int64(dur)
	c.durMu.Unlock()
}

// Snapshot returns a consistent point-in-time view of the metrics.
func (c *Collector) Snapshot() Snapshot {
	c.durMu.Lock()
	ns := c.totalNs
	c.durMu.Unlock()
	return Snapshot{
		Total:    c.total.Load(),
		Success:  c.success.Load(),
		Failure:  c.failure.Load(),
		TotalDur: time.Duration(ns),
	}
}
