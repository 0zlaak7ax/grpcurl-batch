// Package correlation provides request correlation ID generation and propagation
// for tracing batched gRPC requests across retries and concurrent executions.
package correlation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type contextKey struct{}

// ID represents a correlation identifier attached to a request.
type ID string

// String returns the string representation of the correlation ID.
func (id ID) String() string { return string(id) }

// Generator creates new correlation IDs.
type Generator struct {
	prefix string
}

// New returns a Generator. If prefix is empty, IDs are generated without one.
func New(prefix string) *Generator {
	return &Generator{prefix: prefix}
}

// Generate produces a new random correlation ID.
func (g *Generator) Generate() (ID, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("correlation: generate random bytes: %w", err)
	}
	hex := hex.EncodeToString(b)
	if g.prefix != "" {
		return ID(g.prefix + "-" + hex), nil
	}
	return ID(hex), nil
}

// WithID attaches a correlation ID to the provided context.
func WithID(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext retrieves the correlation ID from ctx.
// Returns an empty ID and false if none is set.
func FromContext(ctx context.Context) (ID, bool) {
	v, ok := ctx.Value(contextKey{}).(ID)
	return v, ok && v != ""
}

// MustFromContext retrieves the correlation ID or returns a placeholder.
func MustFromContext(ctx context.Context) ID {
	if id, ok := FromContext(ctx); ok {
		return id
	}
	return ID("unknown")
}
