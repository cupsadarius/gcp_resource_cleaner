package errors

import (
	"errors"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNotInitialized",
			err:      ErrNotInitialized,
			expected: "not initialized",
		},
		{
			name:     "ErrTargetNotPointer",
			err:      ErrTargetNotPointer,
			expected: "target is not pointer",
		},
		{
			name:     "ErrFileDoesNotExist",
			err:      ErrFileDoesNotExist,
			expected: "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message to be '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that our errors are proper error types
	var err error

	err = ErrNotInitialized
	if err == nil {
		t.Error("ErrNotInitialized should not be nil")
	}

	err = ErrTargetNotPointer
	if err == nil {
		t.Error("ErrTargetNotPointer should not be nil")
	}

	err = ErrFileDoesNotExist
	if err == nil {
		t.Error("ErrFileDoesNotExist should not be nil")
	}
}

func TestErrorComparison(t *testing.T) {
	// Test that errors can be compared using errors.Is
	var err = ErrNotInitialized

	if !errors.Is(err, ErrNotInitialized) {
		t.Error("ErrNotInitialized should match itself using errors.Is")
	}

	if errors.Is(err, ErrTargetNotPointer) {
		t.Error("ErrNotInitialized should not match ErrTargetNotPointer")
	}

	if errors.Is(err, ErrFileDoesNotExist) {
		t.Error("ErrNotInitialized should not match ErrFileDoesNotExist")
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test that our errors can be wrapped
	wrappedErr := errors.Join(ErrNotInitialized, errors.New("additional context"))

	if !errors.Is(wrappedErr, ErrNotInitialized) {
		t.Error("Wrapped error should still match original error using errors.Is")
	}

	if wrappedErr.Error() == ErrNotInitialized.Error() {
		t.Error("Wrapped error should have different message than original")
	}
}
