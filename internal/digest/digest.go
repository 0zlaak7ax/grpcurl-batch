// Package digest provides request fingerprinting for grpcurl-batch.
// It computes a stable hash over a request's method, address, and body
// so that identical requests can be identified across runs.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Request holds the fields used to compute a digest.
type Request struct {
	Address string
	Method  string
	Headers map[string]string
	Body    string
}

// Compute returns a deterministic SHA-256 hex digest for the given Request.
// Header keys are sorted before hashing so that insertion order does not
// affect the result.
func Compute(r Request) string {
	h := sha256.New()

	fmt.Fprintf(h, "address:%s\n", strings.TrimSpace(r.Address))
	fmt.Fprintf(h, "method:%s\n", strings.TrimSpace(r.Method))

	keys := make([]string, 0, len(r.Headers))
	for k := range r.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(h, "header:%s=%s\n", strings.ToLower(k), r.Headers[k])
	}

	normBody := normaliseBody(r.Body)
	fmt.Fprintf(h, "body:%s\n", normBody)

	return hex.EncodeToString(h.Sum(nil))
}

// normaliseBody attempts to marshal the body as JSON so that whitespace
// differences do not affect the digest. If the body is not valid JSON the
// raw trimmed string is used instead.
func normaliseBody(body string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}
	var v interface{}
	if err := json.Unmarshal([]byte(trimmed), &v); err != nil {
		return trimmed
	}
	out, err := json.Marshal(v)
	if err != nil {
		return trimmed
	}
	return string(out)
}
