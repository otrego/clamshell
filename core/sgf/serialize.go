package sgf

import (
	"fmt"
	"strings"

	"github.com/otrego/clamshell/core/game"
)

// Serialize converts a Game into SGF format.
// Calls serializeHelper.
func Serialize(g *game.Game) string {
	return fmt.Sprintf("(%s)", serializeHelper(g.Root))
}

// serializeHelper is a recursive DFS searching all
// descendant nodes of n.
func serializeHelper(n *game.Node) string {
	var sb strings.Builder
	sb.WriteString(writeNode(n))
	for _, child := range n.Children {
		if len(n.Children) > 1 {
			sb.WriteString(fmt.Sprintf("(%s)", serializeHelper(child)))
		} else {
			sb.WriteString(serializeHelper(child))
		}
	}
	return sb.String()
}

// writeNode writes a node in SGF format
func writeNode(n *game.Node) string {
	var sb strings.Builder

	sb.WriteString(";")
	for key := range n.Properties {
		sb.WriteString(key)
		for _, value := range n.Properties[key] {
			sb.WriteString(fmt.Sprintf("[%s]", value))
		}
	}
	return sb.String()
}
