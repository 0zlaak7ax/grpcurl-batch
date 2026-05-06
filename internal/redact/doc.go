// Package redact scrubs sensitive field values from JSON payloads produced
// by grpcurl responses before they are written to logs, formatters, or
// external reporters.
//
// Usage:
//
//	r := redact.New([]string{"password", "token", "secret"})
//	clean := r.Scrub(rawJSON)
//
// Field matching is case-insensitive and operates recursively through
// nested objects and arrays. Non-JSON strings are passed through unchanged.
package redact
