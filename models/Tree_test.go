package models

import (
	"reflect"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := NewTree()

	if tree == nil {
		t.Fatal("NewTree() returned nil")
	}

	if tree.Root == nil {
		t.Fatal("NewTree() created tree with nil root")
	}

	if tree.Root.Id != "root" {
		t.Errorf("Expected root id to be 'root', got %s", tree.Root.Id)
	}

	if len(tree.Root.Values) != 1 || tree.Root.Values[0] != "root" {
		t.Errorf("Expected root values to be ['root'], got %v", tree.Root.Values)
	}
}

func TestPostOrderTraversal_EmptyTree(t *testing.T) {
	tree := NewTree()

	result := tree.PostOrderTraversal(nil)

	if len(result) != 0 {
		t.Errorf("Expected empty result for nil node, got %d entries", len(result))
	}
}

func TestPostOrderTraversal_SingleNode(t *testing.T) {
	tree := NewTree()
	node := NewNode("folder1", []string{"project1", "project2"})

	result := tree.PostOrderTraversal(node)

	expected := []Entry{
		{Type: EntryTypeProject, Id: "project1"},
		{Type: EntryTypeProject, Id: "project2"},
		{Type: EntryTypeFolder, Id: "folder1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestPostOrderTraversal_MultiLevel(t *testing.T) {
	tree := NewTree()

	// Create a tree structure:
	//       root
	//      /    \
	//   folder1  folder2
	//   /   \       \
	// proj1 proj2  folder3
	//                 |
	//              proj3

	folder3 := NewNode("folder3", []string{"proj3"})
	folder2 := NewNode("folder2", []string{})
	folder2.Children = append(folder2.Children, folder3)

	folder1 := NewNode("folder1", []string{"proj1", "proj2"})

	root := NewNode("root", []string{})
	root.Children = append(root.Children, folder1, folder2)

	result := tree.PostOrderTraversal(root)

	// Post-order should visit: proj1, proj2, folder1, proj3, folder3, folder2, root
	expected := []Entry{
		{Type: EntryTypeProject, Id: "proj1"},
		{Type: EntryTypeProject, Id: "proj2"},
		{Type: EntryTypeFolder, Id: "folder1"},
		{Type: EntryTypeProject, Id: "proj3"},
		{Type: EntryTypeFolder, Id: "folder3"},
		{Type: EntryTypeFolder, Id: "folder2"},
		{Type: EntryTypeFolder, Id: "root"},
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d entries, got %d", len(expected), len(result))
	}

	for i, entry := range result {
		if entry.Type != expected[i].Type || entry.Id != expected[i].Id {
			t.Errorf("Entry %d: expected %+v, got %+v", i, expected[i], entry)
		}
	}
}

func TestPostOrderTraversal_NoProjects(t *testing.T) {
	tree := NewTree()

	// Tree with only folders, no projects
	folder2 := NewNode("folder2", []string{})
	folder1 := NewNode("folder1", []string{})
	folder1.Children = append(folder1.Children, folder2)

	result := tree.PostOrderTraversal(folder1)

	expected := []Entry{
		{Type: EntryTypeFolder, Id: "folder2"},
		{Type: EntryTypeFolder, Id: "folder1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestPostOrderTraversal_DeepNesting(t *testing.T) {
	tree := NewTree()

	// Create deeply nested structure
	folder3 := NewNode("folder3", []string{"proj3"})
	folder2 := NewNode("folder2", []string{"proj2"})
	folder2.Children = append(folder2.Children, folder3)
	folder1 := NewNode("folder1", []string{"proj1"})
	folder1.Children = append(folder1.Children, folder2)

	result := tree.PostOrderTraversal(folder1)

	expected := []Entry{
		{Type: EntryTypeProject, Id: "proj3"},
		{Type: EntryTypeFolder, Id: "folder3"},
		{Type: EntryTypeProject, Id: "proj2"},
		{Type: EntryTypeFolder, Id: "folder2"},
		{Type: EntryTypeProject, Id: "proj1"},
		{Type: EntryTypeFolder, Id: "folder1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
