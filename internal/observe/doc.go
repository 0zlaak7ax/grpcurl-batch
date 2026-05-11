// Package observe provides a lightweight span-style observation layer
// for grpcurl-batch operations.
//
// Usage:
//
//	obs := observe.New(os.Stderr)
//	err := obs.Record(ctx, "grpc.call", map[string]string{"method": m}, func(ctx context.Context) error {
//		return executor.Execute(ctx, req)
//	})
//
// After all calls, inspect results:
//
//	for _, s := range obs.Spans() {
//		fmt.Println(s.Name, s.Duration(), s.OK())
//	}
//
// The Observer is safe for concurrent use. Span output is written to
// the io.Writer supplied to New; pass io.Discard to suppress output.
package observe
