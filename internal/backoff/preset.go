package backoff

import "time"

// DefaultExponential returns the standard exponential backoff used by the
// runner: starts at 200 ms, caps at 30 s, jitter enabled.
func DefaultExponential() *Exponential {
	return New(200*time.Millisecond, 30*time.Second, true)
}

// AggressiveExponential returns a faster-starting exponential backoff suitable
// for low-latency environments: starts at 50 ms, caps at 5 s, no jitter.
func AggressiveExponential() *Exponential {
	return New(50*time.Millisecond, 5*time.Second, false)
}

// DefaultFixed returns a fixed 1-second delay between attempts.
func DefaultFixed() *Fixed {
	return &Fixed{Interval: time.Second}
}

// DefaultLinear returns a linear backoff starting at 500 ms, capped at 10 s.
func DefaultLinear() *Linear {
	return &Linear{
		Base: 500 * time.Millisecond,
		Max:  10 * time.Second,
	}
}

// ConservativeExponential returns a slow-starting exponential backoff suitable
// for rate-limited or quota-sensitive environments: starts at 1 s, caps at
// 5 min, jitter enabled to spread retries across multiple clients.
func ConservativeExponential() *Exponential {
	return New(time.Second, 5*time.Minute, true)
}
