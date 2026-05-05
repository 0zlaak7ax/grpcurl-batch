// Package middleware provides composable interceptors for gRPC batch requests.
//
// Middlewares wrap a Handler (a function that executes a single gRPC call) and
// can perform work before and/or after the underlying call — for example
// logging, timeout enforcement, header injection, or metrics recording.
//
// Usage:
//
//	handler := middleware.Apply(
//		myBaseHandler,
//		middleware.NewLogging(os.Stderr),
//		middleware.NewTimeout(5*time.Second),
//	)
//
// Middlewares are applied in declaration order; the first listed middleware is
// the outermost wrapper and therefore runs first on the way in and last on the
// way out.
package middleware
