// Package treepath provides functionality for manipulating treepaths.
//
// Treepaths are string-encodings of paths through trees.
//
// Treepaths come from
// https://github.com/Kashomon/glift-core/blob/master/src/rules/treepath.js
package treepath

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/otrego/clamshell/core/game"
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
//
// but there are a couple different string short-hands that make using treepaths
// a little easier
//
//    0.1:2   Take the 0th variation, then repeat taking the 1st varation twice
//
// There are two types of treepaths discussed below -- A *treepath fragment*
// (which is what we have been describing) and an *initial treepath*.
//
// ## Treepath
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

// Parse parses a treepath from a string to an array of ints.
func Parse(path string) (Treepath, error) {
	out := Treepath{}
	curState := variationState
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

	// Use a sentinal newline character so that it's clear when to flush the final
	// charactrs.
	path += "\n"

	//
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
	//
	for idx, c := range path {
		if unicode.IsDigit(c) {
			buf.WriteRune(c)
		} else if c == '.' || c == '\n' {
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
				return nil, fmt.Errorf("error parsing fragment %q at index %d: unknown state %v", path, idx, curState)
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
				return nil, fmt.Errorf("error parsing fragment %q at index %d: repeat char ':' cannot be followed by another repeat ':'", path, idx)
			} else {
				return nil, fmt.Errorf("error parsing fragment %q at index %d: unknown state %v", path, idx, curState)
			}
		} else {
			return nil, fmt.Errorf("error parsing fragment %q at index %d: unexpected character %q", path, idx, c)
		}
	}
	return out, nil
}

// Apply applies a treepath to a node, moving down a tree of moves. This
// takes the variation listed in the treepath until either:
//
// 1. There are no more variation numbers in the treepath.
// 2. There is not a child with the given variation number.
func (tp Treepath) Apply(n *game.Node) *game.Node {
	curNode := n
	for _, v := range tp {
		if v < len(curNode.Children) {
			fmt.Printf("var %v\n", v)
			// Assume there are no gaps. If there are, parsing failed us.
			curNode = curNode.Children[v]
		} else {
			break
		}
	}
	return curNode
}
