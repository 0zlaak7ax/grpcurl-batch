package sampling_test

import (
	"context"
	"testing"

	"github.com/user/grpcurl-batch/internal/sampling"
)

func TestNew_RateZero_NeverSamples(t *testing.T) {
	s := sampling.New(sampling.Config{Rate: 0})
	for i := 0; i < 100; i++ {
		if s.Sample() {
			t.Fatal("expected no samples at rate 0")
		}
	}
}

func TestNew_RateOne_AlwaysSamples(t *testing.T) {
	s := sampling.New(sampling.Config{Rate: 1})
	for i := 0; i < 100; i++ {
		if !s.Sample() {
			t.Fatal("expected every call sampled at rate 1")
		}
	}
}

func TestNew_RateClamped_BelowZero(t *testing.T) {
	s := sampling.New(sampling.Config{Rate: -5})
	for i := 0; i < 50; i++ {
		if s.Sample() {
			t.Fatal("negative rate should clamp to 0")
		}
	}
}

func TestNew_RateClamped_AboveOne(t *testing.T) {
	s := sampling.New(sampling.Config{Rate: 99})
	for i := 0; i < 50; i++ {
		if !s.Sample() {
			t.Fatal("rate > 1 should clamp to 1")
		}
	}
}

func TestPreset_Off(t *testing.T) {
	s := sampling.Preset("off")
	if s.Sample() {
		t.Fatal("off preset should never sample")
	}
}

func TestPreset_Full(t *testing.T) {
	s := sampling.Preset("full")
	if !s.Sample() {
		t.Fatal("full preset should always sample")
	}
}

func TestPreset_Unknown_FallsBackToStandard(t *testing.T) {
	// Standard is 1%; just verify it returns a non-nil sampler.
	s := sampling.Preset("unknown-preset")
	if s == nil {
		t.Fatal("expected non-nil sampler for unknown preset")
	}
}

func TestMiddleware_SetsContextValue(t *testing.T) {
	// Use rate=1 so every request is sampled.
	s := sampling.New(sampling.Config{Rate: 1})
	mw := sampling.Middleware(s)

	var sampled bool
	handler := func(ctx context.Context) error {
		sampled = sampling.IsSampled(ctx)
		return nil
	}

	if err := mw(handler)(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sampled {
		t.Fatal("expected context to be marked as sampled")
	}
}

func TestIsSampled_MissingKey_ReturnsFalse(t *testing.T) {
	if sampling.IsSampled(context.Background()) {
		t.Fatal("expected false when no sampling decision in context")
	}
}
