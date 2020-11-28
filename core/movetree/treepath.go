package movetree

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// A Treepath is a list of variations that says how to travel through a tree of
// moves. And has two forms a list version and a string version. First, the
// list-version:
//
//    [0,1,0]
//
// Means we will first take the 0th variation, then we will take the 1ist
// variation, and lastly we will take the 0th variation again. For
// convenience, the treepath can also be specified by a string, which is where
// the fun begins. At its simpliest,
//
//    [0,0,0] becomes 0.0.0
type Treepath []int

type parseState int

const (
	// variationState means we are looking for numbers for the variation number. Can
	// either transition to either separator (flush) or repeat (flush and then
	// repeat N times
	variationState parseState = iota
	// repeat means we repeat the previous variation 'n' times.
	repeatState
)

// ParsePath parses a treepath from a string to an array of ints (Treepath).
//
// There are a couple different string short-hands that make using treepaths
// a little easier
//
//    0.1:2   Take the 0th variation, then repeat taking the 1st varation twice
//
// Treepaths say how to get from position n to position m.  Thus the numbers are
// always variations, separated by '.' except in the case of A:B syntax, which
// means repeat A for B times.
//
// A leading . is an optional. It can be useful to indicate the root node.
//
// Some examples:
//
//    .             becomes []
//    0             becomes [0]
//    .0            becomes [0]
//    1             becomes [1]
//    53            becomes [53] (the 53rd variation)
//    2.3           becomes [2,3]
//    0.0.0.0       becomes [0,0,0,0]
//    0:4           becomes [0,0,0,0]
//    1:4           becomes [1,1,1,1]
//    1.2:1.0.2:3   becomes [1,2,0,2,2,2]
func ParsePath(path string) (Treepath, error) {
	out := Treepath{}
	curState := variationState
	buf := strings.Builder{}
	prevChar := 0

	// As a special case, just return the empty path as a degenerate case.
	if path == "." || path == "" {
		return out, nil
	}

	convertBuffer := func(idx int) (int, error) {
		if buf.Len() == 0 {
			return 0, fmt.Errorf("error parsing path %v at index %d: separators '.', ':', EOL must be proceeded by digit", path, idx)
		}
		n, err := strconv.Atoi(buf.String())
		if err != nil {
			return 0, fmt.Errorf("error parsing number in path %v at index %d: %v", path, idx, err)
		}
		buf = strings.Builder{}
		return n, nil
	}

	// Use a sentinal newline character so that it's clear when to flush the final
	// charactrs.
	path += "\n"

	// The rough approach is to have a simple parser modeled as a 2-state DFA.
	//
	// VARIATION: means we're parsing a normal variation number. We flush the
	// buffer data when we hit a '.' rune.
	//
	// REPEAT: means we have seen a ':' instead of a '.', so we're going to repeat
	// the previous number once we parse the next number.
	//
	//  --'.'--- <------------
	//  |      |             |
	//  V      ^             ^
	// VARIATION --':'-->  REPEAT
	//   |                   V
	//   EOL <----------------
	//   |
	//   V
	//   End
	for idx, c := range path {
		if unicode.IsDigit(c) {
			buf.WriteRune(c)
		} else if c == '.' || c == '\n' {
			if idx == 0 {
				// . is supported as an optional leading charactor. Just ignore.
				continue
			}
			n, err := convertBuffer(idx)
			if err != nil {
				return nil, err
			}

			// State transitions
			if curState == variationState {
				// Variation => Variation (& flush)
				out = append(out, n)
			} else if curState == repeatState {
				// Repeat => Variation (& flush)
				for i := 0; i < n; i++ {
					out = append(out, prevChar)
				}
				curState = variationState
				prevChar = 0
			} else {
				return nil, fmt.Errorf("error parsing path %q at index %d: unknown state %v", path, idx, curState)
			}
		} else if c == ':' {
			n, err := convertBuffer(idx)
			if err != nil {
				return nil, err
			}

			// State transitions
			if curState == variationState {
				// Variation => Repeat
				prevChar = n
				curState = repeatState
			} else if curState == repeatState {
				// Repeate => Repeat: Invalid
				return nil, fmt.Errorf("error parsing path %q at index %d: repeat char ':' cannot be followed by another repeat ':'", path, idx)
			} else {
				return nil, fmt.Errorf("error parsing path %q at index %d: unknown state %v", path, idx, curState)
			}
		} else {
			return nil, fmt.Errorf("error parsing path %q at index %d: unexpected character %q", path, idx, c)
		}
	}
	return out, nil
}

// Apply applies a treepath to a node, moving down a tree of moves. This
// takes the variation listed in the treepath until either:
//
// 1. There are no more variation numbers in the treepath.
// 2. There is not a child with the given variation number.
func (tp Treepath) Apply(n *Node) *Node {
	curNode := n
	for _, v := range tp {
		if v < len(curNode.Children) {
			// Assume there are no gaps. If there are, parsing failed us.
			curNode = curNode.Children[v]
		} else {
			break
		}
	}
	return curNode
}

// String returns the treepath as a string.
// examples:
//      []                  becomes "[]"
//      [1]                 becomes "[1]"
//      [1,2,0,2,2,2]       becomes "[1,2,0,2,2,2]"
func (tp Treepath) String() string {
	var strArr []string
	for _, i := range tp {
		strArr = append(strArr, strconv.Itoa(i))
	}
	return fmt.Sprintf("%v", strArr)
}

// Clone makes a copy of the treepath.
func (tp Treepath) Clone() Treepath {
	// A shallow copy should be sufficient.
	return tp[:]
}

// CompactString returns the treepath as a CompactString (short-hand).
// examples:
//      []                  becomes "."
//      [1]                 becomes ".1"
//      [0,0,0,0]           becomes ".0:4"
//      [1,1,1,1]           becomes ".1:4"
//      [1,2,0,0,2,2,2]     becomes ".1.2.0:2.2:3"
func (tp Treepath) CompactString() string {
	var (
		count, prev int = 1, -1
		sb          strings.Builder
	)

	for _, v := range tp {

		if v == prev {
			//count repeated variation numbers.
			count++
		} else if count != 1 {
			//write the repeated variation number.
			sb.WriteString(fmt.Sprintf(".%d:%d", prev, count))
			count = 1
		} else if prev != -1 {
			//write non repeated variation number.
			sb.WriteString(fmt.Sprintf(".%d", prev))
		}
		prev = v
	}
	if prev == -1 {
		//empty treepath
		sb.WriteString(".")
	} else if count != 1 {
		//end of treepath was a repeated variation number.
		sb.WriteString(fmt.Sprintf(".%d:%d", prev, count))
	} else {
		//end or treepath was a new variation number.
		sb.WriteString(fmt.Sprintf(".%d", prev))
	}
	return sb.String()
}
