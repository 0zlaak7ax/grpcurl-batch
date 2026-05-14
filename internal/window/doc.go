// Package window implements a sliding-window counter suitable for
// rate-measurement and burst detection in the grpcurl-batch pipeline.
//
// # Overview
//
// A Window divides a fixed-duration interval into a configurable number of
// equal-width buckets. Each call to Add records events into the bucket that
// corresponds to the current moment. Buckets whose leading edge falls outside
// the rolling interval are evicted lazily on the next Add or Count call,
// keeping memory usage constant regardless of event volume.
//
// # Usage
//
//	w := window.New(10*time.Second, 10) // 10 buckets of 1 s each
//	w.Add(1)                            // record an event now
//	fmt.Println(w.Count())             // events in the last 10 s
//
// Window is safe for concurrent use by multiple goroutines.
package window
