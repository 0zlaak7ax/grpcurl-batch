package sampling

import (
	"context"

	"github.com/user/grpcurl-batch/internal/middleware"
)

type contextKey struct{}

// Middleware returns a middleware.Middleware that attaches a sampling decision
// to the context. Downstream handlers can inspect it with IsSampled.
func Middleware(s Sampler) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context) error {
			ctx = context.WithValue(ctx, contextKey{}, s.Sample())
			return next(ctx)
		}
	}
}

// IsSampled reports whether the context was marked as sampled by Middleware.
// Returns false if no sampling decision is present.
func IsSampled(ctx context.Context) bool {
	v, ok := ctx.Value(contextKey{}).(bool)
	return ok && v
}
