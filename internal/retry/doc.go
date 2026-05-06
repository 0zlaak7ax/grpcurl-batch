// Package retry provides a context-aware retry loop with exponential backoff
// for use throughout grpcurl-batch.
//
// Basic usage:
//
//	p := retry.DefaultPolicy()
//	err := retry.Do(ctx, p, func(attempt int) error {
//		return callSomething()
//	})
//	if errors.Is(err, retry.ErrExhausted) {
//		// all attempts consumed
//	}
//
// The Policy.Multiplier field controls back-off growth:
//   - 1.0 → constant delay
//   - 2.0 → classic exponential back-off
//
// Delays are capped at Policy.MaxDelay to avoid unbounded waits.
package retry
