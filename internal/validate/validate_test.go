package validate_test

import (
	"strings"
	"testing"

	"github.com/your-org/grpcurl-batch/internal/validate"
)

func TestValidate_NoRules_AlwaysPasses(t *testing.T) {
	v := validate.New(nil)
	if err := v.Validate(map[string]string{"foo": "bar"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_RequiredField_Present(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "address", Required: true}})
	if err := v.Validate(map[string]string{"address": "localhost:50051"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_RequiredField_Missing(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "address", Required: true}})
	err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
	if !validate.IsValidationError(err) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if !strings.Contains(err.Error(), "address") {
		t.Errorf("error should mention field name, got: %v", err)
	}
}

func TestValidate_RequiredField_Blank(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "method", Required: true}})
	err := v.Validate(map[string]string{"method": "   "})
	if err == nil {
		t.Fatal("expected error for blank required field")
	}
}

func TestValidate_MinLen_Violated(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "token", MinLen: 8}})
	err := v.Validate(map[string]string{"token": "abc"})
	if err == nil {
		t.Fatal("expected min-length error")
	}
	if !strings.Contains(err.Error(), "at least 8") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MaxLen_Violated(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "label", MaxLen: 5}})
	err := v.Validate(map[string]string{"label": "toolongvalue"})
	if err == nil {
		t.Fatal("expected max-length error")
	}
	if !strings.Contains(err.Error(), "at most 5") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MultipleErrors_Collected(t *testing.T) {
	v := validate.New([]validate.Rule{
		{Field: "address", Required: true},
		{Field: "method", Required: true},
	})
	err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected errors")
	}
	ve, ok := err.(*validate.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(ve.Errs), ve.Errs)
	}
}

func TestValidate_OptionalField_AbsentSkipsLenCheck(t *testing.T) {
	v := validate.New([]validate.Rule{{Field: "tag", MinLen: 3}})
	if err := v.Validate(map[string]string{}); err != nil {
		t.Fatalf("absent optional field should not trigger length check, got %v", err)
	}
}
