// Package backoff provides pluggable backoff strategies used when scheduling
// retry delays between gRPC call attempts.
//
// Three strategies are available:
//
//   - Exponential: doubles the delay each attempt (with optional jitter)
//     and caps at a configurable maximum.
//   - Fixed: returns a constant delay on every attempt.
//   - Linear: increases the delay linearly with each attempt.
//
// All strategies implement the Strategy interface so they can be swapped
// without changing the retry logic.
//
// Example:
//
//	strategy := backoff.New(100*time.Millisecond, 5*time.Second, true)
//	for attempt := 0; attempt < maxRetries; attempt++ {
//		time.Sleep(strategy.Delay(attempt))
//	}
package backoff
