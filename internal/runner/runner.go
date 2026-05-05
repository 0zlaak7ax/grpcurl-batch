// Package runner orchestrates the execution of batched gRPC requests,
// applying retries and rate limiting as configured.
package runner

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yourorg/grpcurl-batch/internal/config"
	"github.com/yourorg/grpcurl-batch/internal/executor"
	"github.com/yourorg/grpcurl-batch/internal/formatter"
	"github.com/yourorg/grpcurl-batch/internal/ratelimit"
)

// Runner executes a batch of gRPC requests.
type Runner struct {
	cfg     *config.Config
	exec    executor.Executor
	limiter *ratelimit.Limiter
}

// New constructs a Runner from the given configuration and executor.
func New(cfg *config.Config, exec executor.Executor) *Runner {
	l := ratelimit.New(ratelimit.Config{
		MaxConcurrent: cfg.MaxConcurrent,
		Interval:      cfg.RateInterval,
	})
	return &Runner{cfg: cfg, exec: exec, limiter: l}
}

// Run executes all requests in the config concurrently (bounded by the limiter)
// and returns a slice of results.
func (r *Runner) Run(ctx context.Context) []formatter.Result {
	results := make([]formatter.Result, len(r.cfg.Requests))
	var wg sync.WaitGroup

	for i, req := range r.cfg.Requests {
		wg.Add(1)
		go func(idx int, req config.Request) {
			defer wg.Done()
			if err := r.limiter.Acquire(ctx); err != nil {
				results[idx] = formatter.Result{Name: req.Name, Err: err}
				return
			}
			defer r.limiter.Release()
			results[idx] = r.executeWithRetry(ctx, req)
		}(i, req)
	}

	wg.Wait()
	return results
}

func (r *Runner) executeWithRetry(ctx context.Context, req config.Request) formatter.Result {
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxRetries; attempt++ {
		out, err := r.exec.Execute(ctx, r.cfg, req)
		if err == nil {
			return formatter.Result{Name: req.Name, Output: out, Attempts: attempt}
		}
		lastErr = err
		log.Printf("[%s] attempt %d/%d failed: %v", req.Name, attempt, r.cfg.MaxRetries, err)
		if attempt < r.cfg.MaxRetries {
			select {
			case <-time.After(r.cfg.RetryDelay):
			case <-ctx.Done():
				return formatter.Result{Name: req.Name, Err: ctx.Err(), Attempts: attempt}
			}
		}
	}
	return formatter.Result{Name: req.Name, Err: lastErr, Attempts: r.cfg.MaxRetries}
}
