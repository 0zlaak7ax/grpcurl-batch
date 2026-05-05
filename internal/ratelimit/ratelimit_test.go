package ratelimit_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/grpcurl-batch/internal/ratelimit"
)

func TestNew_DefaultsConcurrency(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxConcurrent: 0})
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l.Release()
}

func TestAcquireRelease_Concurrency(t *testing.T) {
	const max = 3
	l := ratelimit.New(ratelimit.Config{MaxConcurrent: max})

	var inflight int64
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(context.Background()); err != nil {
				t.Errorf("acquire error: %v", err)
				return
			}
			current := atomic.AddInt64(&inflight, 1)
			if current > max {
				t.Errorf("inflight %d exceeded max %d", current, max)
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&inflight, -1)
			l.Release()
		}()
	}
	wg.Wait()
}

func TestAcquire_ContextCancelled(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxConcurrent: 1})
	ctx := context.Background()

	// Exhaust the single token.
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	cancel_ctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	err := l.Acquire(cancel_ctx)
	if err == nil {
		t.Fatal("expected error when context cancelled")
	}
}

func TestAcquire_IntervalThrottles(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{
		MaxConcurrent: 2,
		Interval:      30 * time.Millisecond,
	})

	start := time.Now()
	for i := 0; i < 2; i++ {
		if err := l.Acquire(context.Background()); err != nil {
			t.Fatalf("acquire %d: %v", i, err)
		}
		l.Release()
	}
	elapsed := time.Since(start)
	if elapsed < 30*time.Millisecond {
		t.Errorf("expected at least 30ms elapsed, got %v", elapsed)
	}
}
