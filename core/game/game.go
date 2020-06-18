package game

import (
    "strings"
    "strconv"
    "errors"
)

// Path is used for tree-pathing
type Path string

// Game contains the game tree information
type Game struct {
	Root *Node
}

// New creates a Game struct.
func New() *Game {
	g := &Game{
		Root: NewNode(),
	}
	return g
}

// TreePath returns the node found at a given path
func (g *Game) TreePath(p Path) (*Node, error) {
    treePath, err := parseTreePath(p)
    if err != nil {
        return nil, err
    }
    if len(treePath) == 0 {
        return nil, errors.New("Empty tree path")
    }
    current := g.Root.Children[0]
    for i := 0; i < treePath[0]; i++ {
        current = current.Children[0]
    }
    if len(treePath) > 1 {
        for i := 1; i < len(treePath); i++ {
            current = current.Children[treePath[i]]
        }
    }

    return current, nil
}

// parseTreePath currently only handles the simplest possible paths
// will fix in another branch
func parseTreePath(p Path) ([]int, error) {
    tokens := strings.Split(string(p), ".")
    result := make([]int, len(tokens))
    for i,t := range tokens {
        n, err := strconv.Atoi(t)
        if err != nil {
            return []int{}, nil
        }
        result[i] = n
    }
    return result, nil
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
