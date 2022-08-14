package sgf

import (
	"strings"

	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/prop"
)

// Serialize converts a Game into SGF format.
// Calls serializeHelper.
func Serialize(g *movetree.MoveTree) (string, error) {
	s, err := serializeHelper(g.Root)
	if err != nil {
		return "", err
	}
	return "(" + s + ")", nil
}

// serializeHelper is a recursive DFS searching all
// descendant nodes of n.
func serializeHelper(n *movetree.Node) (string, error) {
	var sb strings.Builder
	s, err := writeNode(n)
	if err != nil {
		return "", nil
	}
	sb.WriteString(s)

	for _, child := range n.Children {
		s, err := serializeHelper(child)
		if err != nil {
			return "", err
		}
		if len(n.Children) > 1 {
			sb.WriteString("(" + s + ")")
		} else {
			sb.WriteString(s)
		}
	}
	return sb.String(), nil
}

// writeNode writes a node in SGF format
func writeNode(n *movetree.Node) (string, error) {
	s, err := prop.ConvertNode(n)
	if err != nil {
		return s, err
	}
	return ";" + s, nil
}
