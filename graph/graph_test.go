package graph

import (
	"testing"
)

// graph with diamond (test case)
// n1 := graph.NewNode("foo")
// n2 := graph.NewNode("bar")
// n3 := graph.NewNode("joe")
// n4 := graph.NewNode("bob")
// n1.Children = append(n1.Children, n2, n3)
// n2.Children = append(n2.Children, n4)
// n3.Children = append(n3.Children, n4)

func TestNodeCreation(t *testing.T) {
	node := NewNode("mynode")

	if node.Name != "mynode" || node.children != nil || node.dirty != false || node.visited != false {
		t.Errorf("Node '%v' not properly initialized", node)
	}
}

func TestAddChilren(t *testing.T) {
	n1 := NewNode("n1")
	n2 := NewNode("n2")

	n1.AddChildren(n2)

	if n1.children[0] != n2 {
		t.Errorf("Wrong children list: got '%v' instead", n1.children)
	}
}

func TestReverseTopoSortEmpty(t *testing.T) {
	sorted, err := ReverseTopoSort(nil)
	if err != nil {
		t.Error(err)
	}
	if len(sorted) != 0 {
		t.Error("sorted list was not empty!")
	}
}

func TestReverseTopoSortSingle(t *testing.T) {
	node := NewNode("one-node-tree")
	sorted, err := ReverseTopoSort([]*Node{node})
	if err != nil {
		t.Error(err)
	}
	if sorted[0] != "one-node-tree" {
		t.Error("Improperly sorted tree!")
	}
}

func TestReverseTopoSortDiamond(t *testing.T) {
	n1 := NewNode("foo")
	n2 := NewNode("bar")
	n3 := NewNode("joe")
	n4 := NewNode("bob")
	n1.AddChildren(n2, n3)
	n2.AddChildren(n4)
	n3.AddChildren(n4)

	sorted, err := ReverseTopoSort([]*Node{n4, n2, n3, n1})
	if err != nil {
		t.Error(err)
	}

	if sorted[0] != "bob" || sorted[1] != "bar" || sorted[2] != "joe" || sorted[3] != "foo" {
		t.Errorf("Improperly reverse-sorted tree: '%v'", sorted)
	}
}

func TestReverseTopoSortAny(t *testing.T) {
	n1 := NewNode("foo")
	n2 := NewNode("bar")
	n3 := NewNode("baz")
	n4 := NewNode("bux")

	sorted, err := ReverseTopoSort([]*Node{n4, n2, n3, n1})
	if err != nil {
		t.Error(err)
	}

	if sorted[0] != "bux" || sorted[1] != "bar" || sorted[2] != "baz" || sorted[3] != "foo" {
		t.Errorf("Improperly reverse-sorted tree: '%v'", sorted)
	}
}

func TestReverseSortedLoopFail(t *testing.T) {
	n1 := NewNode("a")
	n2 := NewNode("b")

	n1.AddChildren(n2)
	n2.AddChildren(n1)

	_, err := ReverseTopoSort([]*Node{n1, n2})
	if err == nil || err.Error() != "Pipeline is not a DAG!" {
		t.Error("Expected a DAG error.")
	}
}

// Benchmark DAG sorting performance
func BenchmarkSort5(b *testing.B) {
	n1 := NewNode("foo")
	n2 := NewNode("bar")
	n3 := NewNode("joe")
	n4 := NewNode("bob")
	n5 := NewNode("baz")
	n1.AddChildren(n2, n3)
	n2.AddChildren(n4)
	n3.AddChildren(n4)
	n4.AddChildren(n5)

	for i := 0; i < b.N; i++ {
		_, err := ReverseTopoSort([]*Node{n5, n4, n2, n3, n1})
		if err != nil {
			b.Error(err)
		}
	}
}
