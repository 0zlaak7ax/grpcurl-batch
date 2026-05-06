package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grpcurl-batch/internal/timeout"
)

func TestDefaultConfig(t *testing.T) {
	cfg := timeout.DefaultConfig()
	if cfg.PerRequest != 30*time.Second {
		t.Errorf("expected PerRequest=30s, got %v", cfg.PerRequest)
	}
	if cfg.Total != 5*time.Minute {
		t.Errorf("expected Total=5m, got %v", cfg.Total)
	}
}

func TestNew_ZeroValuesUseDefaults(t *testing.T) {
	l := timeout.New(timeout.Config{})
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
	// Verify defaults applied by checking that contexts are created without panic.
	ctx, cancel := l.WithRequest(context.Background())
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestNew_CustomValues(t *testing.T) {
	l := timeout.New(timeout.Config{
		PerRequest: 100 * time.Millisecond,
		Total:      200 * time.Millisecond,
	})
	ctx, cancel := l.WithRequest(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 100*time.Millisecond {
		t.Errorf("deadline too far in future: %v", time.Until(deadline))
	}
}

func TestWithTotal_DeadlineSet(t *testing.T) {
	l := timeout.New(timeout.Config{Total: 500 * time.Millisecond})
	ctx, cancel := l.WithTotal(context.Background())
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected total deadline to be set")
	}
}

func TestWithRequest_ExpiresAfterDuration(t *testing.T) {
	l := timeout.New(timeout.Config{PerRequest: 50 * time.Millisecond})
	ctx, cancel := l.WithRequest(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context should have expired")
	}
}

func TestWrapError_DeadlineExceeded(t *testing.T) {
	err := timeout.WrapError("my-op", context.DeadlineExceeded)
	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Errorf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestWrapError_OtherError(t *testing.T) {
	original := errors.New("some other error")
	err := timeout.WrapError("my-op", original)
	if !errors.Is(err, original) {
		t.Errorf("expected original error to be preserved, got %v", err)
	}
}

func TestWrapError_NilPassthrough(t *testing.T) {
	if timeout.WrapError("op", nil) != nil {
		t.Error("expected nil to be returned unchanged")
	}
}
