// Package symbol contains an inventory of intersection symbol.
package symbol

// Symbol represents an intersection on the board.
type Symbol int

const (
	Empty Symbol = 0

	TopLeft     Symbol = 10
	TopRight    Symbol = 11
	BottomLeft  Symbol = 12
	BottomRight Symbol = 13

	TopEdge    Symbol = 14
	BottomEdge Symbol = 15
	LeftEdge   Symbol = 16
	RightEdge  Symbol = 17
	Center     Symbol = 18
	Starpoint  Symbol = 19

	BlackStone Symbol = 20
	WhiteStone Symbol = 21

	Triangle Symbol = 30
	Square   Symbol = 31
	Circle   Symbol = 32
	Xmark    Symbol = 33

	TextLabel Symbol = 34
)

// UnicodeChar converts a symbol representation to a single char. This method is
// primarily for debugging.
func (s Symbol) UnicodeChar() string {
	switch s {
	case TopLeft:
		return "┏"
	case TopRight:
		return "┓"
	case BottomLeft:
		return "┗"
	case BottomRight:
		return "┛"
	case TopEdge:
		return "┳"
	case BottomEdge:
		return "┻"
	case LeftEdge:
		return "┣"
	case RightEdge:
		return "┫"
	case Center:
		return "╋"
	case Starpoint:
		return "✻"
	case BlackStone:
		return "●"
	case WhiteStone:
		return "○"
	case Triangle:
		return "▴"
	case Square:
		return "□"
	case Circle:
		return "⊙"
	case Xmark:
		return "☓"
	case TextLabel:
		return "☒"
	}
	return "?"
}
