package prop

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/otrego/clamshell/go/movetree"
)

var ErrSize = errors.New("error converting size property SZ")

// sizeConv converts the size property SZ.
var sizeConv = &SGFConverter{
	Props: []Prop{"SZ"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if l := len(data); l != 1 {
			return fmt.Errorf("data must be exactly 1, was %d: %w", l, ErrSize)
		}
		sz, err := strconv.Atoi(data[0])
		if err != nil {
			return fmt.Errorf("parsing data %v as integer %v: %w", data, err, ErrSize)
		}
		if sz < 1 || sz > 25 {
			return fmt.Errorf("size was %d, but must be between 1 and 25 %w", sz, ErrSize)
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
			return "", fmt.Errorf("size was %d but only values between 1 and 25 are allowed: %w", sz, ErrSize)
		}
		return "SZ[" + strconv.Itoa(sz) + "]", nil
	},
}
