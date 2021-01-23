package prop

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/otrego/clamshell/core/movetree"
)

// ranksConv converts the rank properties BR and WR.
var ranksConv = &SGFConverter{
	Props: []Prop{"BR", "WR"},
	Scope: RootScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		if len(data) != 1 {
			return fmt.Errorf("%s propertie requires exactly 1 value, but had %d", prop, len(data))
		}

		rank := data[0]
		err := isValid(rank)
		if err != nil {
			return err
		}

		if n.GameInfo == nil {
			// For safety, make sure to set create gameinfo if it doesn't exist.
			n.GameInfo = &movetree.GameInfo{}
		}
		if prop == "BR" {
			n.GameInfo.BlackRank = rank
		} else {
			n.GameInfo.WhiteRank = rank
		}
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if n.GameInfo == nil {
			return "", nil
		}

		var out strings.Builder
		if n.GameInfo.BlackRank != "" {
			rank := n.GameInfo.BlackRank
			err := isValid(rank)
			if err != nil {
				return "", err
			}
			out.WriteString("BR[" + rank + "]")
		}

		if n.GameInfo.WhiteRank != "" {
			rank := n.GameInfo.WhiteRank
			err := isValid(rank)
			if err != nil {
				return "", err
			}
			out.WriteString("WR[" + rank + "]")
		}

		if out.String() == "" {
			return "", nil
		}
		return out.String(), nil
	},
}

func isValid(rank string) error {
	r, _ := regexp.Compile("(k|kyu|d|dan|p|pro)\\b")
	i := r.FindStringIndex(rank)
	if len(i) == 0 {
		return fmt.Errorf("Invalid Rank: %s", rank)
	}

	s := rank[0:i[0]]
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	if num > 30 || num < 1 || (num > 9 && string(rank[i[0]]) != "k") {
		return fmt.Errorf("Invalid number %d for rank %s", num, rank[i[0]:i[1]])
	}
	return nil
}
