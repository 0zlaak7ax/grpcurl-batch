package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/grpcurl-batch/internal/config"
)

// GrpcurlExecutor runs requests using the grpcurl CLI binary.
type GrpcurlExecutor struct {
	address  string
	insecure bool
	binary   string
}

// New creates a GrpcurlExecutor targeting the given address.
func New(cfg *config.Config) *GrpcurlExecutor {
	binary := "grpcurl"
	if cfg.GrpcurlBinary != "" {
		binary = cfg.GrpcurlBinary
	}
	return &GrpcurlExecutor{
		address:  cfg.Address,
		insecure: cfg.Insecure,
		binary:   binary,
	}
}

// Execute builds and runs a grpcurl command for the given request.
func (e *GrpcurlExecutor) Execute(ctx context.Context, req config.Request) (string, error) {
	args := e.buildArgs(req)
	cmd := exec.CommandContext(ctx, e.binary, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("grpcurl error: %w — %s", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (e *GrpcurlExecutor) buildArgs(req config.Request) []string {
	args := []string{}
	if e.insecure {
		args = append(args, "-insecure")
	}
	for k, v := range req.Metadata {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if req.Data != "" {
		args = append(args, "-d", req.Data)
	}
	args = append(args, e.address, req.Method)
	return args
}
