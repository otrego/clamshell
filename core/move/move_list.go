package move

import "sort"

// MoveList is a convenience helper for a list of moves.
type MoveList []*Move

// Sort the list of moves (in-place), in order to have a well-ordered
// collection of moves.
//
// The rules:
// - Black moves are sorted before White and then Empty (although we don't
//   generally allow empty moves.
// - Y values are compared next
// - X values are compared last
func (m MoveList) Sort() {
	sort.Slice(m, func(i, j int) bool {
		if m[i].Color() != m[j].Color() {
			return m[i].Color().Ordinal() < m[j].Color().Ordinal()
		}
		if m[i].Point().Y() != m[j].Point().Y() {
			return m[i].Point().Y() < m[j].Point().Y()
		}
		return m[i].Point().X() < m[j].Point().X()
	})
}

// String returns the string for of the MoveList.
func (m MoveList) String() string {
	var sb strings.StringBuilder
	for _, m := range {
		sb
	}
}
