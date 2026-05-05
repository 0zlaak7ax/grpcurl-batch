// Package hooks provides pre/post request lifecycle hooks for batch execution.
package hooks

import (
	"context"
	"log"
	"time"
)

// Event represents the lifecycle stage at which a hook fires.
type Event string

const (
	EventBeforeRequest Event = "before_request"
	EventAfterRequest  Event = "after_request"
	EventOnFailure     Event = "on_failure"
)

// Payload carries contextual data passed to each hook.
type Payload struct {
	Event     Event
	Method    string
	Attempt   int
	Elapsed   time.Duration
	Success   bool
	Output    string
	Error     string
}

// HookFunc is a function that receives a Payload and may return an error to
// abort further processing.
type HookFunc func(ctx context.Context, p Payload) error

// Registry holds named hooks keyed by lifecycle Event.
type Registry struct {
	hooks map[Event][]HookFunc
	logger *log.Logger
}

// New returns an initialised Registry.
func New(logger *log.Logger) *Registry {
	return &Registry{
		hooks:  make(map[Event][]HookFunc),
		logger: logger,
	}
}

// Register appends a HookFunc for the given Event.
func (r *Registry) Register(event Event, fn HookFunc) {
	r.hooks[event] = append(r.hooks[event], fn)
}

// Fire executes all hooks registered for the given Event in order.
// If any hook returns an error, execution stops and the error is returned.
func (r *Registry) Fire(ctx context.Context, p Payload) error {
	for _, fn := range r.hooks[p.Event] {
		if err := fn(ctx, p); err != nil {
			r.logger.Printf("hook error [%s] method=%s attempt=%d: %v",
				p.Event, p.Method, p.Attempt, err)
			return err
		}
	}
	return nil
}

// Len returns the number of hooks registered for an Event.
func (r *Registry) Len(event Event) int {
	return len(r.hooks[event])
}
