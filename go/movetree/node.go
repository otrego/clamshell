package movetree

import (
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
)

// GameInfo contains typed game properties that can exist only on the root.
type GameInfo struct {
	// Size of the board, where 19 = 19x19. Between 1 and 25 inclusive. A value of
	// 0 should be taken to mean 'unspecified' and treated as 19x19.
	Size int

	// Komi are points added to the player with the white stones as compensation for playing second.
	// Komi must have a decimal value of .0 or .5 (ex: 6.5)
	Komi *float64

	// Initial player turn. This is traditionally the player with the black stones
	Player color.Color
}

// Node contains Properties, Children nodes, and Parent node.
type Node struct {
	// moveNum is the move and indicates the current move number or depth for this
	// node. The root node should always be 0.
	moveNum int

	// varNum is the current variation number. The primary variation (or mainline
	// variation) should always be 0.
	varNum int

	// Children of this position.
	Children []*Node

	// Parent of this node.
	Parent *Node

	// Move indicates a move in the game. A move with a color + no intersection is
	// used to indicate a pass.
	Move *move.Move

	// Placements are stones that are used for setup, but actual moves. For
	// example, handicap stones will be in in placements.
	Placements move.List

	// Comment is the comment for the current node.
	Comment string

	// GameInfo contains properties only found on the root. Should be nil on
	// non-root nodes.
	GameInfo *GameInfo

	// SGFProperties contain all the raw/unprocessed properties
	SGFProperties map[string][]string

	// analysisData contains arbitrary/untyped AnalysisData that is attached to
	// this node.
	analysisData interface{}
}

// NewNode creates a Node.
func NewNode() *Node {
	return &Node{
		SGFProperties: make(map[string][]string),
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
