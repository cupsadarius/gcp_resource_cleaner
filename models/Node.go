package models

import (
	"fmt"

	"github.com/xlab/treeprint"
)

type Node struct {
	Current  *Entry
	Values   []Entry
	Children []*Node
}

func NewNode(current *Entry, values []Entry) *Node {
	return &Node{
		Current:  current,
		Values:   values,
		Children: make([]*Node, 0),
	}
}

func (n *Node) Print(node treeprint.Tree) {
	folder := node.AddBranch(fmt.Sprintf("%s (%s)", n.Current.Name, n.Current.Id))
	for _, value := range n.Values {
		folder.AddNode(fmt.Sprintf("%s (%s)", value.Name, value.Id))
	}
	for _, child := range n.Children {
		child.Print(folder)
	}
}
