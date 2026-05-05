package middleware

import (
	"context"
	"fmt"
	"time"
)

// NewTimeout returns a Middleware that enforces a per-request deadline.
// If d is zero or negative the middleware is a no-op pass-through.
func NewTimeout(d time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			if d <= 0 {
				return next(ctx, req)
			}

			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()

			type result struct {
				resp *Response
				err  error
			}
			ch := make(chan result, 1)
			go func() {
				resp, err := next(ctx, req)
				ch <- result{resp, err}
			}()

			select {
			case r := <-ch:
				return r.resp, r.err
			case <-ctx.Done():
				return nil, fmt.Errorf("request timed out after %s: %w", d, ctx.Err())
			}
		}
	}
}
