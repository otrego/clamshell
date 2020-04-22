package game

// Game contains the game tree information
type Game struct {
	Root *Node
}

// NewGame creates a Game struct
func NewGame() *Game {
	g := &Game{
		Root: NewNode(),
	}
	return g
}

// Node contains Properties, Children nodes, and Parent node
type Node struct {
	Properties map[string][]string
	Children   []*Node
	Parent     *Node
}

// NewNode creates a Node
func NewNode() *Node {
	return &Node{
		Properties: make(map[string][]string),
	}
}
