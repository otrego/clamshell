package move

import (
	"sort"
	"strings"
)

// List is a convenience helper for a list of moves.
type List []*Move

// Sort the list of moves (in-place), in order to have a well-ordered
// collection of moves.
//
// The rules:
// - Black moves are sorted before White and then Empty (although we don't
//   generally allow empty moves.
// - X values are compared next
// - Y values are compared last
func (m List) Sort() {
	sort.Slice(m, func(i, j int) bool {
		if m[i].Color() != m[j].Color() {
			return m[i].Color().Ordinal() < m[j].Color().Ordinal()
		}
		if m[i].Point().X() != m[j].Point().X() {
			return m[i].Point().X() < m[j].Point().X()
		}
		return m[i].Point().Y() < m[j].Point().Y()
	})
}

// String returns the string for of the List.
func (m List) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, mv := range m {
		sb.WriteString(mv.String())
		if i < len(m)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
