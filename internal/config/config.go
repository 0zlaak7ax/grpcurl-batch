package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// RetryConfig controls retry behaviour for failed requests.
type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts"`
	Delay       time.Duration `yaml:"delay"`
}

// Request represents a single gRPC call definition.
type Request struct {
	Name     string            `yaml:"name"`
	Method   string            `yaml:"method"`
	Data     string            `yaml:"data"`
	Metadata map[string]string `yaml:"metadata"`
}

// Config is the top-level configuration loaded from a YAML file.
type Config struct {
	Address       string      `yaml:"address"`
	Insecure      bool        `yaml:"insecure"`
	GrpcurlBinary string      `yaml:"grpcurl_binary"`
	Retry         RetryConfig `yaml:"retry"`
	Requests      []Request   `yaml:"requests"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if cfg.Retry.MaxAttempts == 0 {
		cfg.Retry.MaxAttempts = 3
	}
	if cfg.Retry.Delay == 0 {
		cfg.Retry.Delay = 500 * time.Millisecond
	}

	return &cfg, nil
}

// validate checks that all required fields are present and that request
// definitions are well-formed.
func (c *Config) validate() error {
	if c.Address == "" {
		return errors.New("config: address is required")
	}
	for i, r := range c.Requests {
		if r.Method == "" {
			return fmt.Errorf("config: request[%d] (%q): method is required", i, r.Name)
		}
	}
	return nil
}
