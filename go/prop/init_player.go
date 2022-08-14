package prop

import (
	"errors"
	"fmt"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/movetree"
)

var ErrInitPlayer = errors.New("error converting initial player property PL")

var initPlayerConv = &SGFConverter{
	Props: []Prop{"PL"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if len(data) != 1 {
			return fmt.Errorf("requires exactly 1 Value, but had %d: %w", len(data), ErrInitPlayer)
		}
		if n.GameInfo == nil {
			// For safety, make sure to set create gameinfo if it doesn't exist.
			n.GameInfo = &movetree.GameInfo{}
		}

		switch data[0] {
		case "B", "b":
			n.GameInfo.Player = color.Black
		case "W", "w":
			n.GameInfo.Player = color.White
		default:
			return fmt.Errorf("invalid value of %s: %w", data[0], ErrInitPlayer)
		}

		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if n.GameInfo == nil {
			return "", nil
		}

		switch n.GameInfo.Player {
		case "":
			return "", nil
		case "W", "B":
			return fmt.Sprintf("PL[%s]", n.GameInfo.Player), nil
		}
		return "", fmt.Errorf("can only have value W or B, but was %s: %w", n.GameInfo.Player, ErrInitPlayer)
	},
}
