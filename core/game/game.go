package game

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
