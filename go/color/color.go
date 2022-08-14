// Package color contains utilities related to player and stone color.
package color

import (
	"errors"
	"fmt"
)

// Color is typed string indicating player or stone color
type Color string

const (
	// Black player or stone color
	Black Color = "B"
	// White player or stone color
	White Color = "W"
	// Empty player or stone color
	Empty Color = ""
)

// Opposite returns the opposite color. If the color is unknown or empty, just
// return the same color.
func (c Color) Opposite() Color {
	if c == Black {
		return White
	} else if c == White {
		return Black
	}
	return c
}

// Ordinal returns an ordinal number for the color, for the purposes of sorting.
func (c Color) Ordinal() int {
	switch c {
	case Black:
		return 1
	case White:
		return 2
	case Empty:
		return 3
	default:
		// This shouldn't happen normally, but here for good measure.
		return 4
	}
}

// ErrColorConversion is an err
var ErrColorConversion = errors.New("color conversion error")

// FromSGFProp returns the color from a SGF property that's color related)
func FromSGFProp(prop string) (Color, error) {
	switch prop {
	case "B", "AB":
		return Black, nil
	case "W", "AW":
		return White, nil
	default:
		return Empty, fmt.Errorf("%w: converting property %q", ErrColorConversion, prop)
	}
}
