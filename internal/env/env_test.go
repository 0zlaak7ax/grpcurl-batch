package env_test

import (
	"testing"
	"time"

	"github.com/your-org/grpcurl-batch/internal/env"
)

func TestGetString_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_STR", "hello")
	if got := env.GetString("TEST_STR", "default"); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestGetString_ReturnsDefault_WhenUnset(t *testing.T) {
	t.Setenv("TEST_STR_UNSET", "")
	if got := env.GetString("TEST_STR_UNSET", "fallback"); got != "fallback" {
		t.Fatalf("expected %q, got %q", "fallback", got)
	}
}

func TestGetInt_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	if got := env.GetInt("TEST_INT", 0); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestGetInt_ReturnsDefault_OnInvalid(t *testing.T) {
	t.Setenv("TEST_INT_BAD", "not-a-number")
	if got := env.GetInt("TEST_INT_BAD", 7); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestGetInt_ReturnsDefault_WhenUnset(t *testing.T) {
	if got := env.GetInt("TEST_INT_MISSING_XYZ", 99); got != 99 {
		t.Fatalf("expected 99, got %d", got)
	}
}

func TestGetBool_TruthyValues(t *testing.T) {
	for _, v := range []string{"1", "true", "TRUE", "yes", "on"} {
		t.Setenv("TEST_BOOL", v)
		if got := env.GetBool("TEST_BOOL", false); !got {
			t.Fatalf("expected true for value %q", v)
		}
	}
}

func TestGetBool_FalsyValue(t *testing.T) {
	t.Setenv("TEST_BOOL_F", "false")
	if got := env.GetBool("TEST_BOOL_F", true); got {
		t.Fatal("expected false")
	}
}

func TestGetBool_ReturnsDefault_WhenUnset(t *testing.T) {
	if got := env.GetBool("TEST_BOOL_MISSING_XYZ", true); !got {
		t.Fatal("expected default true")
	}
}

func TestGetDuration_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_DUR", "5s")
	if got := env.GetDuration("TEST_DUR", time.Second); got != 5*time.Second {
		t.Fatalf("expected 5s, got %v", got)
	}
}

func TestGetDuration_ReturnsDefault_OnInvalid(t *testing.T) {
	t.Setenv("TEST_DUR_BAD", "nope")
	if got := env.GetDuration("TEST_DUR_BAD", 3*time.Second); got != 3*time.Second {
		t.Fatalf("expected 3s, got %v", got)
	}
}

func TestRequire_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_REQ", "present")
	v, err := env.Require("TEST_REQ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "present" {
		t.Fatalf("expected %q, got %q", "present", v)
	}
}

func TestRequire_ReturnsError_WhenUnset(t *testing.T) {
	t.Setenv("TEST_REQ_MISSING", "")
	_, err := env.Require("TEST_REQ_MISSING")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
