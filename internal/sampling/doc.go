// Package sampling provides probabilistic request sampling for grpcurl-batch.
//
// A Sampler is created with a rate in [0.0, 1.0] and its Sample method
// returns true for approximately that fraction of calls. The package also
// provides:
//
//   - Preset – named samplers for common scenarios (off, debug, standard, full)
//   - Middleware – a middleware.Middleware that stamps the sampling decision
//     onto the request context so downstream components can branch on it
//   - IsSampled – a helper to read the decision from a context
//
// Example:
//
//	s := sampling.Preset("debug")  // 10% sampling
//	chain := middleware.Chain(sampling.Middleware(s), otherMiddleware)
package sampling
