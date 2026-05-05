package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"grpcurl-batch/internal/middleware"
)

func nopHandler(_ context.Context, req *middleware.Request) (*middleware.Response, error) {
	return &middleware.Response{Output: "ok", Duration: 1}, nil
}

func errHandler(_ context.Context, _ *middleware.Request) (*middleware.Response, error) {
	return nil, errors.New("boom")
}

func TestChain_Order(t *testing.T) {
	var order []string
	mk := func(name string) middleware.Middleware {
		return func(next middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req *middleware.Request) (*middleware.Response, error) {
				order = append(order, name+"-in")
				r, err := next(ctx, req)
				order = append(order, name+"-out")
				return r, err
			}
		}
	}

	h := middleware.Apply(nopHandler, mk("A"), mk("B"))
	h(context.Background(), &middleware.Request{Method: "Test"})

	want := []string{"A-in", "B-in", "B-out", "A-out"}
	for i, v := range want {
		if order[i] != v {
			t.Errorf("order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestLogging_WritesLines(t *testing.T) {
	var buf bytes.Buffer
	h := middleware.Apply(nopHandler, middleware.NewLogging(&buf))
	h(context.Background(), &middleware.Request{Address: "localhost:50051", Method: "pkg.Svc/Method"})

	out := buf.String()
	if !strings.Contains(out, "-->") || !strings.Contains(out, "<--") {
		t.Errorf("expected arrow markers in log output, got: %s", out)
	}
}

func TestLogging_LogsError(t *testing.T) {
	var buf bytes.Buffer
	h := middleware.Apply(errHandler, middleware.NewLogging(&buf))
	h(context.Background(), &middleware.Request{Method: "Fail"})

	if !strings.Contains(buf.String(), "ERR") {
		t.Errorf("expected ERR in log output, got: %s", buf.String())
	}
}

func TestTimeout_PassThrough(t *testing.T) {
	h := middleware.Apply(nopHandler, middleware.NewTimeout(0))
	resp, err := h(context.Background(), &middleware.Request{})
	if err != nil || resp == nil {
		t.Fatalf("unexpected err=%v resp=%v", err, resp)
	}
}

func TestTimeout_Exceeded(t *testing.T) {
	slow := func(_ context.Context, _ *middleware.Request) (*middleware.Response, error) {
		time.Sleep(200 * time.Millisecond)
		return &middleware.Response{Output: "late"}, nil
	}
	h := middleware.Apply(slow, middleware.NewTimeout(20*time.Millisecond))
	_, err := h(context.Background(), &middleware.Request{})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApply_NoMiddlewares(t *testing.T) {
	h := middleware.Apply(nopHandler)
	resp, err := h(context.Background(), &middleware.Request{})
	if err != nil || resp.Output != "ok" {
		t.Fatalf("unexpected result: resp=%v err=%v", resp, err)
	}
}
