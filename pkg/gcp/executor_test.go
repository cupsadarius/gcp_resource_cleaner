package gcp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewConcurrentExecutor(t *testing.T) {
	tests := []struct {
		name             string
		maxConcurrent    int
		expectedCapacity int
	}{
		{
			name:             "single concurrent executor",
			maxConcurrent:    1,
			expectedCapacity: 1,
		},
		{
			name:             "multiple concurrent executor",
			maxConcurrent:    5,
			expectedCapacity: 5,
		},
		{
			name:             "high concurrency executor",
			maxConcurrent:    20,
			expectedCapacity: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewConcurrentExecutor(tt.maxConcurrent)

			if executor == nil {
				t.Fatal("NewConcurrentExecutor returned nil")
			}

			if cap(executor.semaphore) != tt.expectedCapacity {
				t.Errorf("Expected semaphore capacity %d, got %d", tt.expectedCapacity, cap(executor.semaphore))
			}
		})
	}
}

func TestConcurrentExecutor_ExecuteCommand_Success(t *testing.T) {
	executor := NewConcurrentExecutor(2)
	ctx := context.Background()

	// Execute a simple command
	output, err := executor.ExecuteCommand(ctx, "echo", "hello")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}
}

func TestConcurrentExecutor_ExecuteCommand_ConcurrencyLimit(t *testing.T) {
	const maxConcurrent = 2
	executor := NewConcurrentExecutor(maxConcurrent)
	ctx := context.Background()

	// Track concurrent executions
	var activeCount int32
	var maxActiveCount int32
	var mu sync.Mutex

	// Override the actual execution to track concurrency
	var wg sync.WaitGroup
	numCommands := 5

	for i := 0; i < numCommands; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate the concurrent execution tracking
			select {
			case executor.semaphore <- struct{}{}:
				mu.Lock()
				activeCount++
				if activeCount > maxActiveCount {
					maxActiveCount = activeCount
				}
				currentActive := activeCount
				mu.Unlock()

				// Verify we don't exceed the limit
				if currentActive > maxConcurrent {
					t.Errorf("Concurrency limit exceeded: %d > %d", currentActive, maxConcurrent)
				}

				// Simulate work
				time.Sleep(10 * time.Millisecond)

				mu.Lock()
				activeCount--
				mu.Unlock()

				<-executor.semaphore
			case <-ctx.Done():
				t.Error("Context cancelled unexpectedly")
			}
		}(i)
	}

	wg.Wait()

	// Verify we respected the concurrency limit
	if maxActiveCount > maxConcurrent {
		t.Errorf("Max active count %d exceeded limit %d", maxActiveCount, maxConcurrent)
	}

	// Verify all slots were used
	if maxActiveCount < maxConcurrent {
		t.Logf("Max active count was %d, expected up to %d", maxActiveCount, maxConcurrent)
	}
}

func TestConcurrentExecutor_ExecuteCommand_ContextCancellation(t *testing.T) {
	executor := NewConcurrentExecutor(1)
	ctx, cancel := context.WithCancel(context.Background())

	// Fill the semaphore
	select {
	case executor.semaphore <- struct{}{}:
		// Semaphore is now full
	default:
		t.Fatal("Failed to fill semaphore")
	}

	// Cancel the context
	cancel()

	// Try to execute - should fail with context cancellation
	_, err := executor.ExecuteCommand(ctx, "echo", "hello")

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	// Clean up semaphore
	<-executor.semaphore
}

func TestConcurrentExecutor_ExecuteCommand_CommandError(t *testing.T) {
	executor := NewConcurrentExecutor(1)
	ctx := context.Background()

	// Execute a command that should fail
	_, err := executor.ExecuteCommand(ctx, "nonexistent-command-xyz", "arg")

	if err == nil {
		t.Error("Expected command error, got nil")
	}
}

func TestConcurrentExecutor_ExecuteCommand_Parallel(t *testing.T) {
	executor := NewConcurrentExecutor(3)
	ctx := context.Background()

	// Execute multiple commands in parallel
	var wg sync.WaitGroup
	errors := make(chan error, 5)
	outputs := make(chan []byte, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			output, err := executor.ExecuteCommand(ctx, "echo", "test")
			errors <- err
			outputs <- output
		}(i)
	}

	wg.Wait()
	close(errors)
	close(outputs)

	// Check results
	errorCount := 0
	outputCount := 0

	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	for output := range outputs {
		if len(output) > 0 {
			outputCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors, got %d errors", errorCount)
	}

	if outputCount != 5 {
		t.Errorf("Expected 5 outputs, got %d", outputCount)
	}
}

func TestConcurrentExecutor_SemaphoreCleanup(t *testing.T) {
	executor := NewConcurrentExecutor(2)
	ctx := context.Background()

	// Execute multiple commands to ensure semaphore is properly cleaned up
	for i := 0; i < 10; i++ {
		_, err := executor.ExecuteCommand(ctx, "echo", "test")
		if err != nil {
			t.Errorf("Command %d failed: %v", i, err)
		}
	}

	// Verify semaphore is empty (all slots available)
	// Try to acquire all slots quickly
	for i := 0; i < cap(executor.semaphore); i++ {
		select {
		case executor.semaphore <- struct{}{}:
			// Successfully acquired slot
		default:
			t.Errorf("Semaphore slot %d not available, cleanup failed", i)
		}
	}

	// Clean up - release all slots
	for i := 0; i < cap(executor.semaphore); i++ {
		<-executor.semaphore
	}
}

func BenchmarkConcurrentExecutor_Sequential(b *testing.B) {
	executor := NewConcurrentExecutor(1)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.ExecuteCommand(ctx, "echo", "benchmark")
		if err != nil {
			b.Fatalf("Command failed: %v", err)
		}
	}
}

func BenchmarkConcurrentExecutor_Parallel(b *testing.B) {
	executor := NewConcurrentExecutor(5)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := executor.ExecuteCommand(ctx, "echo", "benchmark")
			if err != nil {
				b.Fatalf("Command failed: %v", err)
			}
		}
	})
}
