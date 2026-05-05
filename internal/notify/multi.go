package notify

import (
	"context"
	"errors"
)

// Multi fans a single Notify call out to several backends.
// All notifiers are attempted; errors are joined and returned together.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi notifier wrapping the provided backends.
func NewMulti(nn ...Notifier) *Multi {
	return &Multi{notifiers: nn}
}

// Add appends a notifier to the fan-out list.
func (m *Multi) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Notify calls every registered notifier and collects errors.
func (m *Multi) Notify(ctx context.Context, s Summary) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(ctx, s); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Len returns the number of registered notifiers.
func (m *Multi) Len() int { return len(m.notifiers) }
