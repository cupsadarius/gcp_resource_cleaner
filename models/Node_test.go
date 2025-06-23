package models

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		values   []string
		expected *Node
	}{
		{
			name:   "basic node creation",
			id:     "test-folder",
			values: []string{"project1", "project2"},
			expected: &Node{
				Id:       "test-folder",
				Values:   []string{"project1", "project2"},
				Children: make([]*Node, 0),
			},
		},
		{
			name:   "empty values",
			id:     "empty-folder",
			values: []string{},
			expected: &Node{
				Id:       "empty-folder",
				Values:   []string{},
				Children: make([]*Node, 0),
			},
		},
		{
			name:   "nil values",
			id:     "nil-folder",
			values: nil,
			expected: &Node{
				Id:       "nil-folder",
				Values:   nil,
				Children: make([]*Node, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNode(tt.id, tt.values)

			if node.Id != tt.expected.Id {
				t.Errorf("Expected id %s, got %s", tt.expected.Id, node.Id)
			}

			if len(node.Values) != len(tt.expected.Values) {
				t.Errorf("Expected values length %d, got %d", len(tt.expected.Values), len(node.Values))
			}

			for i, value := range node.Values {
				if value != tt.expected.Values[i] {
					t.Errorf("Expected value[%d] to be %s, got %s", i, tt.expected.Values[i], value)
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
	parent := NewNode("parent", []string{"project1"})
	child1 := NewNode("child1", []string{"project2"})
	child2 := NewNode("child2", []string{"project3"})

	parent.Children = append(parent.Children, child1, child2)

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].Id != "child1" {
		t.Errorf("Expected first child id to be 'child1', got %s", parent.Children[0].Id)
	}

	if parent.Children[1].Id != "child2" {
		t.Errorf("Expected second child id to be 'child2', got %s", parent.Children[1].Id)
	}
}
