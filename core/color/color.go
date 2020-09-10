// Package color contains utilities related to player and stone color.
package color

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
