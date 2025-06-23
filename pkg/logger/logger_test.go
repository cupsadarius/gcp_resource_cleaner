package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestNew(t *testing.T) {
	logger := New("test-section", "test-action")

	if logger.section != "test-section" {
		t.Errorf("Expected section to be 'test-section', got %s", logger.section)
	}

	if logger.action != "test-action" {
		t.Errorf("Expected action to be 'test-action', got %s", logger.action)
	}
}

func TestLogger_Zerolog(t *testing.T) {
	logger := New("test-section", "test-action")

	// Create a test event
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf)
	event := testLogger.Info()

	// Apply the logger's zerolog method
	enrichedEvent := logger.Zerolog(event)

	// Execute the event to capture output
	enrichedEvent.Msg("test message")

	output := buf.String()

	// Check that section and action were added
	if !strings.Contains(output, `"section":"test-section"`) {
		t.Error("Expected section to be added to log event")
	}

	if !strings.Contains(output, `"action":"test-action"`) {
		t.Error("Expected action to be added to log event")
	}

	if !strings.Contains(output, `"message":"test message"`) {
		t.Error("Expected message to be preserved in log event")
	}
}

func TestStringLevelToZerologLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected zerolog.Level
	}{
		{"trace", zerolog.TraceLevel},
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"panic", zerolog.PanicLevel},
		{"TRACE", zerolog.TraceLevel},  // Test case insensitivity
		{"Debug", zerolog.DebugLevel},  // Test mixed case
		{"invalid", zerolog.InfoLevel}, // Test invalid input defaults to info
		{"", zerolog.InfoLevel},        // Test empty string defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stringLevelToZerologLevel(tt.input)
			if result != tt.expected {
				t.Errorf("Expected level %v for input '%s', got %v", tt.expected, tt.input, result)
			}
		})
	}
}

func TestConfig_Init(t *testing.T) {
	// Test pretty format initialization
	cfg := Config{
		Source: "test-app",
		Level:  "debug",
		Format: "pretty",
	}

	// Capture the original logger to restore later
	originalLogger := log.Logger

	Init(cfg)

	// Verify that global log level was set
	if zerolog.GlobalLevel() != zerolog.DebugLevel {
		t.Errorf("Expected global level to be debug, got %v", zerolog.GlobalLevel())
	}

	// Restore original logger
	log.Logger = originalLogger
}

func TestConfig_InitJSON(t *testing.T) {
	// Test JSON format (default) initialization
	cfg := Config{
		Source: "test-app",
		Level:  "warn",
		Format: "json",
	}

	// Capture the original logger to restore later
	originalLogger := log.Logger

	Init(cfg)

	// Verify that global log level was set
	if zerolog.GlobalLevel() != zerolog.WarnLevel {
		t.Errorf("Expected global level to be warn, got %v", zerolog.GlobalLevel())
	}

	// Restore original logger
	log.Logger = originalLogger
}

func TestLogger_Methods(t *testing.T) {
	// Redirect log output to capture it
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)
	
	// Set global log level to trace to capture all messages
	originalLevel := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	defer zerolog.SetGlobalLevel(originalLevel)

	logger := New("test-section", "test-action")

	// Test Info
	logger.Info("info message")
	fmt.Println(buf.String())
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message not logged")
	}
	buf.Reset()

	// Test Debug
	logger.Debug("debug message")
	fmt.Println(buf.String())
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message not logged")
	}
	buf.Reset()

	// Test DebugWithExtra
	extra := map[string]any{
		"key1": "value1",
		"key2": 42,
	}
	logger.DebugWithExtra("debug with extra", extra)
	fmt.Println(buf.String())
	output := buf.String()
	if !strings.Contains(output, "debug with extra") {
		t.Error("DebugWithExtra message not logged")
	}
	if !strings.Contains(output, "value1") {
		t.Error("Extra data not logged")
	}
	buf.Reset()

	// Test Warn
	logger.Warn("warn message")
	fmt.Println(buf.String())
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn message not logged")
	}
	buf.Reset()
}
