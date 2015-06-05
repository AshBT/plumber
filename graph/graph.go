package graph

import (
	"errors"
)

type Node struct {
	Name     string
	children []*Node
	dirty    bool
	visited  bool
}

// Creates a new node with the given name
func NewNode(name string) *Node {
	n := Node{name, nil, false, false}
	return &n
}

// Trajan's algorithm for a topological sort
// will return a list of node names in *reverse* topologically sorted
// order
func ReverseTopoSort(nodes []*Node) ([]string, error) {
	var err error
	// creates a slide with pre-allocated storage
	topoSorted := make([]string, 0, len(nodes))

	for _, node := range nodes {
		if !node.visited {
			err = node.visit(&topoSorted)
			if err != nil {
				return nil, err
			}
		}
	}
	return topoSorted, nil
}

func (n *Node) AddChildren(children... *Node) {
	n.children = append(n.children, children...)
}

// TODO: make this accessor read-only somehow
// func (n *Node) Children() []*Node {
// 	return n.children
// }

func (n *Node) visit(result *[]string) error {
	if n.dirty {
		return errors.New("Pipeline is not a DAG!")
	}
	if !n.visited {
		n.dirty = true
		for _, child := range n.children {
			if err := child.visit(result); err != nil {
				return err
			}
		}
		n.visited = true
		n.dirty = false
		*result = append(*result, n.Name)
	}
	return nil
}
