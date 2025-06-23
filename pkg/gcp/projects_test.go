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

func TestGetProjects_Success(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("project1\nproject2\nproject3\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	projects, err := GetProjects(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"project1", "project2", "project3"}
	if len(projects) != len(expected) {
		t.Errorf("Expected %d projects, got %d", len(expected), len(projects))
	}

	for i, project := range projects {
		if project != expected[i] {
			t.Errorf("Expected project[%d] to be %s, got %s", i, expected[i], project)
		}
	}

	// Verify the correct command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}

	expectedArgs := []string{"projects", "list", "--filter", "parent.id:12345", "--format", "csv[no-heading](projectId)"}
	if len(lastCall.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
	}
}

func TestGetProjects_EmptyResult(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  nil,
	}

	ctx := context.Background()
	projects, err := GetProjects(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if projects != nil {
		t.Errorf("Expected nil projects for empty output, got %v", projects)
	}
}

func TestGetProjects_CommandError(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  errors.New("command failed"),
	}

	ctx := context.Background()
	projects, err := GetProjects(ctx, "12345", mockExec)

	if err == nil {
		t.Error("Expected error when command fails, got nil")
	}

	if projects != nil {
		t.Errorf("Expected nil projects when command fails, got %v", projects)
	}
}

func TestGetProjects_WithNewlines(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("project1\n\nproject2\n\n\nproject3\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	projects, err := GetProjects(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"project1", "project2", "project3"}
	if len(projects) != len(expected) {
		t.Errorf("Expected %d projects, got %d", len(expected), len(projects))
	}

	for i, project := range projects {
		if project != expected[i] {
			t.Errorf("Expected project[%d] to be %s, got %s", i, expected[i], project)
		}
	}
}

func TestDeleteProject_DryRun(t *testing.T) {
	mockExec := &MockExecutor{}

	ctx := context.Background()
	err := DeleteProject(ctx, "test-project", true, mockExec)

	if err != nil {
		t.Errorf("Expected no error in dry run mode, got %v", err)
	}

	// In dry run mode, no command should be executed
	if mockExec.GetCallCount() != 0 {
		t.Errorf("Expected 0 command calls in dry run mode, got %d", mockExec.GetCallCount())
	}
}

func TestDeleteProject_Success(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Project deleted successfully\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	err := DeleteProject(ctx, "test-project", false, mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the correct command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}

	lastCall := mockExec.GetLastCall()
	if lastCall.Name != "gcloud" {
		t.Errorf("Expected command to be 'gcloud', got %s", lastCall.Name)
	}

	expectedArgs := []string{"projects", "delete", "test-project", "--quiet"}
	if len(lastCall.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
	}

	for i, arg := range lastCall.Args {
		if arg != expectedArgs[i] {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expectedArgs[i], arg)
		}
	}
}

func TestDeleteProject_CommandError(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Error deleting project\n"),
		MockError:  errors.New("delete command failed"),
	}

	ctx := context.Background()
	err := DeleteProject(ctx, "test-project", false, mockExec)

	if err == nil {
		t.Error("Expected error when delete command fails, got nil")
	}
}
