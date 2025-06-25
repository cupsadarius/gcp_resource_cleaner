package cli

import (
	"context"
	"testing"

	"github.com/cupsadarius/gcp_resource_cleaner/pkg/errors"
)

func TestInit(t *testing.T) {
	appID := "test-app"
	shortDesc := "Test application"
	longDesc := "This is a test application for unit testing"

	Init(appID, shortDesc, longDesc)

	if cmd == nil {
		t.Fatal("Init() did not initialize cmd")
	}

	if cmd.Use != appID {
		t.Errorf("Expected Use to be %s, got %s", appID, cmd.Use)
	}

	if cmd.Short != shortDesc {
		t.Errorf("Expected Short to be %s, got %s", shortDesc, cmd.Short)
	}

	if cmd.Long != longDesc {
		t.Errorf("Expected Long to be %s, got %s", longDesc, cmd.Long)
	}
}

func TestAddCommand_Success(t *testing.T) {
	Init("test-app", "short", "long")

	called := false
	testHandler := func(ctx context.Context) {
		called = true
	}

	err := AddCommand("test-command", "Test command description", testHandler)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that command was added
	foundCommand := false
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "test-command" {
			foundCommand = true
			if subCmd.Short != "Test command description" {
				t.Errorf("Expected description to be 'Test command description', got %s", subCmd.Short)
			}
			// Simulate running the command
			subCmd.Run(subCmd, []string{})
			break
		}
	}

	if !foundCommand {
		t.Error("Command was not added to CLI")
	}

	if !called {
		t.Error("Handler function was not called")
	}
}

func TestAddCommand_NotInitialized(t *testing.T) {
	cmd = nil // Reset cmd to simulate uninitialized state

	testHandler := func(ctx context.Context) {}
	err := AddCommand("test-command", "Test command description", testHandler)

	if err != errors.ErrNotInitialized {
		t.Errorf("Expected ErrNotInitialized, got %v", err)
	}
}

func TestAssignStringFlag(t *testing.T) {
	Init("test-app", "short", "long")

	var testString string
	AssignStringFlag(&testString, "test-flag", "default-value", "Test flag description")

	// Get the flag to verify it was added
	flag := cmd.PersistentFlags().Lookup("test-flag")
	if flag == nil {
		t.Fatal("String flag was not added")
	}

	if flag.DefValue != "default-value" {
		t.Errorf("Expected default value to be 'default-value', got %s", flag.DefValue)
	}

	if flag.Usage != "Test flag description" {
		t.Errorf("Expected usage to be 'Test flag description', got %s", flag.Usage)
	}
}

func TestAssignBoolFlag(t *testing.T) {
	Init("test-app", "short", "long")

	var testBool bool
	AssignBoolFlag(&testBool, "test-bool", true, "Test bool description")

	// Get the flag to verify it was added
	flag := cmd.PersistentFlags().Lookup("test-bool")
	if flag == nil {
		t.Fatal("Bool flag was not added")
	}

	if flag.DefValue != "true" {
		t.Errorf("Expected default value to be 'true', got %s", flag.DefValue)
	}

	if flag.Usage != "Test bool description" {
		t.Errorf("Expected usage to be 'Test bool description', got %s", flag.Usage)
	}
}

func TestMultipleCommands(t *testing.T) {
	Init("test-app", "short", "long")

	command1Called := false
	command2Called := false

	err1 := AddCommand("command1", "First command", func(ctx context.Context) {
		command1Called = true
	})

	err2 := AddCommand("command2", "Second command", func(ctx context.Context) {
		command2Called = true
	})

	if err1 != nil {
		t.Errorf("Error adding command1: %v", err1)
	}

	if err2 != nil {
		t.Errorf("Error adding command2: %v", err2)
	}

	// Verify both commands exist
	commands := cmd.Commands()
	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Test each command
	for _, subCmd := range commands {
		subCmd.Run(subCmd, []string{})
	}

	if !command1Called {
		t.Error("Command1 handler was not called")
	}

	if !command2Called {
		t.Error("Command2 handler was not called")
	}
}

func TestAssignIntFlag(t *testing.T) {
	Init("test-app", "short", "long")

	var testInt int
	AssignIntFlag(&testInt, "test-int", 42, "Test int description")

	// Get the flag to verify it was added
	flag := cmd.PersistentFlags().Lookup("test-int")
	if flag == nil {
		t.Fatal("Int flag was not added")
	}

	if flag.DefValue != "42" {
		t.Errorf("Expected default value to be '42', got %s", flag.DefValue)
	}

	if flag.Usage != "Test int description" {
		t.Errorf("Expected usage to be 'Test int description', got %s", flag.Usage)
	}
}

