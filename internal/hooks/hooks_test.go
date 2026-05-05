package hooks_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/grpcurl-batch/internal/hooks"
)

func newRegistry() *hooks.Registry {
	return hooks.New(log.New(os.Stderr, "", 0))
}

func TestRegistry_Empty(t *testing.T) {
	r := newRegistry()
	p := hooks.Payload{Event: hooks.EventBeforeRequest, Method: "pkg.Svc/Method"}
	if err := r.Fire(context.Background(), p); err != nil {
		t.Fatalf("unexpected error on empty registry: %v", err)
	}
}

func TestRegistry_Len(t *testing.T) {
	r := newRegistry()
	r.Register(hooks.EventBeforeRequest, func(_ context.Context, _ hooks.Payload) error { return nil })
	r.Register(hooks.EventBeforeRequest, func(_ context.Context, _ hooks.Payload) error { return nil })
	if got := r.Len(hooks.EventBeforeRequest); got != 2 {
		t.Fatalf("expected 2 hooks, got %d", got)
	}
	if got := r.Len(hooks.EventAfterRequest); got != 0 {
		t.Fatalf("expected 0 hooks for after_request, got %d", got)
	}
}

func TestFire_HooksCalledInOrder(t *testing.T) {
	r := newRegistry()
	var order []int
	for i := 0; i < 3; i++ {
		i := i
		r.Register(hooks.EventAfterRequest, func(_ context.Context, _ hooks.Payload) error {
			order = append(order, i)
			return nil
		})
	}
	r.Fire(context.Background(), hooks.Payload{Event: hooks.EventAfterRequest})
	for idx, v := range order {
		if v != idx {
			t.Fatalf("expected order %d, got %d", idx, v)
		}
	}
}

func TestFire_StopsOnError(t *testing.T) {
	r := newRegistry()
	called := 0
	r.Register(hooks.EventOnFailure, func(_ context.Context, _ hooks.Payload) error {
		called++
		return errors.New("hook failed")
	})
	r.Register(hooks.EventOnFailure, func(_ context.Context, _ hooks.Payload) error {
		called++
		return nil
	})
	err := r.Fire(context.Background(), hooks.Payload{Event: hooks.EventOnFailure})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if called != 1 {
		t.Fatalf("expected only 1 hook called, got %d", called)
	}
}

func TestFire_PayloadFieldsVisible(t *testing.T) {
	r := newRegistry()
	var received hooks.Payload
	r.Register(hooks.EventAfterRequest, func(_ context.Context, p hooks.Payload) error {
		received = p
		return nil
	})
	sent := hooks.Payload{
		Event:   hooks.EventAfterRequest,
		Method:  "pkg.Svc/DoThing",
		Attempt: 2,
		Elapsed: 150 * time.Millisecond,
		Success: true,
		Output:  `{"result":"ok"}`,
	}
	r.Fire(context.Background(), sent)
	if received.Method != sent.Method || received.Attempt != sent.Attempt || received.Elapsed != sent.Elapsed {
		t.Fatalf("payload mismatch: got %+v", received)
	}
}
