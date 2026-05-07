// Package throttle implements a token-bucket throttler for controlling
// the rate of outgoing gRPC requests in grpcurl-batch.
//
// Basic usage:
//
//	t := throttle.New(throttle.Config{Rate: 10, Burst: 20})
//	if err := t.Wait(ctx); err != nil {
//	    // context cancelled or deadline exceeded
//	}
//
// Preset configurations are available via throttle.Preset:
//
//	t := throttle.Preset("medium") // 50 req/s, burst 100
package throttle
