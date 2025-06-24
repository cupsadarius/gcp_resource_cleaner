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

	// NewTree creates an empty tree, Root should be nil initially
	if tree.Root != nil {
		t.Error("NewTree() should create tree with nil root")
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
	node := NewNode(
		NewEntry("folder1", "Folder 1", EntryTypeFolder),
		[]Entry{
			*NewEntry("project1", "Project 1", EntryTypeProject),
			*NewEntry("project2", "Project 2", EntryTypeProject),
		},
	)

	result := tree.PostOrderTraversal(node)

	expected := []Entry{
		{Type: EntryTypeProject, Id: "project1", Name: "Project 1"},
		{Type: EntryTypeProject, Id: "project2", Name: "Project 2"},
		{Type: EntryTypeFolder, Id: "folder1", Name: "Folder 1"},
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

	folder3 := NewNode(
		NewEntry("folder3", "Folder 3", EntryTypeFolder),
		[]Entry{*NewEntry("proj3", "Project 3", EntryTypeProject)},
	)
	folder2 := NewNode(
		NewEntry("folder2", "Folder 2", EntryTypeFolder),
		[]Entry{},
	)
	folder2.Children = append(folder2.Children, folder3)

	folder1 := NewNode(
		NewEntry("folder1", "Folder 1", EntryTypeFolder),
		[]Entry{
			*NewEntry("proj1", "Project 1", EntryTypeProject),
			*NewEntry("proj2", "Project 2", EntryTypeProject),
		},
	)

	root := NewNode(
		NewEntry("root", "Root", EntryTypeFolder),
		[]Entry{},
	)
	root.Children = append(root.Children, folder1, folder2)

	result := tree.PostOrderTraversal(root)

	// Post-order should visit: proj1, proj2, folder1, proj3, folder3, folder2, root
	expected := []Entry{
		{Type: EntryTypeProject, Id: "proj1", Name: "Project 1"},
		{Type: EntryTypeProject, Id: "proj2", Name: "Project 2"},
		{Type: EntryTypeFolder, Id: "folder1", Name: "Folder 1"},
		{Type: EntryTypeProject, Id: "proj3", Name: "Project 3"},
		{Type: EntryTypeFolder, Id: "folder3", Name: "Folder 3"},
		{Type: EntryTypeFolder, Id: "folder2", Name: "Folder 2"},
		{Type: EntryTypeFolder, Id: "root", Name: "Root"},
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d entries, got %d", len(expected), len(result))
	}

	for i, entry := range result {
		if entry.Type != expected[i].Type || entry.Id != expected[i].Id || entry.Name != expected[i].Name {
			t.Errorf("Entry %d: expected %+v, got %+v", i, expected[i], entry)
		}
	}
}

func TestPostOrderTraversal_NoProjects(t *testing.T) {
	tree := NewTree()

	// Tree with only folders, no projects
	folder2 := NewNode(NewEntry("folder2", "Folder 2", EntryTypeFolder), []Entry{})
	folder1 := NewNode(NewEntry("folder1", "Folder 1", EntryTypeFolder), []Entry{})
	folder1.Children = append(folder1.Children, folder2)

	result := tree.PostOrderTraversal(folder1)

	expected := []Entry{
		{Type: EntryTypeFolder, Id: "folder2", Name: "Folder 2"},
		{Type: EntryTypeFolder, Id: "folder1", Name: "Folder 1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestPostOrderTraversal_DeepNesting(t *testing.T) {
	tree := NewTree()

	// Create deeply nested structure
	folder3 := NewNode(NewEntry("folder3", "Folder 3", EntryTypeFolder), []Entry{*NewEntry("proj3", "Project 3", EntryTypeProject)})
	folder2 := NewNode(NewEntry("folder2", "Folder 2", EntryTypeFolder), []Entry{*NewEntry("proj2", "Project 2", EntryTypeProject)})
	folder2.Children = append(folder2.Children, folder3)
	folder1 := NewNode(NewEntry("folder1", "Folder 1", EntryTypeFolder), []Entry{*NewEntry("proj1", "Project 1", EntryTypeProject)})
	folder1.Children = append(folder1.Children, folder2)

	result := tree.PostOrderTraversal(folder1)

	expected := []Entry{
		{Type: EntryTypeProject, Id: "proj3", Name: "Project 3"},
		{Type: EntryTypeFolder, Id: "folder3", Name: "Folder 3"},
		{Type: EntryTypeProject, Id: "proj2", Name: "Project 2"},
		{Type: EntryTypeFolder, Id: "folder2", Name: "Folder 2"},
		{Type: EntryTypeProject, Id: "proj1", Name: "Project 1"},
		{Type: EntryTypeFolder, Id: "folder1", Name: "Folder 1"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
