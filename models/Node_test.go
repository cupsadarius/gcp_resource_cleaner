package models

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	tests := []struct {
		name     string
		current  *Entry
		values   []Entry
		expected *Node
	}{
		{
			name:    "basic node creation",
			current: NewEntry("test-folder", "Test Folder", EntryTypeFolder),
			values: []Entry{
				*NewEntry("project1", "Project 1", EntryTypeProject),
				*NewEntry("project2", "Project 2", EntryTypeProject),
			},
			expected: &Node{
				Current: NewEntry("test-folder", "Test Folder", EntryTypeFolder),
				Values: []Entry{
					*NewEntry("project1", "Project 1", EntryTypeProject),
					*NewEntry("project2", "Project 2", EntryTypeProject),
				},
				Children: make([]*Node, 0),
			},
		},
		{
			name:    "empty values",
			current: NewEntry("empty-folder", "Empty Folder", EntryTypeFolder),
			values:  []Entry{},
			expected: &Node{
				Current:  NewEntry("empty-folder", "Empty Folder", EntryTypeFolder),
				Values:   []Entry{},
				Children: make([]*Node, 0),
			},
		},
		{
			name:    "nil values",
			current: NewEntry("nil-folder", "Nil Folder", EntryTypeFolder),
			values:  nil,
			expected: &Node{
				Current:  NewEntry("nil-folder", "Nil Folder", EntryTypeFolder),
				Values:   nil,
				Children: make([]*Node, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNode(tt.current, tt.values)

			if node.Current.Id != tt.expected.Current.Id {
				t.Errorf("Expected current id %s, got %s", tt.expected.Current.Id, node.Current.Id)
			}

			if node.Current.Name != tt.expected.Current.Name {
				t.Errorf("Expected current name %s, got %s", tt.expected.Current.Name, node.Current.Name)
			}

			if len(node.Values) != len(tt.expected.Values) {
				t.Errorf("Expected values length %d, got %d", len(tt.expected.Values), len(node.Values))
			}

			for i, value := range node.Values {
				if value.Id != tt.expected.Values[i].Id || value.Name != tt.expected.Values[i].Name {
					t.Errorf("Expected value[%d] to be %+v, got %+v", i, tt.expected.Values[i], value)
				}
			}

			if node.Children == nil {
				t.Error("Expected Children to be initialized, got nil")
			}

			if len(node.Children) != 0 {
				t.Errorf("Expected Children to be empty, got %d items", len(node.Children))
			}
		})
	}
}

func TestNode_AddChildren(t *testing.T) {
	parent := NewNode(NewEntry("parent", "Parent", EntryTypeFolder), []Entry{*NewEntry("project1", "Project 1", EntryTypeProject)})
	child1 := NewNode(NewEntry("child1", "Child 1", EntryTypeFolder), []Entry{*NewEntry("project2", "Project 2", EntryTypeProject)})
	child2 := NewNode(NewEntry("child2", "Child 2", EntryTypeFolder), []Entry{*NewEntry("project3", "Project 3", EntryTypeProject)})

	parent.Children = append(parent.Children, child1, child2)

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].Current.Id != "child1" {
		t.Errorf("Expected first child id to be 'child1', got %s", parent.Children[0].Current.Id)
	}

	if parent.Children[1].Current.Id != "child2" {
		t.Errorf("Expected second child id to be 'child2', got %s", parent.Children[1].Current.Id)
	}
}
