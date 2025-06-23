package models

type Tree struct {
	Root *Node `json:"root"`
}

func NewTree() *Tree {
	return &Tree{
		Root: NewNode("root", []string{"root"}),
	}
}

func (t *Tree) PostOrderTraversal(node *Node) []Entry {
	var result []Entry

	if node == nil {
		return result
	}

	for _, node := range node.Children {
		result = append(result, t.PostOrderTraversal(node)...)
	}
	for _, child := range node.Values {
		result = append(result, *NewEntry(child, EntryTypeProject))
	}
	result = append(result, *NewEntry(node.Id, EntryTypeFolder))

	return result
}

func (t *Tree) Print() {
	t.Root.Print()
}

