package snapshot

import (
	"github.com/otrego/clamshell/core/point"
	"github.com/otrego/clamshell/core/snapshot/symbol"
)

// Intersection contains information for a single intersection
type Intersection struct {
	// The Point for this intersection.
	Point *point.Point

	// The Base symbol is the first layer of the intesection, indicating the
	// underlying board symbol.
	Base symbol.Symbol

	// The Symbol symbol is the second layer of the intesection, indicating the
	// black or white stones.
	Stone symbol.Symbol

	// The Mark symbol is the third layer of the intesection, indicating the
	// marks on top of the board or stones.
	Mark symbol.Symbol

	// Label for the intersection. Label should only be set when Mark == TextLabel
	// or a similar label-mark.
	Label string
}

// TopLayerUnicodeString outputs a single character for the intersection, based
// on the "Top" symbol. This method is primarily for debugging.
func (n *Intersection) TopLayerUnicodeString() string {
	switch {
	case n.Mark != symbol.Empty:
		return n.Mark.UnicodeString()
	case n.Stone != symbol.Empty:
		return n.Stone.UnicodeString()
	case n.Base != symbol.Empty:
		return n.Base.UnicodeString()
	}
	return " "
}
