package runner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/grpcurl-batch/internal/config"
	"github.com/user/grpcurl-batch/internal/runner"
)

type mockExecutor struct {
	callCount int
	failTimes int
	output    string
}

func (m *mockExecutor) Execute(_ context.Context, _ config.Request) (string, error) {
	m.callCount++
	if m.callCount <= m.failTimes {
		return "", errors.New("mock error")
	}
	return m.output, nil
}

func baseCfg() *config.Config {
	return &config.Config{
		Retry: config.RetryConfig{
			MaxAttempts: 3,
			Delay:       1 * time.Millisecond,
		},
		Requests: []config.Request{
			{Name: "test-req", Method: "pkg.Svc/Method"},
		},
	}
}

func TestRunner_SuccessFirstAttempt(t *testing.T) {
	exec := &mockExecutor{output: `{"ok":true}`}
	r := runner.New(baseCfg(), exec)
	results := r.Run(context.Background())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Success {
		t.Errorf("expected success")
	}
	if results[0].Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", results[0].Attempts)
	}
}

func TestRunner_SuccessAfterRetry(t *testing.T) {
	exec := &mockExecutor{failTimes: 2, output: "ok"}
	r := runner.New(baseCfg(), exec)
	results := r.Run(context.Background())
	if !results[0].Success {
		t.Errorf("expected success after retry")
	}
	if results[0].Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", results[0].Attempts)
	}
}

func TestRunner_AllAttemptsFail(t *testing.T) {
	exec := &mockExecutor{failTimes: 10}
	r := runner.New(baseCfg(), exec)
	results := r.Run(context.Background())
	if results[0].Success {
		t.Errorf("expected failure")
	}
	if results[0].Err == nil {
		t.Errorf("expected non-nil error")
	}
}

func TestRunner_ContextCancelled(t *testing.T) {
	exec := &mockExecutor{failTimes: 10}
	cfg := baseCfg()
	cfg.Retry.Delay = 500 * time.Millisecond
	r := runner.New(cfg, exec)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	results := r.Run(ctx)
	if results[0].Success {
		t.Errorf("expected failure on cancelled context")
	}
}
