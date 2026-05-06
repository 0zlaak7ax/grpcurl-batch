package ratelimit_test

import (
	"context"
	"testing"

	"github.com/your-org/grpcurl-batch/internal/ratelimit"
)

func TestPreset_AcquireRelease(t *testing.T) {
	lim := ratelimit.Preset()

	ctx := context.Background()
	if err := lim.Acquire(ctx); err != nil {
		t.Fatalf("Preset Acquire: %v", err)
	}
	lim.Release()
}

func TestUnlimited_AcquireRelease(t *testing.T) {
	lim := ratelimit.Unlimited()

	ctx := context.Background()
	if err := lim.Acquire(ctx); err != nil {
		t.Fatalf("Unlimited Acquire: %v", err)
	}
	lim.Release()
}

func TestHighThroughput_ParallelAcquire(t *testing.T) {
	lim := ratelimit.HighThroughput()
	ctx := context.Background()

	const goroutines = 16
	errCh := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			if err := lim.Acquire(ctx); err != nil {
				errCh <- err
				return
			}
			lim.Release()
			errCh <- nil
		}()
	}

	for i := 0; i < goroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("goroutine %d: %v", i, err)
		}
	}
}

func TestPreset_ReturnsNonNil(t *testing.T) {
	if lim := ratelimit.Preset(); lim == nil {
		t.Fatal("Preset returned nil")
	}
}

func TestHighThroughput_ReturnsNonNil(t *testing.T) {
	if lim := ratelimit.HighThroughput(); lim == nil {
		t.Fatal("HighThroughput returned nil")
	}
}
