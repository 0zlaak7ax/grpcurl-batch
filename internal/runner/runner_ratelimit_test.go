package runner_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/grpcurl-batch/internal/config"
	"github.com/yourorg/grpcurl-batch/internal/executor"
	"github.com/yourorg/grpcurl-batch/internal/runner"
)

type countingExecutor struct {
	inflight int64
	peak     int64
	delay    time.Duration
}

func (c *countingExecutor) Execute(ctx context.Context, cfg *config.Config, req config.Request) (string, error) {
	current := atomic.AddInt64(&c.inflight, 1)
	for {
		pk := atomic.LoadInt64(&c.peak)
		if current <= pk || atomic.CompareAndSwapInt64(&c.peak, pk, current) {
			break
		}
	}
	time.Sleep(c.delay)
	atomic.AddInt64(&c.inflight, -1)
	return `{"ok":true}`, nil
}

var _ executor.Executor = (*countingExecutor)(nil)

func TestRunner_RespectsConcurrencyLimit(t *testing.T) {
	const maxConcurrent = 2
	exec := &countingExecutor{delay: 40 * time.Millisecond}

	cfg := &config.Config{
		Address:       "localhost:50051",
		MaxRetries:    1,
		RetryDelay:    0,
		MaxConcurrent: maxConcurrent,
		Requests: []config.Request{
			{Name: "r1", Method: "pkg.Svc/M"},
			{Name: "r2", Method: "pkg.Svc/M"},
			{Name: "r3", Method: "pkg.Svc/M"},
			{Name: "r4", Method: "pkg.Svc/M"},
		},
	}

	r := runner.New(cfg, exec)
	results := r.Run(context.Background())

	if int64(exec.peak) > maxConcurrent {
		t.Errorf("peak inflight %d exceeded max concurrent %d", exec.peak, maxConcurrent)
	}
	for _, res := range results {
		if res.Err != nil {
			t.Errorf("unexpected error for %s: %v", res.Name, res.Err)
		}
	}
}

func TestRunner_ContextCancelledDuringAcquire(t *testing.T) {
	exec := &countingExecutor{delay: 100 * time.Millisecond}

	cfg := &config.Config{
		Address:       "localhost:50051",
		MaxRetries:    1,
		MaxConcurrent: 1,
		Requests: []config.Request{
			{Name: "r1", Method: "pkg.Svc/M"},
			{Name: "r2", Method: "pkg.Svc/M"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	r := runner.New(cfg, exec)
	results := r.Run(ctx)

	errCount := 0
	for _, res := range results {
		if res.Err != nil {
			errCount++
		}
	}
	if errCount == 0 {
		t.Error("expected at least one error due to context cancellation")
	}
}
