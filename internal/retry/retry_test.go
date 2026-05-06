package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/retry"
)

func fastPolicy(attempts int) retry.Policy {
	return retry.Policy{
		MaxAttempts:  attempts,
		InitialDelay: time.Millisecond,
		Multiplier:   1,
		MaxDelay:     10 * time.Millisecond,
	}
}

func TestDo_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), func(attempt int) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_SuccessAfterRetry(t *testing.T) {
	var calls int32
	err := retry.Do(context.Background(), fastPolicy(3), func(attempt int) error {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AllAttemptsFail(t *testing.T) {
	sentinel := errors.New("boom")
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), func(attempt int) error {
		calls++
		return sentinel
	})
	if !errors.Is(err, retry.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel wrapped in error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelledBetweenAttempts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := retry.Do(ctx, fastPolicy(5), func(attempt int) error {
		calls++
		if calls == 2 {
			cancel()
		}
		return errors.New("fail")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDefaultPolicy(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("MaxAttempts = %d, want 3", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("Multiplier = %f, want 2.0", p.Multiplier)
	}
}

func TestDo_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 0, InitialDelay: time.Millisecond, Multiplier: 1}
	_ = retry.Do(context.Background(), p, func(attempt int) error {
		calls++
		return errors.New("fail")
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
