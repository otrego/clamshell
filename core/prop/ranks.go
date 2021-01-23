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

		if (num >= 1 && num <= 30 && string(rank[i[0]]) == "k") || (num >= 1 && num <= 9) {
			if n.GameInfo == nil {
				// For safety, make sure to set create gameinfo if it doesn't exist.
				n.GameInfo = &movetree.GameInfo{}
			}
			if prop == "BR" {
				n.GameInfo.BlackRank = rank
			} else {
				n.GameInfo.WhiteRank = rank
			}

		} else {
			return fmt.Errorf("Invalid number %d for rank %s", num, rank[i[0]:i[1]])
		}
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if n.GameInfo == nil {
			return "", nil
		}

		r, _ := regexp.Compile("(k|kyu|d|dan|p|pro)\\b")
		var out strings.Builder
		if n.GameInfo.BlackRank != "" {
			rank := n.GameInfo.BlackRank
			i := r.FindStringIndex(rank)
			if len(i) == 0 {
				return "", fmt.Errorf("Invalid Rank: %s", rank)
			}

			s := rank[0:i[0]]
			num, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if (num >= 1 && num <= 30 && string(rank[i[0]]) == "k") || (num >= 1 && num <= 9) {
				out.WriteString("BR[" + n.GameInfo.BlackRank + "]")
			} else {
				return "", fmt.Errorf("Invalid number %d for rank %s", num, rank[i[0]:i[1]])
			}
		}
		if n.GameInfo.WhiteRank != "" {
			rank := n.GameInfo.WhiteRank
			i := r.FindStringIndex(rank)
			if len(i) == 0 {
				return "", fmt.Errorf("Invalid Rank: %s", rank)
			}

			s := rank[0:i[0]]
			num, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if (num >= 1 && num <= 30 && string(rank[i[0]]) == "k") || (num >= 1 && num <= 9) {
				out.WriteString("WR[" + n.GameInfo.WhiteRank + "]")
			} else {
				return "", fmt.Errorf("Invalid number %d for rank %s", num, rank[i[0]:i[1]])
			}
		}

		if out.String() == "" {
			return "", nil
		}
		return out.String(), nil
	},
}
