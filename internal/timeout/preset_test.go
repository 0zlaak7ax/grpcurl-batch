package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/grpcurl-batch/internal/timeout"
)

func TestPreset_Default(t *testing.T) {
	l := timeout.Preset("default")
	if l == nil {
		t.Fatal("expected non-nil Limiter for 'default'")
	}
}

func TestPreset_Unknown_FallsBackToDefault(t *testing.T) {
	l := timeout.Preset("unknown-preset")
	if l == nil {
		t.Fatal("expected non-nil Limiter for unknown preset")
	}
}

func TestPreset_Fast_ShorterDeadline(t *testing.T) {
	fast := timeout.Preset("fast")
	def := timeout.Preset("default")

	fastCtx, cancelFast := fast.WithRequest(context.Background())
	defer cancelFast()
	defCtx, cancelDef := def.WithRequest(context.Background())
	defer cancelDef()

	fastDL, _ := fastCtx.Deadline()
	defDL, _ := defCtx.Deadline()

	if !fastDL.Before(defDL) {
		t.Errorf("fast preset deadline should be before default: fast=%v default=%v",
			time.Until(fastDL), time.Until(defDL))
	}
}

func TestPreset_Slow_LongerDeadline(t *testing.T) {
	slow := timeout.Preset("slow")
	def := timeout.Preset("default")

	slowCtx, cancelSlow := slow.WithRequest(context.Background())
	defer cancelSlow()
	defCtx, cancelDef := def.WithRequest(context.Background())
	defer cancelDef()

	slowDL, _ := slowCtx.Deadline()
	defDL, _ := defCtx.Deadline()

	if !slowDL.After(defDL) {
		t.Errorf("slow preset deadline should be after default: slow=%v default=%v",
			time.Until(slowDL), time.Until(defDL))
	}
}

func TestPreset_ReturnsNonNil(t *testing.T) {
	for _, name := range []string{"fast", "default", "slow", ""} {
		if timeout.Preset(name) == nil {
			t.Errorf("Preset(%q) returned nil", name)
		}
	}
}
