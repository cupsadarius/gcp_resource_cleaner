package gcp

import (
	"context"
	"os/exec"
)

// CommandExecutor defines the interface for executing external commands
type CommandExecutor interface {
	ExecuteCommand(ctx context.Context, name string, args ...string) ([]byte, error)
}

// GCloudExecutor is the real implementation that executes gcloud commands
type GCloudExecutor struct{}

// ExecuteCommand executes the actual gcloud command
func (g *GCloudExecutor) ExecuteCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.CombinedOutput()
}

// MockExecutor is a test implementation that returns predefined responses
type MockExecutor struct {
	MockOutput []byte
	MockError  error
	CallLog    []CommandCall // For verifying what commands were called
}

// CommandCall represents a command that was executed
type CommandCall struct {
	Name string
	Args []string
}

// ExecuteCommand returns the mock response without executing anything
func (m *MockExecutor) ExecuteCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	// Log the call for verification in tests
	m.CallLog = append(m.CallLog, CommandCall{
		Name: name,
		Args: args,
	})

	return m.MockOutput, m.MockError
}

// Reset clears the call log (useful between tests)
func (m *MockExecutor) Reset() {
	m.CallLog = nil
	m.MockOutput = nil
	m.MockError = nil
}

// GetLastCall returns the most recent command call
func (m *MockExecutor) GetLastCall() *CommandCall {
	if len(m.CallLog) == 0 {
		return nil
	}
	return &m.CallLog[len(m.CallLog)-1]
}

// GetCallCount returns the number of commands executed
func (m *MockExecutor) GetCallCount() int {
	return len(m.CallLog)
}
