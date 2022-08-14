package prop

import (
	"errors"
	"fmt"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/movetree"
)

var ErrMove = errors.New("error converting move property B or W")

// movesConv is an SGF converter for moves B,W.
var movesConv = &SGFConverter{
	Props: []Prop{"B", "W"},
	Scope: AllScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		col, err := color.FromSGFProp(prop)
		if err != nil {
			return err
		}
		if n.Move != nil {
			return fmt.Errorf("found two moves on one node at move: %w", ErrMove)
		}
		if len(data) != 1 && len(data) != 0 {
			return fmt.Errorf("expected black move data to have exactly one value or zero values: %w", ErrMove)
		}
		if len(data) == 0 {
			data = []string{""}
		}
		move, err := move.FromSGFPoint(col, data[0])
		if err != nil {
			return err
		}
		n.Move = move
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		mv := n.Move
		if mv == nil {
			return "", nil
		}
		var col string
		if mv.Color() == color.Black {
			col = "B"
		} else if mv.Color() == color.White {
			col = "W"
		}
		if mv.IsPass() {
			// Return non-nil slice to indicate it should be stored.
			return col + "[]", nil
		}
		sgfPt, err := mv.Point().ToSGF()
		if err != nil {
			return "", err
		}
		return col + "[" + sgfPt + "]", nil
	},
}
