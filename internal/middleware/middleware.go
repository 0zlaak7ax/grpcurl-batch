// Package middleware provides a chain of request/response interceptors
// that can be applied to gRPC batch executions before and after each call.
package middleware

import "context"

// Request holds the data for a single gRPC call before execution.
type Request struct {
	Method  string
	Address string
	Body    string
	Headers map[string]string
}

// Response holds the result of a single gRPC call after execution.
type Response struct {
	Output   string
	Err      error
	Duration int64 // milliseconds
}

// Handler is a function that processes a request and returns a response.
type Handler func(ctx context.Context, req *Request) (*Response, error)

// Middleware wraps a Handler to add pre/post processing.
type Middleware func(next Handler) Handler

// Chain composes multiple middlewares into a single Handler wrapper.
// Middlewares are applied in order: first middleware is outermost.
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply wraps the given handler with the provided middleware chain.
func Apply(handler Handler, middlewares ...Middleware) Handler {
	if len(middlewares) == 0 {
		return handler
	}
	return Chain(middlewares...)(handler)
}
