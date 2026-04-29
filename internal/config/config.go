package config

import (
	"errors"
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

	if cfg.Address == "" {
		return nil, errors.New("config: address is required")
	}

	if cfg.Retry.MaxAttempts == 0 {
		cfg.Retry.MaxAttempts = 3
	}
	if cfg.Retry.Delay == 0 {
		cfg.Retry.Delay = 500 * time.Millisecond
	}

	return &cfg, nil
}
