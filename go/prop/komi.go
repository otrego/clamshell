package prop

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/otrego/clamshell/go/movetree"
)

var ErrKomi = errors.New("error converting komi proprtey KM")

// komiConv converts the komi property KM.
var komiConv = &SGFConverter{
	Props: []Prop{"KM"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		komi, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			return err
		}
		_, fp := math.Modf(komi)
		if !(fp == 0.5 || fp == 0.0) {
			return fmt.Errorf("value was %f, but the only decimal-value allowed for komi is .0 or .5: %w", komi, ErrKomi)
		}
		if n.GameInfo == nil {
			// For safety, make sure to set create gameinfo if it doesn't exist.
			n.GameInfo = &movetree.GameInfo{}
		}
		n.GameInfo.Komi = new(float64)
		*n.GameInfo.Komi = komi
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if n.GameInfo == nil {
			return "", nil
		}
		if n.GameInfo.Komi == nil {
			return "", nil
		}
		komi := *n.GameInfo.Komi
		_, fp := math.Modf(komi)
		if !(fp == 0.5 || fp == 0.0) {
			return "", fmt.Errorf("komi value %f, but the only decimal-value allowed for komi is .0 or .5: %w", komi, ErrKomi)
		}
		s := strconv.FormatFloat(komi, 'f', 1, 64)
		return fmt.Sprintf("KM[%s]", s), nil
	},
}
