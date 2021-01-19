package prop

import (
	"fmt"
	"strconv"

	"github.com/otrego/clamshell/core/movetree"
)

// sizeConv converts the size property SZ.
var sizeConv = &SGFConverter{
	Props: []Prop{"SZ"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if l := len(data); l != 1 {
			return fmt.Errorf("for prop %s, data  must be exactly 1, was %d", prop, l)
		}
		sz, err := strconv.Atoi(data[0])
		if err != nil {
			return fmt.Errorf("for prop %s, error parsing data %v as integer: %v", prop, data, err)
		}
		if sz < 1 || sz > 25 {
			return fmt.Errorf("for prop %s, size was %d, but must be between 1 and 25", prop, sz)
		}
		if n.GameInfo == nil {
			// For safety, make sure to set create gameinfo if it doesn't exist.
			n.GameInfo = &movetree.GameInfo{}
		}
		n.GameInfo.Size = sz
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if n.GameInfo == nil {
			return "", nil
		}
		sz := n.GameInfo.Size
		if sz == 0 {
			// BoardSize is unspecified.
			return "", nil
		}
		if sz < 1 || sz > 25 {
			return "", fmt.Errorf("invalid board size: %d, but only values between 1 and 25 are allowed", sz)
		}
		return "SZ[" + strconv.Itoa(sz) + "]", nil
	},
}
