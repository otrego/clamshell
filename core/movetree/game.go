package movetree

// MoveTree contains the game tree information for a go game.
type MoveTree struct {
	Root *Node
}

// New creates a MoveTree.
func New() *MoveTree {
	g := &MoveTree{
		Root: NewNode(),
	}
	return g
}
