package sgf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/otrego/clamshell/core/movetree"
)

// Serialize converts a Game into SGF format.
// Calls serializeHelper.
func Serialize(g *movetree.MoveTree) string {
	return fmt.Sprintf("(%s)", serializeHelper(g.Root))
}

// serializeHelper is a recursive DFS searching all
// descendant nodes of n.
func serializeHelper(n *movetree.Node) string {
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
func writeNode(n *movetree.Node) string {
	var sb strings.Builder

	sb.WriteString(";")
	keys := make([]string, 0)
	for key := range n.Properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for key := range n.Properties {
		sb.WriteString(key)
		for _, value := range n.Properties[key] {
			sb.WriteString(fmt.Sprintf("[%s]", value))
		}
	}
	return sb.String()
}
