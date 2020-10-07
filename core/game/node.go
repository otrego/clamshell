package game

// Node contains Properties, Children nodes, and Parent node.
type Node struct {
	// moveNum is the move and indicates the current move number or depth for this
	// node. The root node should always be 0.
	moveNum int

	// Placements are stones that are used for setup, but actual moves. For
	// example, handicap stones will be in in placements.
	Placements []*Move

	// Move indicates a move in the game.
	Move *Move

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
	n.Children = append(n.Children, nn)
}

//
func (n *Node) SetAnalysisData() {

}
