package models

import (
	"testing"
)

func TestNewEntry(t *testing.T) {
	tests := []struct {
		testName  string
		id        string
		name      string
		entryType EntryType
		expected  *Entry
	}{
		{
			testName:  "create project entry",
			id:        "test-project-123",
			name:      "Test Project 123",
			entryType: EntryTypeProject,
			expected: &Entry{
				Type: EntryTypeProject,
				Id:   "test-project-123",
				Name: "Test Project 123",
			},
		},
		{
			testName:  "create folder entry",
			id:        "test-folder-456",
			name:      "Test Folder 456",
			entryType: EntryTypeFolder,
			expected: &Entry{
				Type: EntryTypeFolder,
				Id:   "test-folder-456",
				Name: "Test Folder 456",
			},
		},
		{
			testName:  "empty id",
			id:        "",
			name:      "",
			entryType: EntryTypeProject,
			expected: &Entry{
				Type: EntryTypeProject,
				Id:   "",
				Name: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			entry := NewEntry(tt.id, tt.name, tt.entryType)

			if entry.Type != tt.expected.Type {
				t.Errorf("Expected type %d, got %d", tt.expected.Type, entry.Type)
			}

			if entry.Id != tt.expected.Id {
				t.Errorf("Expected id %s, got %s", tt.expected.Id, entry.Id)
			}
			if entry.Name != tt.expected.Name {
				t.Errorf("Expected id %s, got %s", tt.expected.Name, entry.Name)
			}
		})
	}
}

func TestEntryTypes_Constants(t *testing.T) {
	tests := []struct {
		entryType EntryType
		expected  string
	}{
		{EntryTypeProject, "project"},
		{EntryTypeFolder, "folder"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if value, exists := EntryTypes[tt.entryType]; !exists {
				t.Errorf("EntryType %d not found in EntryTypes map", tt.entryType)
			} else if value != tt.expected {
				t.Errorf("Expected EntryTypes[%d] to be %s, got %s", tt.entryType, tt.expected, value)
			}
		})
	}
}

func TestEntryType_Values(t *testing.T) {
	// Test that the EntryType constants have expected values
	if EntryTypeProject != 0 {
		t.Errorf("Expected EntryTypeProject to be 0, got %d", EntryTypeProject)
	}

	if EntryTypeFolder != 1 {
		t.Errorf("Expected EntryTypeFolder to be 1, got %d", EntryTypeFolder)
	}
}
