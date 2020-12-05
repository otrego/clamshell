package movetree

import (
	"github.com/otrego/clamshell/core/move"
)

// Node contains Properties, Children nodes, and Parent node.
type Node struct {
	// moveNum is the move and indicates the current move number or depth for this
	// node. The root node should always be 0.
	moveNum int

	// varNum is the current variation number. The primary variation (or mainline
	// variation) should always be 0.
	varNum int

	// Placements are stones that are used for setup, but actual moves. For
	// example, handicap stones will be in in placements.
	Placements []*move.Move

	// Move indicates a move in the game.
	Move *move.Move

	// Properties contain all the raw/unprocessed properties
	Properties map[string][]string

	// Children of this position.
	Children []*Node

	// Parent of this node.
	Parent *Node

	// analysisData contains arbitrary AnalysisData
	analysisData interface{}
}

// NewNode creates a Node.
func NewNode() *Node {
	return &Node{
		Properties: make(map[string][]string),
	}
}

// AddChild adds a child node.
func (n *Node) AddChild(nn *Node) {
	nn.moveNum = n.moveNum + 1
	nn.varNum = len(n.Children)
	n.Children = append(n.Children, nn)
}

// Next gets the next node, given the variation number, returning nil if no node
// is available.
func (n *Node) Next(variation int) *Node {
	// We assume there are no gaps in the Children slice.
	if variation < len(n.Children) {
		return n.Children[variation]
	}
	return nil
}

// MoveNum returns the current move number.
func (n *Node) MoveNum() int {
	return n.moveNum
}

// VarNum returns the variation number.
func (n *Node) VarNum() int {
	return n.varNum
}

// SetAnalysisData sets the analysis data.
func (n *Node) SetAnalysisData(an interface{}) {
	n.analysisData = an
}

// AnalysisData gets the attached analysis data, returning nil if the
// analysisData is empty.
func (n *Node) AnalysisData() interface{} {
	return n.analysisData
}

// Traverse Traverses the tree using BFS.
func (n *Node) Traverse(fn func(node *Node)) {
	queue := make([]*Node, 0)
	queue = append(queue, n)
	for len(queue) != 0 {
		value := queue[0]
		fn(value)
		queue = queue[1:]                        // dequeue
		queue = append(queue, value.Children...) // enqueue
	}
}

// TraverseMainBranch Traverses the 0th variation nodes (Main Branch).
// This Traversal uses BFS.
func (n *Node) TraverseMainBranch(fn func(node *Node)) {
	queue := make([]*Node, 0)
	queue = append(queue, n)
	for len(queue) != 0 {
		value := queue[0]
		queue = queue[1:] // dequeue
		if value.varNum == 0 {
			fn(value)
			queue = append(queue, value.Children...) // enqueue
		}
	}
}
