package game

import (
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

// A Move is location + a color. A pass is represented by a Move with defined
// color but no point specified.
type Move struct {
	Color color.Color
	Point *point.Point
}

// Game contains the game tree information for a go game.
type Game struct {
	Root *Node
}

// New creates a Game.
func New() *Game {
	g := &Game{
		Root: NewNode(),
	}
	return g
}

// Node contains Properties, Children nodes, and Parent node.
type Node struct {
	// Placements are stones that are used for setup, but actual moves. For
	// example, handicap stones will be in in placements.
	Placements []*Move

	// Move indicates a move in the game.
	Move *Move

	// Properties contain all the raw/unprocessed properties
	Properties map[string][]string

	// Children of this position.
	Children []*Node

	//
	Parent *Node
}

// NewNode creates a Node.
func NewNode() *Node {
	return &Node{
		Properties: make(map[string][]string),
	}
}

// AddChild adds a child node.
func (n *Node) AddChild(nn *Node) {
	n.Children = append(n.Children, nn)
}
