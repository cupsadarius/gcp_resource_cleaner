package models

import "fmt"

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

func (n *Node) Print() {
	fmt.Println(n.Id)
	fmt.Println(n.Values)
	for _, child := range n.Children {
		child.Print()
	}
}
