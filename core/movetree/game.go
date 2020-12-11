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

	// Sensible Defaults go here.
	g.Root.Properties.Size = 19
	g.Root.SGFProperties["GM"] = []string{"1"} // GM[1]=go
	return g
}
