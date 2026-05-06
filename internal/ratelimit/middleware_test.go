package ratelimit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/grpcurl-batch/internal/middleware"
	"github.com/your-org/grpcurl-batch/internal/ratelimit"
)

func TestMiddleware_PassesThrough(t *testing.T) {
	lim := ratelimit.New(ratelimit.Options{Concurrency: 2})
	mw := ratelimit.Middleware(lim)

	handler := mw(func(_ context.Context, req any) (any, error) {
		return "ok", nil
	})

	res, err := handler(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "ok" {
		t.Fatalf("expected 'ok', got %v", res)
	}
}

func TestMiddleware_PropagatesHandlerError(t *testing.T) {
	lim := ratelimit.New(ratelimit.Options{Concurrency: 2})
	mw := ratelimit.Middleware(lim)
	want := errors.New("handler error")

	handler := mw(func(_ context.Context, _ any) (any, error) {
		return nil, want
	})

	_, err := handler(context.Background(), nil)
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestMiddleware_ContextCancelledBeforeAcquire(t *testing.T) {
	lim := ratelimit.New(ratelimit.Options{Concurrency: 1})

	// Fill the single slot so the next Acquire will block.
	ctxFill := context.Background()
	if err := lim.Acquire(ctxFill); err != nil {
		t.Fatalf("fill Acquire: %v", err)
	}
	defer lim.Release()

	mw := ratelimit.Middleware(lim)
	handler := mw(func(_ context.Context, _ any) (any, error) {
		return "should not reach", nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := handler(ctx, nil)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestMiddleware_ComposesWithChain(t *testing.T) {
	lim := ratelimit.New(ratelimit.Options{Concurrency: 4})

	called := false
	chain := middleware.Chain(
		ratelimit.Middleware(lim),
	)
	final := middleware.Apply(chain, func(_ context.Context, _ any) (any, error) {
		called = true
		return nil, nil
	})

	if _, err := final(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("inner handler was not called")
	}
}
