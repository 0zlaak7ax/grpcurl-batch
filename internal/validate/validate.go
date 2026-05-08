// Package validate provides request payload validation for gRPC batch definitions.
package validate

import (
	"errors"
	"fmt"
	"strings"
)

// Rule is a single validation rule applied to a request field.
type Rule struct {
	Field    string
	Required bool
	MinLen   int
	MaxLen   int
}

// Validator checks request payloads against a set of rules.
type Validator struct {
	rules []Rule
}

// New creates a Validator with the given rules.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// ValidationError holds all field-level errors found during validation.
type ValidationError struct {
	Errs []string
}

func (e *ValidationError) Error() string {
	return "validation failed: " + strings.Join(e.Errs, "; ")
}

// Validate checks the provided fields map against the configured rules.
// fields is a map of field name to string value extracted from the request payload.
func (v *Validator) Validate(fields map[string]string) error {
	var errs []string

	for _, r := range v.rules {
		val, ok := fields[r.Field]

		if r.Required && (!ok || strings.TrimSpace(val) == "") {
			errs = append(errs, fmt.Sprintf("field %q is required", r.Field))
			continue
		}

		if !ok {
			continue
		}

		if r.MinLen > 0 && len(val) < r.MinLen {
			errs = append(errs, fmt.Sprintf("field %q must be at least %d characters", r.Field, r.MinLen))
		}

		if r.MaxLen > 0 && len(val) > r.MaxLen {
			errs = append(errs, fmt.Sprintf("field %q must be at most %d characters", r.Field, r.MaxLen))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errs: errs}
	}
	return nil
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
