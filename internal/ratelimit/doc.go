// Package ratelimit provides concurrency and interval-based throttling
// for outbound gRPC requests.
//
// A Limiter controls two dimensions of rate limiting:
//
//  1. Concurrency – the maximum number of requests that may be in-flight
//     simultaneously, enforced via a semaphore channel.
//
//  2. Interval – a minimum duration between successive Acquire calls,
//     enforced via a time.Ticker.  Set to zero to disable interval
//     throttling.
//
// Typical usage:
//
//	lim := ratelimit.New(ratelimit.Options{
//		Concurrency: 4,
//		Interval:    50 * time.Millisecond,
//	})
//
//	if err := lim.Acquire(ctx); err != nil {
//		return err
//	}
//	defer lim.Release()
package ratelimit
