package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/grpcurl-batch/internal/retry"
)

// TestDo_ExponentialBackoff verifies that successive delays grow.
func TestDo_ExponentialBackoff(t *testing.T) {
	policy := retry.Policy{
		MaxAttempts:  4,
		InitialDelay: 10 * time.Millisecond,
		Multiplier:   2.0,
		MaxDelay:     200 * time.Millisecond,
	}

	timestamps := make([]time.Time, 0, 4)
	_ = retry.Do(context.Background(), policy, func(attempt int) error {
		timestamps = append(timestamps, time.Now())
		return errors.New("always fail")
	})

	if len(timestamps) != 4 {
		t.Fatalf("expected 4 timestamps, got %d", len(timestamps))
	}

	d1 := timestamps[1].Sub(timestamps[0])
	d2 := timestamps[2].Sub(timestamps[1])

	// d2 should be meaningfully larger than d1 due to back-off.
	if d2 < d1 {
		t.Errorf("expected d2 (%v) >= d1 (%v) for exponential back-off", d2, d1)
	}
}

// TestDo_MaxDelayCap ensures delays never exceed MaxDelay.
func TestDo_MaxDelayCap(t *testing.T) {
	policy := retry.Policy{
		MaxAttempts:  5,
		InitialDelay: 50 * time.Millisecond,
		Multiplier:   10.0,
		MaxDelay:     60 * time.Millisecond,
	}

	start := time.Now()
	_ = retry.Do(context.Background(), policy, func(attempt int) error {
		return errors.New("fail")
	})
	elapsed := time.Since(start)

	// 4 delays capped at 60ms each → max ~240ms; add generous margin.
	if elapsed > 600*time.Millisecond {
		t.Errorf("elapsed %v exceeds expected cap; MaxDelay not respected", elapsed)
	}
}

// TestDo_ContextDeadline checks that a tight deadline aborts retries.
func TestDo_ContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	policy := retry.Policy{
		MaxAttempts:  10,
		InitialDelay: 20 * time.Millisecond,
		Multiplier:   1,
		MaxDelay:     20 * time.Millisecond,
	}

	err := retry.Do(ctx, policy, func(attempt int) error {
		return errors.New("fail")
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
