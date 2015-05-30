package graph

import (
	"errors"
)

type Node struct {
	Name string
	Children []*Node
	dirty bool
	visited bool
}

func NewNode(name string) (*Node) {
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

func (n *Node) visit(result *[]string) error {
	if n.dirty {
		return errors.New("Pipeline is not a DAG!")
	}
	if !n.visited {
		n.dirty = true
		for _, child := range n.Children {
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
