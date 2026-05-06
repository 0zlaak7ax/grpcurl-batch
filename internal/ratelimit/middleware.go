package ratelimit

import (
	"context"

	"github.com/your-org/grpcurl-batch/internal/middleware"
)

// Middleware wraps a middleware.Handler so that each invocation must
// first Acquire a slot from the Limiter.  The slot is Released once
// the inner handler returns, regardless of outcome.
//
// This allows a Limiter to be composed into a middleware.Chain alongside
// logging, timeout, and other cross-cutting concerns.
//
//	chain := middleware.Chain(
//		ratelimit.Middleware(lim),
//		middleware.NewLogging(log.Default()),
//		middleware.NewTimeout(5*time.Second),
//	)
func Middleware(lim *Limiter) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if err := lim.Acquire(ctx); err != nil {
				return nil, err
			}
			defer lim.Release()
			return next(ctx, req)
		}
	}
}
