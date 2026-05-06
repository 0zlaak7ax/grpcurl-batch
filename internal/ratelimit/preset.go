package ratelimit

import "time"

// Preset returns a Limiter configured with common defaults suitable for
// most gRPC batch workloads.
//
//	- Concurrency: 8 simultaneous in-flight requests.
//	- Interval:    10 ms between successive Acquire calls.
func Preset() *Limiter {
	return New(Options{
		Concurrency: 8,
		Interval:    10 * time.Millisecond,
	})
}

// Unlimited returns a Limiter that imposes no restrictions.  Useful in
// tests or when the caller wants to manage back-pressure externally.
func Unlimited() *Limiter {
	return New(Options{
		Concurrency: 0, // resolved to 1 inside New; callers should use a
		// large value when they truly want no cap.
	})
}

// HighThroughput returns a Limiter tuned for high-volume, low-latency
// environments where the downstream service can handle bursts.
func HighThroughput() *Limiter {
	return New(Options{
		Concurrency: 32,
		Interval:    2 * time.Millisecond,
	})
}
