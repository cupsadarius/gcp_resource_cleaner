package version

import (
	"testing"
)

func TestVersionVariables(t *testing.T) {
	// Test that version variables exist and can be set
	originalAppVersion := AppVersion
	originalGitCommit := GitCommit

	// Test setting values
	AppVersion = "1.0.0"
	GitCommit = "abc123def456"

	if AppVersion != "1.0.0" {
		t.Errorf("Expected AppVersion to be '1.0.0', got '%s'", AppVersion)
	}

	if GitCommit != "abc123def456" {
		t.Errorf("Expected GitCommit to be 'abc123def456', got '%s'", GitCommit)
	}

	// Restore original values
	AppVersion = originalAppVersion
	GitCommit = originalGitCommit
}

func TestDefaultValues(t *testing.T) {
	// Test that variables have empty defaults when not set during build
	// Note: In real builds, these would be set via -ldflags

	// We can't really test the "unset" state in unit tests since the variables
	// are package-level globals, but we can test that they're strings

	var appVersionType string = AppVersion
	var gitCommitType string = GitCommit

	// This test mainly ensures the variables exist and are the right type
	_ = appVersionType
	_ = gitCommitType

	// In actual usage, these would be set via build flags like:
	// go build -ldflags "-X github.com/cupsadarius/gcp_resource_cleaner/pkg/version.AppVersion=1.0.0"
}

func TestVersionStringConcatenation(t *testing.T) {
	// Test that versions can be used in string operations
	AppVersion = "1.0.0"
	GitCommit = "abc123"

	versionString := "AppVersion=" + AppVersion + ", GitCommit=" + GitCommit
	expected := "AppVersion=1.0.0, GitCommit=abc123"

	if versionString != expected {
		t.Errorf("Expected version string to be '%s', got '%s'", expected, versionString)
	}
}
