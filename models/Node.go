package models

import (
	"github.com/xlab/treeprint"
)

type Node struct {
	Id       string   `json:"id"`
	Values   []string `json:"values"`
	Children []*Node  `json:"children"`
}

func NewNode(id string, values []string) *Node {
	return &Node{
		Id:       id,
		Values:   values,
		Children: make([]*Node, 0),
	}
}

func (n *Node) Print(node treeprint.Tree) {
	folder := node.AddBranch(n.Id)
	for _, value := range n.Values {
		folder.AddNode(value)
	}
	for _, child := range n.Children {
		child.Print(folder)
	}
}
