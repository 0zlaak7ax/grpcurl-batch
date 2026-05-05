package middleware

import (
	"context"
	"fmt"
	"io"
	"time"
)

// NewLogging returns a Middleware that logs each request method and duration.
func NewLogging(w io.Writer) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *Request) (*Response, error) {
			start := time.Now()
			fmt.Fprintf(w, "[middleware/logging] --> %s %s\n", req.Address, req.Method)

			resp, err := next(ctx, req)

			elapsed := time.Since(start).Milliseconds()
			status := "OK"
			if err != nil {
				status = fmt.Sprintf("ERR: %v", err)
			}
			fmt.Fprintf(w, "[middleware/logging] <-- %s %s (%dms) %s\n",
				req.Address, req.Method, elapsed, status)

			return resp, err
		}
	}
}
