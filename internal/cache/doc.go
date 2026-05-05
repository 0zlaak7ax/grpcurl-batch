// Package cache implements a lightweight, thread-safe in-memory result cache
// for grpcurl-batch.
//
// When the same gRPC request (identified by address + method + body) appears
// multiple times in a batch file, the cache can short-circuit subsequent
// executions and return the previously recorded output without invoking
// grpcurl again.
//
// Usage:
//
//	c := cache.New(30 * time.Second) // zero disables caching
//	key := cache.Key(address, method, body)
//	if e := c.Get(key); e != nil {
//		// use e.Output / e.Err
//	}
//	// ... execute grpcurl ...
//	c.Set(key, output, err)
//
// The TTL is evaluated lazily on Get; call Purge periodically to reclaim
// memory from expired entries.
package cache
