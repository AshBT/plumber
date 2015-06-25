/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

func (n *Node) AddChildren(children ...*Node) {
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
