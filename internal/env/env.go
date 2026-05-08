// Package env provides helpers for resolving configuration values from
// environment variables with optional defaults and type coercion.
package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetString returns the value of the named environment variable, or
// defaultVal if the variable is unset or empty.
func GetString(name, defaultVal string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultVal
}

// GetInt returns the integer value of the named environment variable.
// If the variable is unset, empty, or cannot be parsed, defaultVal is
// returned.
func GetInt(name string, defaultVal int) int {
	v := os.Getenv(name)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return defaultVal
	}
	return n
}

// GetBool returns the boolean value of the named environment variable.
// Accepted truthy values: "1", "true", "yes", "on" (case-insensitive).
// Any other non-empty value is treated as false. If unset, defaultVal
// is returned.
func GetBool(name string, defaultVal bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	if v == "" {
		return defaultVal
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// GetDuration returns the time.Duration value of the named environment
// variable parsed via time.ParseDuration. If unset, empty, or invalid,
// defaultVal is returned.
func GetDuration(name string, defaultVal time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(name))
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}

// Require returns the value of the named environment variable or returns
// an error if the variable is unset or empty.
func Require(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("env: required variable %q is not set", name)
	}
	return v, nil
}
