package models

import (
	"testing"
)

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		entryType EntryType
		expected  *Entry
	}{
		{
			name:      "create project entry",
			id:        "test-project-123",
			entryType: EntryTypeProject,
			expected: &Entry{
				Type: EntryTypeProject,
				Id:   "test-project-123",
			},
		},
		{
			name:      "create folder entry",
			id:        "test-folder-456",
			entryType: EntryTypeFolder,
			expected: &Entry{
				Type: EntryTypeFolder,
				Id:   "test-folder-456",
			},
		},
		{
			name:      "empty id",
			id:        "",
			entryType: EntryTypeProject,
			expected: &Entry{
				Type: EntryTypeProject,
				Id:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := NewEntry(tt.id, tt.entryType)

			if entry.Type != tt.expected.Type {
				t.Errorf("Expected type %d, got %d", tt.expected.Type, entry.Type)
			}

			if entry.Id != tt.expected.Id {
				t.Errorf("Expected id %s, got %s", tt.expected.Id, entry.Id)
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
