package movetree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/otrego/clamshell/go/board"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
)

var ErrParseTreepath = errors.New("error parsing treepath")

var ErrApplyTreepath = errors.New("error applying treepath")

// A Path is a list of variations that says how to travel through a tree of
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
type Path []int

type parseState int

const (
	// variationState means we are looking for numbers for the variation number. Can
	// either transition to either separator (flush) or repeat (flush and then
	// repeat N times
	variationState parseState = iota
	// repeat means we repeat the previous variation 'n' times.
	repeatState
)

// ParsePath parses a treepath from a string to an array of ints (Path).
//
// There are a couple different string short-hands that make using treepaths
// a little easier
//
//    0-1x2   Take the 0th variation, then repeat taking the 1st varation twice
//
// Paths say how to get from position n to position m.  Thus the numbers are
// always variations, separated by '.' except in the case of A:B syntax, which
// means repeat A for B times.
//
// A leading . is an optional. It can be useful to indicate the root node.
//
// Some examples:
//
//    -             becomes []
//    0             becomes [0]
//    -0            becomes [0]
//    1             becomes [1]
//    53            becomes [53] (the 53rd variation)
//    2-3           becomes [2,3]
//    0-0-0-0       becomes [0,0,0,0]
//    0x4           becomes [0,0,0,0]
//    1x4           becomes [1,1,1,1]
//    1-2x1-0-2x3   becomes [1,2,0,2,2,2]
func ParsePath(path string) (Path, error) {
	out := Path{}
	curState := variationState
	buf := strings.Builder{}
	prevChar := 0

	// As a special case, just return the empty path as a degenerate case.
	if path == "-" || path == "" {
		return out, nil
	}

	convertBuffer := func(idx int) (int, error) {
		if buf.Len() == 0 {
			return 0, fmt.Errorf("path %v at index %d: separators '-', 'x', EOL must be proceeded by digit: %w", path, idx, ErrParseTreepath)
		}
		n, err := strconv.Atoi(buf.String())
		if err != nil {
			return 0, fmt.Errorf("bad number in path %v at index %d: %v: %w", path, idx, err, ErrParseTreepath)
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
	// buffer data when we hit a '-' rune.
	//
	// REPEAT: means we have seen a 'x' instead of a '-', so we're going to repeat
	// the previous number once we parse the next number.
	//
	//  --'-'--- <------------
	//  |      |             |
	//  V      ^             ^
	// VARIATION --'x'-->  REPEAT
	//   |                   V
	//   EOL <----------------
	//   |
	//   V
	//   End
	for idx, c := range path {
		if unicode.IsDigit(c) {
			buf.WriteRune(c)
		} else if c == '-' || c == '\n' {
			if idx == 0 {
				// - is supported as an optional leading charactor. Just ignore.
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
				return nil, fmt.Errorf("path %q at index %d: unknown state %v: %w", path, idx, curState, ErrParseTreepath)
			}
		} else if c == 'x' {
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
				return nil, fmt.Errorf("path %q at index %d: repeat char ':' cannot be followed by another repeat ':': %w", path, idx, ErrParseTreepath)
			} else {
				return nil, fmt.Errorf("path %q at index %d: unknown state %v: %w", path, idx, curState, ErrParseTreepath)
			}
		} else {
			return nil, fmt.Errorf("path %q at index %d: unexpected character %q: %w", path, idx, c, ErrParseTreepath)
		}
	}
	return out, nil
}

// Apply applies a treepath to a node, moving down a tree of moves. This
// takes the variation listed in the treepath until either:
//
// 1. There are no more variation numbers in the treepath.
// 2. There is not a child with the given variation number.
func (tp Path) Apply(n *Node) *Node {
	curNode := n
	for _, v := range tp {
		if v < len(curNode.Children) {
			curNode = curNode.Children[v]
		} else {
			break
		}
	}
	return curNode
}

// ApplyToBoard applies a treepath to a Go-Board, returning the captured stones,
// or an error if the application was unsuccessful.
//
// A board copy, and the relevant captures are returned
func (tp Path) ApplyToBoard(n *Node, b *board.Board) (*board.Board, move.List, error) {
	b = b.Clone()

	applyStones := func(n *Node, bb *board.Board) (move.List, error) {
		err := bb.SetPlacements(n.Placements)
		if err != nil {
			return nil, err
		}
		if n.Move != nil && n.Move.Color() != color.Empty {
			return bb.PlaceStone(n.Move)
		}
		return nil, nil
	}

	var traversed Path
	var captures move.List
	for i := 0; n != nil; i++ {
		ml, err := applyStones(n, b)
		if err != nil {
			return nil, nil, fmt.Errorf("at traversed path %v: %v: %w", traversed, err, err)
		}
		captures = append(captures, ml...)
		if i >= len(tp) {
			n = nil
			continue
		}
		nextVar := tp[i]

		if nextVar < len(n.Children) {
			n = n.Children[nextVar]
		} else {
			n = nil
		}
	}
	captures.Sort()
	return b, captures, nil
}

// String returns the treepath as a string.
// examples:
//      []                  becomes "[]"
//      [1]                 becomes "[1]"
//      [1,2,0,2,2,2]       becomes "[1,2,0,2,2,2]"
func (tp Path) String() string {
	var strArr []string
	for _, i := range tp {
		strArr = append(strArr, strconv.Itoa(i))
	}
	return fmt.Sprintf("%v", strArr)
}

// Clone makes a copy of the treepath.
func (tp Path) Clone() Path {
	newPath := make(Path, len(tp))
	for i := range tp {
		newPath[i] = tp[i]
	}
	return newPath
}

// CompactString returns the treepath as a CompactString (short-hand).
// examples:
//      []                  becomes "-"
//      [1]                 becomes "-1"
//      [0,0,0,0]           becomes "-0x4"
//      [1,1,1,1]           becomes "-1x4"
//      [1,2,0,0,2,2,2]     becomes "-1-2-0x2-2x3"
func (tp Path) CompactString() string {
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
			sb.WriteString(fmt.Sprintf("-%dx%d", prev, count))
			count = 1
		} else if prev != -1 {
			//write non repeated variation number.
			sb.WriteString(fmt.Sprintf("-%d", prev))
		}
		prev = v
	}
	if prev == -1 {
		// empty path
		sb.WriteString("-")
	} else if count != 1 {
		// end of treepath was a repeated variation number.
		sb.WriteString(fmt.Sprintf("-%dx%d", prev, count))
	} else {
		// end or treepath was a new variation number.
		sb.WriteString(fmt.Sprintf("-%d", prev))
	}
	return sb.String()
}
