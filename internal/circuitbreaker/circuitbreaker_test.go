package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/your-org/grpcurl-batch/internal/circuitbreaker"
)

func TestNew_Defaults(t *testing.T) {
	b := circuitbreaker.New(0, 0)
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected StateClosed, got %v", b.State())
	}
}

func TestAllow_ClosedAlwaysAllows(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	for i := 0; i < 10; i++ {
		if !b.Allow() {
			t.Errorf("iteration %d: expected Allow()=true on closed breaker", i)
		}
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected still closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Errorf("expected StateOpen after 3 failures, got %v", b.State())
	}
}

func TestAllow_OpenRejectsBeforeCooldown(t *testing.T) {
	b := circuitbreaker.New(1, time.Minute)
	b.RecordFailure()
	if b.Allow() {
		t.Error("expected Allow()=false while open")
	}
}

func TestAllow_TransitionsToHalfOpenAfterCooldown(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if !b.Allow() {
		t.Error("expected Allow()=true after cooldown (half-open)")
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestRecordSuccess_ResetsToClosed(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	b.Allow() // move to half-open
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected StateClosed after success, got %v", b.State())
	}
}

func TestRecordSuccess_ResetsFailureCount(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected closed: success should reset counter, got %v", b.State())
	}
}
