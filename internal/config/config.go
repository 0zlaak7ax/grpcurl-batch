package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Request defines a single gRPC request to be executed.
type Request struct {
	Name    string            `yaml:"name"`
	Address string            `yaml:"address"`
	Service string            `yaml:"service"`
	Method  string            `yaml:"method"`
	Data    string            `yaml:"data"`
	Headers map[string]string `yaml:"headers"`
}

// RetryConfig holds retry policy settings.
type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts"`
	Delay       time.Duration `yaml:"delay"`
}

// OutputConfig controls how results are formatted.
type OutputConfig struct {
	Format string `yaml:"format"` // json, text
	Verbose bool   `yaml:"verbose"`
}

// BatchConfig is the top-level configuration loaded from a YAML file.
type BatchConfig struct {
	Requests []Request    `yaml:"requests"`
	Retry    RetryConfig  `yaml:"retry"`
	Output   OutputConfig `yaml:"output"`
}

// Load reads and parses a YAML batch config file.
func Load(path string) (*BatchConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg BatchConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks required fields in the config.
func (c *BatchConfig) validate() error {
	for i, r := range c.Requests {
		if r.Address == "" {
			return fmt.Errorf("request[%d] %q: address is required", i, r.Name)
		}
		if r.Service == "" {
			return fmt.Errorf("request[%d] %q: service is required", i, r.Name)
		}
		if r.Method == "" {
			return fmt.Errorf("request[%d] %q: method is required", i, r.Name)
		}
	}
	if c.Retry.MaxAttempts == 0 {
		c.Retry.MaxAttempts = 1
	}
	if c.Output.Format == "" {
		c.Output.Format = "json"
	}
	return nil
}
