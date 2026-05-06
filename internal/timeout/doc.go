// Package timeout provides deadline management for individual gRPC requests
// and entire batch runs within grpcurl-batch.
//
// # Overview
//
// A [Limiter] is constructed from a [Config] that specifies two independent
// deadlines:
//
//   - PerRequest: maximum wall-clock time allowed for a single gRPC call.
//   - Total: maximum wall-clock time allowed for the whole batch.
//
// # Usage
//
//	lim := timeout.New(timeout.Config{
//		PerRequest: 10 * time.Second,
//		Total:      2 * time.Minute,
//	})
//
//	batchCtx, cancelBatch := lim.WithTotal(ctx)
//	defer cancelBatch()
//
//	for _, req := range requests {
//		reqCtx, cancelReq := lim.WithRequest(batchCtx)
//		err := execute(reqCtx, req)
//		cancelReq()
//		if err != nil {
//			log.Println(timeout.WrapError(req.Name, err))
//		}
//	}
//
// Named presets are available via [Preset].
package timeout
