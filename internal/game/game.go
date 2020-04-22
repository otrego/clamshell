package game

// A game record (sgf) struct for now is just a root node
type Game struct {
	Root *Node
}

func NewGame() *Game {
	g := &Game{
		Root: NewNode(),
	}
	return g
}

// A node is comprised of key-value pairs (might hold other info someday)
// keys can be repeated
type Node struct {
	Properties map[string][]string
	Children   []*Node
	Parent     *Node
}

func NewNode() *Node {
	return &Node{
		Properties: make(map[string][]string),
	}
}
