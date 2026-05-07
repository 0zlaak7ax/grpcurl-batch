package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/grpcurl-batch/internal/throttle"
)

func TestNew_DefaultBurst_EqualToRate(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 10})
	if got := th.Available(); got != 10 {
		t.Fatalf("expected burst=10, got %v", got)
	}
}

func TestWait_ImmediateWhenTokensAvailable(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 100, Burst: 100})
	ctx := context.Background()
	start := time.Now()
	if err := th.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("expected immediate return, took %v", elapsed)
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	// Rate of 0.001 means effectively no tokens available quickly.
	th := throttle.New(throttle.Config{Rate: 0.001, Burst: 0.001})
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	if err := th.Wait(ctx); err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestWait_ContextDeadline(t *testing.T) {
	th := throttle.New(throttle.Config{Rate: 0.001, Burst: 0})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := th.Wait(ctx)
	if err == nil {
		t.Fatal("expected deadline error")
	}
	if elapsed := time.Since(start); elapsed > 300*time.Millisecond {
		t.Fatalf("took too long to respect deadline: %v", elapsed)
	}
}

func TestPreset_Unlimited_ReturnsNonNil(t *testing.T) {
	th := throttle.Preset("unlimited")
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}

func TestPreset_Unknown_FallsBackToUnlimited(t *testing.T) {
	th := throttle.Preset("nonexistent")
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
	if got := th.Available(); got <= 0 {
		t.Fatalf("expected positive token count, got %v", got)
	}
}

func TestPreset_Low_Rate(t *testing.T) {
	th := throttle.Preset("low")
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
	// Burst for "low" is 10; drain all tokens then Available should be <1.
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		if err := th.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
	}
	if got := th.Available(); got >= 1 {
		t.Fatalf("expected <1 token after draining burst, got %v", got)
	}
}
