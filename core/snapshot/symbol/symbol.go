// Package symbol contains an inventory of intersection symbol.
package symbol

// Symbol represents an intersection on the board.
type Symbol int

const (
	// Empty is the default symbol.
	Empty Symbol = 0

	// TopLeft corner of the board.
	TopLeft Symbol = 10
	// TopRight corner of the board.
	TopRight Symbol = 11
	// BottomLeft corner of the board.
	BottomLeft Symbol = 12
	// BottomRight corner of the board.
	BottomRight Symbol = 13
	// TopEdge of the board.
	TopEdge Symbol = 14
	// BottomEdge of the board.
	BottomEdge Symbol = 15
	// LeftEdge of the board.
	LeftEdge Symbol = 16
	// RightEdge of the board.
	RightEdge Symbol = 17
	// Center of the board.
	Center Symbol = 18
	// StarPoint of the board.
	StarPoint Symbol = 19

	// BlackStone is a black stone.
	BlackStone Symbol = 20
	// WhiteStone is a white stone.
	WhiteStone Symbol = 21

	// Triangle is a triangle-mark.
	Triangle Symbol = 30
	// Square is a square-mark.
	Square Symbol = 31
	// Circle is a circle-mark.
	Circle Symbol = 32
	// Xmark is an x-mark.
	Xmark Symbol = 33

	// TextLabel represents some arbitrary text-label mark.
	TextLabel Symbol = 34
)

// UnicodeString converts a symbol representation to a single char. This method is
// primarily for debugging.
func (s Symbol) UnicodeString() string {
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
	case StarPoint:
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
