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

func TestGetFolders_Success(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("folder1\nfolder2\nfolder3\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	folders, err := GetFolders(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"folder1", "folder2", "folder3"}
	if len(folders) != len(expected) {
		t.Errorf("Expected %d folders, got %d", len(expected), len(folders))
	}

	for i, folder := range folders {
		if folder != expected[i] {
			t.Errorf("Expected folder[%d] to be %s, got %s", i, expected[i], folder)
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

	expectedArgs := []string{"resource-manager", "folders", "list", "--folder", "12345", "--format", "csv[no-heading](ID)"}
	if len(lastCall.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
	}

	for i, arg := range lastCall.Args {
		if arg != expectedArgs[i] {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expectedArgs[i], arg)
		}
	}
}

func TestGetFolders_EmptyResult(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  nil,
	}

	ctx := context.Background()
	folders, err := GetFolders(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if folders != nil {
		t.Errorf("Expected nil folders for empty output, got %v", folders)
	}

	// Verify command was still called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call even for empty result, got %d", mockExec.GetCallCount())
	}
}

func TestGetFolders_CommandError(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  errors.New("gcloud command failed"),
	}

	ctx := context.Background()
	folders, err := GetFolders(ctx, "12345", mockExec)

	if err == nil {
		t.Error("Expected error when command fails, got nil")
	}

	if folders != nil {
		t.Errorf("Expected nil folders when command fails, got %v", folders)
	}

	// Verify command was attempted
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call attempt, got %d", mockExec.GetCallCount())
	}
}

func TestGetFolders_WithNewlines(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("folder1\n\nfolder2\n\n\nfolder3\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	folders, err := GetFolders(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"folder1", "folder2", "folder3"}
	if len(folders) != len(expected) {
		t.Errorf("Expected %d folders, got %d", len(expected), len(folders))
	}

	for i, folder := range folders {
		if folder != expected[i] {
			t.Errorf("Expected folder[%d] to be %s, got %s", i, expected[i], folder)
		}
	}
}

func TestGetFolders_SingleFolder(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("single-folder\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	folders, err := GetFolders(ctx, "12345", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"single-folder"}
	if len(folders) != 1 {
		t.Errorf("Expected 1 folder, got %d", len(folders))
	}

	if folders[0] != expected[0] {
		t.Errorf("Expected folder to be %s, got %s", expected[0], folders[0])
	}
}

func TestGetFolders_DifferentParentFolder(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("child-folder\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	_, err := GetFolders(ctx, "parent-folder-xyz", mockExec)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the parent folder ID was passed correctly
	lastCall := mockExec.GetLastCall()
	folderArgIndex := -1
	for i, arg := range lastCall.Args {
		if arg == "--folder" && i+1 < len(lastCall.Args) {
			folderArgIndex = i + 1
			break
		}
	}

	if folderArgIndex == -1 {
		t.Error("Expected --folder argument not found")
	} else if lastCall.Args[folderArgIndex] != "parent-folder-xyz" {
		t.Errorf("Expected parent folder to be 'parent-folder-xyz', got %s", lastCall.Args[folderArgIndex])
	}
}

func TestDeleteFolder_DryRun(t *testing.T) {
	mockExec := &MockExecutor{}

	ctx := context.Background()
	err := DeleteFolder(ctx, "test-folder", true, mockExec)

	if err != nil {
		t.Errorf("Expected no error in dry run mode, got %v", err)
	}

	// In dry run mode, no command should be executed
	if mockExec.GetCallCount() != 0 {
		t.Errorf("Expected 0 command calls in dry run mode, got %d", mockExec.GetCallCount())
	}
}

func TestDeleteFolder_Success(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Folder deleted successfully\n"),
		MockError:  nil,
	}

	ctx := context.Background()
	err := DeleteFolder(ctx, "test-folder", false, mockExec)

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

	expectedArgs := []string{"resource-manager", "folders", "delete", "test-folder", "--quiet"}
	if len(lastCall.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
	}

	for i, arg := range lastCall.Args {
		if arg != expectedArgs[i] {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expectedArgs[i], arg)
		}
	}
}

func TestDeleteFolder_SuccessEmptyOutput(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte(""),
		MockError:  nil,
	}

	ctx := context.Background()
	err := DeleteFolder(ctx, "test-folder", false, mockExec)

	if err != nil {
		t.Errorf("Expected no error for empty output, got %v", err)
	}

	// Verify command was called
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call, got %d", mockExec.GetCallCount())
	}
}

func TestDeleteFolder_CommandError(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("Error deleting folder\n"),
		MockError:  errors.New("delete command failed"),
	}

	ctx := context.Background()
	err := DeleteFolder(ctx, "test-folder", false, mockExec)

	if err == nil {
		t.Error("Expected error when delete command fails, got nil")
	}

	// Verify command was attempted
	if mockExec.GetCallCount() != 1 {
		t.Errorf("Expected 1 command call attempt, got %d", mockExec.GetCallCount())
	}
}

func TestDeleteFolder_PermissionDenied(t *testing.T) {
	mockExec := &MockExecutor{
		MockOutput: []byte("ERROR: (gcloud.resource-manager.folders.delete) User does not have permission to access folder"),
		MockError:  errors.New("permission denied"),
	}

	ctx := context.Background()
	err := DeleteFolder(ctx, "test-folder", false, mockExec)

	if err == nil {
		t.Error("Expected permission error, got nil")
	}

	if err.Error() != "permission denied" {
		t.Errorf("Expected 'permission denied' error, got %v", err)
	}
}
