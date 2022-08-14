package prop

import (
	"errors"
	"fmt"
	"strings"

	"github.com/otrego/clamshell/go/movetree"
)

var ErrComment = errors.New("error converting comment property C")

// commentConv is an SGF converter for the comment property C.
var commentConv = &SGFConverter{
	Props: []Prop{"C"},
	Scope: AllScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if n.Comment != "" {
			return fmt.Errorf("%w: already found on node: %q", ErrComment, n.Comment)
		}
		if len(data) == 0 {
			// Edgecase where Comment-property is set, but there is no data.
			return nil
		} else if len(data) != 1 {
			return fmt.Errorf("%w: comment only allows one prop-value, found %v", ErrComment, data)
		}
		// The parser does escaping for the property data on the input side.
		n.Comment = data[0]
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		c := n.Comment
		if c == "" {
			return "", nil
		}
		return "C[" + strings.Replace(c, "]", "\\]", -1) + "]", nil
	},
}
