package gcp

import (
	"context"
	"errors"
	"testing"

	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func init() {
	// Initialize logger for tests
	logger.Init(logger.Config{
		Level:  "error", // Reduce log noise in tests
		Source: "test",
		Format: "json",
	})
}

func TestCheckHealth_Success(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Google Cloud SDK 400.0.0\nbq 2.0.75\ncore 2022.08.19\ngsutil 5.12\n"),
		MockError:  nil,
	}

	ctx := context.Background()

	// CheckHealth doesn't return an error, it just logs
	// This test ensures it doesn't panic and calls the right command
	CheckHealth(ctx, mockExec)

	// Verify the correct command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}

	expectedArgs := []string{"version"}
	if len(lastCall.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
	}

	for i, arg := range lastCall.Args {
		if arg != expectedArgs[i] {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expectedArgs[i], arg)
		}
	}
}

func TestCheckHealth_CommandError(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  errors.New("gcloud not found"),
	}

	ctx := context.Background()

	// CheckHealth doesn't return an error, it just logs
	// This test ensures it doesn't panic even when gcloud fails
	CheckHealth(ctx, mockExec)

	// Verify the command was attempted
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call attempt, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}
}

func TestCheckHealth_EmptyOutput(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  nil,
	}

	ctx := context.Background()

	// CheckHealth doesn't return an error, it just logs
	// This test ensures it handles empty output gracefully
	CheckHealth(ctx, mockExec)

	// Verify the command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}
}

func TestCheckHealth_MultilineOutput(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(`Google Cloud SDK 400.0.0
bq 2.0.75
core 2022.08.19
gsutil 5.12
Updates are available for some Google Cloud CLI components.  To install them,
please run:
  $ gcloud components update
`),
		MockError: nil,
	}

	ctx := context.Background()

	// This test verifies that multiline output is handled correctly
	CheckHealth(ctx, mockExec)

	// Verify the command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}

	if len(lastCall.Args) != 1 || lastCall.Args[0] != "version" {
		t.Errorf("Expected args to be ['version'], got %v", lastCall.Args)
	}
}

func TestCheckHealth_WithContext(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Google Cloud SDK 400.0.0\n"),
		MockError:  nil,
	}

	// Test with a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	CheckHealth(ctx, mockExec)

	// Verify the command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}
}

func TestCheckHealth_ContextCancellation(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Google Cloud SDK 400.0.0\n"),
		MockError:  nil,
	}

	// Test with a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	CheckHealth(ctx, mockExec)

	// Even with a cancelled context, the function should not panic
	// The command might still be called since cancellation is handled in the executor
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}
}

func TestCheckHealth_GcloudNotInstalled(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  errors.New("exec: \"gcloud\": executable file not found in $PATH"),
	}

	ctx := context.Background()

	// This simulates the case where gcloud is not installed
	CheckHealth(ctx, mockExec)

	// Verify the command was attempted
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call attempt, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}
}

func TestCheckHealth_OldGcloudVersion(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Google Cloud SDK 350.0.0\nbq 1.0.75\ncore 2021.08.19\ngsutil 4.12\n"),
		MockError:  nil,
	}

	ctx := context.Background()

	CheckHealth(ctx, mockExec)

	// Verify the command was called successfully
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}
}

func TestCheckHealth_CallSequence(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Google Cloud SDK 400.0.0\n"),
		MockError:  nil,
	}

	ctx := context.Background()

	// Call CheckHealth multiple times
	CheckHealth(ctx, mockExec)
	CheckHealth(ctx, mockExec)
	CheckHealth(ctx, mockExec)

	// Verify all calls were made
	if mockExec.GetCallCount() != 3 {
		t.Errorf("Expected 3 command calls, got %d", mockExec.GetCallCount())
	}

	// Verify all calls were the same
	for i, call := range mockExec.CallLog {
		if call.Name != "gcloud" {
			t.Errorf("Call %d: expected command to be 'gcloud', got %s", i, call.Name)
		}
		if len(call.Args) != 1 || call.Args[0] != "version" {
			t.Errorf("Call %d: expected args to be ['version'], got %v", i, call.Args)
		}
	}
}
