package notify_test

import (
	"context"
	"errors"
	"testing"

	"grpcurl-batch/internal/notify"
)

type stubNotifier struct {
	called bool
	err    error
}

func (s *stubNotifier) Notify(_ context.Context, _ notify.Summary) error {
	s.called = true
	return s.err
}

func TestMulti_Empty(t *testing.T) {
	m := notify.NewMulti()
	if m.Len() != 0 {
		t.Fatalf("expected 0 notifiers, got %d", m.Len())
	}
	if err := m.Notify(context.Background(), exampleSummary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMulti_AllCalled(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	m := notify.NewMulti(a, b)
	if err := m.Notify(context.Background(), exampleSummary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("not all notifiers were called")
	}
}

func TestMulti_CollectsErrors(t *testing.T) {
	errA := errors.New("backend A failed")
	errB := errors.New("backend B failed")
	a := &stubNotifier{err: errA}
	b := &stubNotifier{err: errB}
	m := notify.NewMulti(a, b)
	err := m.Notify(context.Background(), exampleSummary)
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !errors.Is(err, errA) {
		t.Errorf("expected errA in combined error, got: %v", err)
	}
	if !errors.Is(err, errB) {
		t.Errorf("expected errB in combined error, got: %v", err)
	}
}

func TestMulti_Add(t *testing.T) {
	m := notify.NewMulti()
	m.Add(&stubNotifier{})
	m.Add(&stubNotifier{})
	if m.Len() != 2 {
		t.Fatalf("expected 2, got %d", m.Len())
	}
}

func TestMulti_PartialError_StillCallsAll(t *testing.T) {
	a := &stubNotifier{err: errors.New("oops")}
	b := &stubNotifier{}
	m := notify.NewMulti(a, b)
	_ = m.Notify(context.Background(), exampleSummary)
	if !b.called {
		t.Error("second notifier should be called even when first errors")
	}
}
