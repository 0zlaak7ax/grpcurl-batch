package executor_test

import (
	"testing"

	"github.com/user/grpcurl-batch/internal/config"
	"github.com/user/grpcurl-batch/internal/executor"
)

func TestNew_DefaultBinary(t *testing.T) {
	cfg := &config.Config{Address: "localhost:50051"}
	exec := executor.New(cfg)
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
}

func TestNew_CustomBinary(t *testing.T) {
	cfg := &config.Config{
		Address:       "localhost:50051",
		GrpcurlBinary: "/usr/local/bin/grpcurl",
	}
	exec := executor.New(cfg)
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
}

func TestExecute_BinaryNotFound(t *testing.T) {
	cfg := &config.Config{
		Address:       "localhost:50051",
		GrpcurlBinary: "/nonexistent/grpcurl-binary",
	}
	exec := executor.New(cfg)
	_, err := exec.Execute(t.Context(), config.Request{
		Name:   "test",
		Method: "pkg.Svc/Method",
	})
	if err == nil {
		t.Error("expected error when binary not found")
	}
}

func TestExecute_InsecureFlag(t *testing.T) {
	// This test validates that insecure configs don't panic during construction.
	cfg := &config.Config{
		Address:  "localhost:50051",
		Insecure: true,
	}
	exec := executor.New(cfg)
	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
}
