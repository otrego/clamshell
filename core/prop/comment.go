package prop

import (
	"fmt"
	"strings"

	"github.com/otrego/clamshell/core/movetree"
)

// commentConv is an SGF converter for the comment property C.
var commentConv = &SGFConverter{
	Props: []Prop{"C"},
	Scope: AllScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if n.Comment != "" {
			return fmt.Errorf("comment already found on node: %q", n.Comment)
		}
		if len(data) == 0 {
			// Edgecase where Comment-property is set, but there is no data.
			return nil
		} else if len(data) != 1 {
			return fmt.Errorf("comment only allows one prop-value, found %v", data)
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
