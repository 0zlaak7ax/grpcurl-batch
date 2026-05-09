// Package audit provides a structured, concurrency-safe audit logger for
// grpcurl-batch. Each gRPC batch event — request dispatch, response receipt,
// retry attempt, or deliberate skip — is serialised as a JSON line and written
// to a configurable io.Writer (file, stdout, or any sink).
//
// Usage:
//
//	f, _ := os.OpenFile("audit.jsonl", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
//	log := audit.New(f)
//	log.RecordRequest(correlationID, "/pkg.Service/Method", nil)
//	log.RecordResponse(correlationID, "/pkg.Service/Method", 1, 42, nil)
//
// Events are newline-delimited JSON (NDJSON) so they can be streamed directly
// into log aggregators such as Loki, Splunk, or jq pipelines.
package audit