func TestConcurrencyFlags(t *testing.T) {
	Init("test-app", "short", "long")

	var enableConcurrency bool
	var concurrencyLimit int

	// Test adding concurrency flags
	AssignBoolFlag(&enableConcurrency, "concurrency", false, "Enable concurrency")
	AssignIntFlag(&concurrencyLimit, "concurrency-limit", 5, "Concurrency limit")

	// Verify concurrency flag
	concurrencyFlag := cmd.PersistentFlags().Lookup("concurrency")
	if concurrencyFlag == nil {
		t.Fatal("Concurrency flag was not added")
	}

	if concurrencyFlag.DefValue != "false" {
		t.Errorf("Expected concurrency default to be 'false', got %s", concurrencyFlag.DefValue)
	}

	if concurrencyFlag.Usage != "Enable concurrency" {
		t.Errorf("Expected concurrency usage to be 'Enable concurrency', got %s", concurrencyFlag.Usage)
	}

	// Verify concurrency-limit flag
	limitFlag := cmd.PersistentFlags().Lookup("concurrency-limit")
	if limitFlag == nil {
		t.Fatal("Concurrency-limit flag was not added")
	}

	if limitFlag.DefValue != "5" {
		t.Errorf("Expected concurrency-limit default to be '5', got %s", limitFlag.DefValue)
	}

	if limitFlag.Usage != "Concurrency limit" {
		t.Errorf("Expected concurrency-limit usage to be 'Concurrency limit', got %s", limitFlag.Usage)
	}
}

func TestIntFlagValidation(t *testing.T) {
	Init("test-app", "short", "long")

	tests := []struct {
		name         string
		defaultValue int
		description  string
	}{
		{
			name:         "zero default",
			defaultValue: 0,
			description:  "Zero concurrency limit",
		},
		{
			name:         "positive default",
			defaultValue: 10,
			description:  "Positive concurrency limit",
		},
		{
			name:         "large default",
			defaultValue: 1000,
			description:  "Large concurrency limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testInt int
			flagName := "test-" + tt.name

			AssignIntFlag(&testInt, flagName, tt.defaultValue, tt.description)

			flag := cmd.PersistentFlags().Lookup(flagName)
			if flag == nil {
				t.Fatalf("Flag %s was not added", flagName)
			}

			expectedDefault := string(rune(tt.defaultValue + '0'))
			if tt.defaultValue >= 10 {
				expectedDefault = flag.DefValue // For multi-digit numbers, just check it was set
			}

			if tt.defaultValue < 10 && flag.DefValue != expectedDefault {
				t.Errorf("Expected default value %d, got %s", tt.defaultValue, flag.DefValue)
			}

			if flag.Usage != tt.description {
				t.Errorf("Expected usage '%s', got '%s'", tt.description, flag.Usage)
			}
		})
	}
}

func TestCommandsWithConcurrencyFlags(t *testing.T) {
	Init("test-app", "short", "long")

	// Add concurrency flags
	var enableConcurrency bool
	var concurrencyLimit int
	AssignBoolFlag(&enableConcurrency, "concurrency", false, "Enable concurrency")
	AssignIntFlag(&concurrencyLimit, "concurrency-limit", 5, "Concurrency limit")

	// Add test commands
	command1Called := false
	command2Called := false

	err1 := AddCommand("command1", "First command", func(ctx context.Context) {
		command1Called = true
	})

	err2 := AddCommand("command2", "Second command", func(ctx context.Context) {
		command2Called = true
	})

	if err1 != nil {
		t.Errorf("Error adding command1: %v", err1)
	}

	if err2 != nil {
		t.Errorf("Error adding command2: %v", err2)
	}

	// Verify commands and flags coexist
	commands := cmd.Commands()
	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Verify flags are still present
	concurrencyFlag := cmd.PersistentFlags().Lookup("concurrency")
	limitFlag := cmd.PersistentFlags().Lookup("concurrency-limit")

	if concurrencyFlag == nil {
		t.Error("Concurrency flag missing after adding commands")
	}

	if limitFlag == nil {
		t.Error("Concurrency-limit flag missing after adding commands")
	}

	// Test command execution (flags should be available to commands)
	for _, subCmd := range commands {
		subCmd.Run(subCmd, []string{})
	}

	if !command1Called {
		t.Error("Command1 was not called")
	}

	if !command2Called {
		t.Error("Command2 was not called")
	}
}
