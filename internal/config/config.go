// Package config loads and validates the YAML batch configuration file.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Request describes a single gRPC call to execute.
type Request struct {
	Name    string            `yaml:"name"`
	Method  string            `yaml:"method"`
	Data    string            `yaml:"data"`
	Headers map[string]string `yaml:"headers"`
}

// Config is the top-level configuration loaded from a YAML file.
type Config struct {
	Address       string        `yaml:"address"`
	Insecure      bool          `yaml:"insecure"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
	MaxConcurrent int           `yaml:"max_concurrent"`
	RateInterval  time.Duration `yaml:"rate_interval"`
	OutputFormat  string        `yaml:"output_format"`
	JUnitOutput   string        `yaml:"junit_output"`
	Requests      []Request     `yaml:"requests"`
}

// Load reads and validates a Config from the YAML file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Address == "" {
		return nil, errors.New("config: address is required")
	}

	// Apply defaults.
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay == 0 {
		cfg.RetryDelay = 500 * time.Millisecond
	}
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 1
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "text"
	}

	return &cfg, nil
}
