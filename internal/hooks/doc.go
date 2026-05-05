// Package hooks implements a lightweight lifecycle hook registry for
// grpcurl-batch. Hooks can be attached to three events:
//
//   - EventBeforeRequest – fired immediately before a gRPC call is dispatched.
//   - EventAfterRequest  – fired after every attempt, regardless of outcome.
//   - EventOnFailure     – fired only when all retry attempts are exhausted.
//
// Usage:
//
//	reg := hooks.New(logger)
//	reg.Register(hooks.EventAfterRequest, func(ctx context.Context, p hooks.Payload) error {
//		fmt.Printf("method=%s success=%v elapsed=%s\n", p.Method, p.Success, p.Elapsed)
//		return nil
//	})
//
// Hooks are executed synchronously in registration order. Returning a non-nil
// error from a hook aborts the remaining hooks for that event.
package hooks
