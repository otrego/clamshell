// Package treepath provides functionality for manipulating treepaths.
//
// Treepaths come from
// https://github.com/Kashomon/glift-core/blob/master/src/rules/treepath.js
//
// A treepath is a list of variations that says how to travel through a tree of
// moves. And has two forms a list version and astring version. First, the
// list-version:
//
//    [0,1,0]
//
// Means we will first take the 0th variation, then we will take the 1ist
// variation, and lastly we will take the 0th variation again. For
// convenience, the treepath can also be specified by a string, which is where
// the fun begins. At it's simpliest,
//
//    [0,0,0] becomes 0.0.0
//
// but there are a couple different string short-hands that make using treepaths
// a little easier
//
//    0.1+    Take the 0th variation, then the 1st variation, then go to the end
//    0.1:2   Take the 0th variation, then repeat taking the 1st varation twice
//
// There are two types of treepaths discussed below -- A *treepath fragment*
// (which is what we have been describing) and an *initial treepath*.
//
// ## Treepath Fragments
//
// Treepaths say how to get from position n to position m.  Thus the numbers are
// always variations except in the case of AxB syntax, where B is a multiplier
// for a variation.
//
// This is how fragment strings are parsed:
//
//    0             becomes [0]
//    1             becomes [1]
//    53            becomes [53] (the 53rd variation)
//    2.3           becomes [2,3]
//    0.0.0.0       becomes [0,0,0]
//    0:4           becomes [0,0,0,0]
//    1:4           becomes [1,1,1,1]
//    1.2:1.0.2:3   becomes [1,2,0,2,2,2]
//
// ## Initial tree paths.
//
// The initial treepath always treats the first number as a 'move number'
// instead of a variation. Thus
//
//    3.1.0
//
// means start at move 3 (always taking the 0th variation path) and then take
// the path fragment [1,0].
//
// Some examples:
//
//    0         - Start at the 0th move (the root node)
//    53        - Start at the 53rd move (taking the 0th variation)
//    2.3       - Start at the 3rd variation on move 2 (actually move 3)
//    3         - Start at the 3rd move
//    2.0       - Start at the 3rd move
//    0.0.0.0   - Start at the 3rd move
//    0.0:3     - Start at the 3rd move
//
// As with fragments, the init position returned is an array of variation
// numbers traversed through.  The move number is precisely the length of the
// array.
//
// So, for parsing
//
//    0         becomes []
//    1         becomes [0]
//    0.1       becomes [1]
//    53        becomes [0,0,0,...,0] (53 times)
//    2.3       becomes [0,0,3]
//    0.0.0.0   becomes [0,0,0]
//    0.0:3.1:3 becomes [0,0,0,1,1,1]
//
// As mentioned before, '+' is a special symbol which means "go to the end via
// the first variation." This is implemented with a by appending 500 0s to the
// path array.  This is a hack, but in practice games don't go over 500 moves.
package treepath

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// A treepath represents a variation path through a Game-tree
type Treepath []int

// ParseInitialPath Parses an initial treepath.
//
// The difference between an InitialTreePath and a Treepath-fragment
// is that the first move is interprented as a move number.
//
// Ex:
//
// 53 = [0, 0, 0, ... (53 times)]
func ParseInitialPath(initPos string) (Treepath, error) {
	ns := make([]rune, 0, 3) // very unlikely that n > 1000
	idx := 0
	for _, r := range initPos {
		if unicode.IsDigit(r) {
			idx++
			ns = append(ns, r)
		} else if r == '.' {
			idx++
			break
		} else {
			return nil, fmt.Errorf("failed to parse initial path %q, unexpected char at index %d. expected digit or '.'", initPos, idx)
		}
	}
	n, err := strconv.Atoi(string(ns))
	if err != nil {
		return nil, err
	}

	tp := make(Treepath, n, n)
	for i := 0; i < n; i++ {
		tp = append(tp, 0)
	}

	if idx >= len(initPos) {
		// We have captured the full string.
		return tp, nil
	}

	tail, err := ParseFragment(initPos[idx:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse initial path %q: %v", initPos, err)
	}
	tp = append(tp, tail...)

	return tp, nil
}

type parseState int

const (
	// variation means we are looking for numbers for the variation number. Can
	// either transition to either separator (flush) or repeat (flush and then
	// repeat N times
	variation parseState = iota
	// repeat means we repeat the previous variation 'n' times.
	repeat
)

// ParseFragment parses a treepath fragment.
func ParseFragment(path string) (Treepath, error) {
	out := make(Treepath, 0, len(path)/2)
	curState := variation
	buf := &strings.Builder{}
	prevChar := 0

	convertBuffer := func(idx int) (int, error) {
		if buf.Len() == 0 {
			return 0, fmt.Errorf("error parsing fragment %q at index %d: separator '.' must be proceeded by digit", path, idx)
		}
		n, err := strconv.Atoi(buf.String())
		if err != nil {
			return 0, fmt.Errorf("error parsing number in fragment %q at index %d: %v", path, idx, err)
		}
		buf = &strings.Builder{}
		return n, nil
	}

	for idx, c := range path {
		if unicode.IsDigit(c) {
			buf.WriteRune(c)
		} else if c == '.' {
			n, err := convertBuffer(idx)
			if err != nil {
				return nil, err
			}
			if curState == variation {
				out = append(out, n)
			} else if curState == repeat {
				for i := 0; i < n; i++ {
					out = append(out, prevChar)
				}
			}
		} else if c == ':' {
			n, err := convertBuffer(idx)
			if err != nil {
				return nil, err
			}
			if curState == variation {
				prevChar = n
			} else if curState == repeat {
				fmt.Errorf("error parsing fragment %q at index %d: repeat char ':' cannot be followed by another repeat ':'", path, idx)
			} else {
				fmt.Errorf("error parsing fragment %q at index %d: unexpectad character %q", path, idx, c)
			}
		}
	}
	return out, nil
}
