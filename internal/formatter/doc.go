// Package formatter provides output formatting for grpcurl-batch results.
//
// It supports three output formats:
//
//   - text    — human-readable line-based output (default)
//   - json    — one JSON object per result, suitable for log ingestion
//   - summary — minimal ✓/✗ per-method output
//
// Usage:
//
//	f := formatter.New(formatter.FormatText, os.Stdout)
//	f.Write(formatter.Result{
//		Method:   "pkg.Service/Method",
//		Success:  true,
//		Attempts: 1,
//		Duration: 42 * time.Millisecond,
//		Output:   `{"id":"abc"}`,
//	})
//
// Results can be accumulated via Collector and a final summary printed with
// PrintSummary.
package formatter
