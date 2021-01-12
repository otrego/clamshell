package snapshot

import (
	"github.com/otrego/clamshell/core/point"
	"github.com/otrego/clamshell/core/snapshot/symbol"
)

// Intersection contains information for a single intersection
type Intersection struct {
	// The Point for this intersection
	point *point.Point

	// The Base symbol: The first layer of the intesection, which tells.
	base symbol.Symbol

	stone symbol.Symbol

	mark symbol.Symbol

	// Label for the intersection. Label should only be set when Mark == TextLabel
	// or another
	label string
}

// Point returns the point.
func (n *Intersection) Point() *point.Point { return n.point }

// Base returns base layer.
func (n *Intersection) Base() symbol.Symbol { return n.base }

// Stone returns stone layer.
func (n *Intersection) Stone() symbol.Symbol { return n.stone }

// Mark returns the mark layer.
func (n *Intersection) Mark() symbol.Symbol { return n.mark }

// Label returns the label.
func (n *Intersection) Label() string { return n.label }

// TopLayerUnicodeChar outputs a single character for the intersection, based
// on the "Top" symbol. This method is primarily for debugging.
func (n *Intersection) TopLayerUnicodeChar() string {
	switch {
	case n.mark != symbol.Empty:
		return n.mark.UnicodeChar()
	case n.stone != symbol.Empty:
		return n.stone.UnicodeChar()
	case n.base != symbol.Empty:
		return n.base.UnicodeChar()
	}
	return " "
}
