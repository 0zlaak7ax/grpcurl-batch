package runner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/user/grpcurl-batch/internal/config"
)

// Result holds the outcome of a single gRPC request.
type Result struct {
	Name    string
	Success bool
	Output  string
	Err     error
	Attempts int
}

// Executor defines how a single grpcurl call is made.
type Executor interface {
	Execute(ctx context.Context, req config.Request) (string, error)
}

// Runner executes batched gRPC requests from a config.
type Runner struct {
	cfg      *config.Config
	executor Executor
}

// New creates a Runner with the provided config and executor.
func New(cfg *config.Config, executor Executor) *Runner {
	return &Runner{cfg: cfg, executor: executor}
}

// Run executes all requests defined in the config and returns results.
func (r *Runner) Run(ctx context.Context) []Result {
	results := make([]Result, 0, len(r.cfg.Requests))
	for _, req := range r.cfg.Requests {
		res := r.runWithRetry(ctx, req)
		results = append(results, res)
	}
	return results
}

func (r *Runner) runWithRetry(ctx context.Context, req config.Request) Result {
	maxAttempts := r.cfg.Retry.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	result := Result{Name: req.Name}
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result.Attempts = attempt
		output, err := r.executor.Execute(ctx, req)
		if err == nil {
			result.Success = true
			result.Output = output
			return result
		}
		result.Err = err
		log.Printf("[%s] attempt %d/%d failed: %v", req.Name, attempt, maxAttempts, err)
		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				result.Err = fmt.Errorf("context cancelled: %w", ctx.Err())
				return result
			case <-time.After(r.cfg.Retry.Delay):
			}
		}
	}
	return result
}
