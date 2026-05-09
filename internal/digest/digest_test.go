package digest_test

import (
	"testing"

	"github.com/user/grpcurl-batch/internal/digest"
)

func baseRequest() digest.Request {
	return digest.Request{
		Address: "localhost:50051",
		Method:  "pkg.Service/Method",
		Headers: map[string]string{"authorization": "Bearer token"},
		Body:    `{"key":"value"}`,
	}
}

func TestCompute_Deterministic(t *testing.T) {
	r := baseRequest()
	d1 := digest.Compute(r)
	d2 := digest.Compute(r)
	if d1 != d2 {
		t.Errorf("expected same digest, got %q and %q", d1, d2)
	}
}

func TestCompute_NonEmpty(t *testing.T) {
	d := digest.Compute(baseRequest())
	if d == "" {
		t.Error("expected non-empty digest")
	}
	if len(d) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(d))
	}
}

func TestCompute_DifferentMethod(t *testing.T) {
	r1 := baseRequest()
	r2 := baseRequest()
	r2.Method = "pkg.Service/Other"
	if digest.Compute(r1) == digest.Compute(r2) {
		t.Error("expected different digests for different methods")
	}
}

func TestCompute_HeaderOrderIndependent(t *testing.T) {
	r1 := baseRequest()
	r1.Headers = map[string]string{"a": "1", "b": "2"}

	r2 := baseRequest()
	r2.Headers = map[string]string{"b": "2", "a": "1"}

	if digest.Compute(r1) != digest.Compute(r2) {
		t.Error("expected same digest regardless of header insertion order")
	}
}

func TestCompute_BodyWhitespaceNormalised(t *testing.T) {
	r1 := baseRequest()
	r1.Body = `{"key":"value"}`

	r2 := baseRequest()
	r2.Body = `{ "key" :  "value" }`

	if digest.Compute(r1) != digest.Compute(r2) {
		t.Error("expected same digest for semantically equal JSON bodies")
	}
}

func TestCompute_NonJSONBodyUsedRaw(t *testing.T) {
	r1 := baseRequest()
	r1.Body = "plain text body"

	r2 := baseRequest()
	r2.Body = "plain text body"

	if digest.Compute(r1) != digest.Compute(r2) {
		t.Error("expected same digest for identical non-JSON bodies")
	}
}

func TestCompute_EmptyBodyAndNoHeaders(t *testing.T) {
	r := digest.Request{
		Address: "host:443",
		Method:  "svc/M",
	}
	d := digest.Compute(r)
	if d == "" {
		t.Error("expected non-empty digest for minimal request")
	}
}
