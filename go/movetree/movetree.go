// Package movetree contains logic for managing trees of moves in a game or
// problem. A movetree can be serialized to / deserialized from an SGF.
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

	// sensible defaults go here.
	g.Root.GameInfo = &GameInfo{}
	g.Root.GameInfo.Size = 19
	g.Root.SGFProperties["GM"] = []string{"1"}     // GM[1]=go
	g.Root.SGFProperties["FF"] = []string{"4"}     // FF[4]=SGF file format 4
	g.Root.SGFProperties["CA"] = []string{"UTF-8"} // CA[UTF-8]=UTF-8 encoding
	return g
}
