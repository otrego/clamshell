package prop

import (
	"fmt"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/movetree"
)

var initPlayerConv = &SGFConverter{
	Props: []Prop{"PL"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if len(data) != 1 {
			return fmt.Errorf("PL property requires exactly 1 Value, but had %d", len(data))
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
			return fmt.Errorf("Prop PL has invalid value of %s", data[0])
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
		return "", fmt.Errorf("prop PL can only have value W or B, but was %s", n.GameInfo.Player)
	},
}
