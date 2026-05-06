// Package redact provides utilities for scrubbing sensitive field values
// from gRPC request/response payloads before logging or reporting.
package redact

import (
	"encoding/json"
	"strings"
)

const placeholder = "[REDACTED]"

// Redactor scrubs a configured set of field names from JSON payloads.
type Redactor struct {
	fields map[string]struct{}
}

// New returns a Redactor that will mask the given field names.
// Field matching is case-insensitive.
func New(fields []string) *Redactor {
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		m[strings.ToLower(f)] = struct{}{}
	}
	return &Redactor{fields: m}
}

// Scrub parses src as JSON, replaces sensitive field values with
// [REDACTED], and returns the re-encoded JSON. If src is not valid
// JSON the original string is returned unchanged.
func (r *Redactor) Scrub(src string) string {
	if len(r.fields) == 0 {
		return src
	}
	var obj interface{}
	if err := json.Unmarshal([]byte(src), &obj); err != nil {
		return src
	}
	obj = r.walk(obj)
	b, err := json.Marshal(obj)
	if err != nil {
		return src
	}
	return string(b)
}

func (r *Redactor) walk(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, child := range val {
			if _, sensitive := r.fields[strings.ToLower(k)]; sensitive {
				out[k] = placeholder
			} else {
				out[k] = r.walk(child)
			}
		}
		return out
	case []interface{}:
		for i, item := range val {
			val[i] = r.walk(item)
		}
		return val
	}
	return v
}
