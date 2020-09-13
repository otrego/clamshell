package sgf

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/otrego/clamshell/core/game"
)

// Parser parses SGFs into game trees
type Parser struct {
	rdr io.RuneReader
}

// FromString creates a parser from a string.
func FromString(sgf string) *Parser {
	return FromReader(strings.NewReader(sgf))
}

// FromReader creates a parser from a string reader.
func FromReader(r *strings.Reader) *Parser {
	return &Parser{
		rdr: r,
	}
}

// stateData contains the current parser state.
type stateData struct {
	idx, row, col int
	curchar       rune
	prevchar      rune
	curstate      parseState

	// tmp buffer for holding an escape char '\' during property data.
	holdChar rune
	buf      strings.Builder

	branches []*game.Node
	curnode  *game.Node
}

func (sd *stateData) addBranch(n *game.Node) {
	sd.branches = append(sd.branches, n)
}

func (sd *stateData) popBranch() (*game.Node, error) {
	if len(sd.branches) == 0 {
		return nil, sd.parseError("unable to pop-branch; likely due to empty variation")
	}
	parent := sd.branches[len(sd.branches)-1]
	sd.branches = sd.branches[:len(sd.branches)-1]
	return parent, nil
}

func (sd *stateData) addToBuf(c rune) {
	sd.buf.WriteRune(c)
}

func (sd *stateData) flushBuf() string {
	o := sd.buf.String()
	sd.buf = strings.Builder{}
	return o
}

// parseError creates a parsing error, which fives index, line, col and
// charector context
func (sd *stateData) parseError(msg string) error {
	return fmt.Errorf(
		"error during state %v, at index %v, line %v, column %v, char '%v'. %v",
		sd.curstate, sd.idx, sd.row, sd.col, string(sd.curchar), msg)
}

// propBuffer contains a buffer of property data that has yet to be flushed.
type propBuffer struct {
	prop     string
	propdata []string
}

func (b *propBuffer) flush(n *game.Node) {
	if b.prop != "" && len(b.propdata) != 0 {
		// Properties cannot be empty, propdata must be non-zero
		//
		// Validation should happen here.
		n.Properties[b.prop] = b.propdata
	}
	b.prop = ""
	b.propdata = []string{}
}

func (b *propBuffer) addToData(s string) {
	b.propdata = append(b.propdata, s)
}

// special chars, used to delimit sections of the SGF.
const (
	lparen    rune = '('
	rparen    rune = ')'
	lbrace    rune = '['
	rbrace    rune = ']'
	scolon    rune = ';'
	newline   rune = '\n'
	backslash rune = '\\'
)

// parseState is used to denate parseState
type parseState int

const (
	beginningState parseState = iota
	propertyState
	propDataState
	betweenState
)

// String converts the parseState to a readable string
func (p parseState) String() string {
	switch p {
	case beginningState:
		return "beginning"
	case propertyState:
		return "property"
	case propDataState:
		return "propertyData"
	case betweenState:
		return "between"
	default:
		return "unknown state"
	}
}

// Parse parses a game into a tree of moves, return a game or a parsing error.
func (p *Parser) Parse() (*game.Game, error) {
	g := game.New()
	sd := &stateData{}
	pbuf := &propBuffer{}

	var c rune
	var err error
	for c, _, err = p.rdr.ReadRune(); err == nil; c, _, err = p.rdr.ReadRune() {
		sd.idx++
		sd.col++
		sd.curchar = c
		if sd.curchar == newline {
			sd.row++
			sd.col = 0
		}

		switch sd.curstate {
		case beginningState:
			err := handleBeginning(sd, pbuf, g)
			if err != nil {
				return nil, err
			}

		case betweenState:
			err := handleBetween(sd, pbuf)
			if err != nil {
				return nil, err
			}

		case propertyState:
			err := handleProperty(sd, pbuf)
			if err != nil {
				return nil, err
			}

		case propDataState:
			err := handlePropData(sd, pbuf)
			if err != nil {
				return nil, err
			}

		default:
			// This is unlkely to happen unless we messed up our parser correctness.
			return nil, sd.parseError("unexpected parsing state")
		}
		sd.prevchar = c
	}

	if err == nil || !errors.Is(err, io.EOF) {
		return nil, sd.parseError(fmt.Sprintf("expected to end on EOF; got %v", err))
	} else if sd.prevchar != rparen {
		return nil, sd.parseError(fmt.Sprintf("expected to end on ')' got %c", sd.prevchar))
	} else if len(sd.branches) != 0 {
		return nil, sd.parseError("expected to end on root branch, but ended in nested condition")
	}

	return g, nil
}

// handleBeginning handles the beginning state, initializing the first (root)
// node.
//
// Transitions:
//
//     beginning => between
func handleBeginning(sd *stateData, pbuf *propBuffer, g *game.Game) error {
	if unicode.IsSpace(sd.curchar) {
		return nil // We can safely ignore whitespace here.
	} else if sd.curchar == lparen {
		// (;AW[aw][bw]
		// ^
		sd.branches = append(sd.branches, g.Root)
		return nil
	} else if sd.curchar == scolon {
		// (;AW[aw][bw]
		//  ^
		sd.curstate = betweenState
		sd.curnode = g.Root
		return nil
	}
	return sd.parseError("unexpected char")
}

// handleBetween handles the between state, transitioning to one of several
// other states based on the next characters.
//
// Transitions:
//
//     between => propertyState        ex: AW
//     between => propDataState        ex: AW[ab][
//     between => between, add branch  ex: B[ab](
//     between => between, pop branch  ex: B[ab](;W[ac])
//     between => between, add node    ex: B[ab];
func handleBetween(sd *stateData, pbuf *propBuffer) error {
	if unicode.IsSpace(sd.curchar) {
		// We can safely ignore whitespace here.
		return nil
	} else if unicode.IsUpper(sd.curchar) {
		// AW[aw][bw]
		// ^
		pbuf.flush(sd.curnode)
		sd.addToBuf(sd.curchar)
		sd.curstate = propertyState
		return nil
	} else if sd.curchar == lbrace {
		// AW[aw][bw]
		//   ^   ^
		sd.curstate = propDataState
		return nil
	} else if sd.curchar == lparen {
		// AW[aw][bw] (;B[ab]
		//            ^
		pbuf.flush(sd.curnode)
		sd.addBranch(sd.curnode)
		return nil
	} else if sd.curchar == scolon {
		// AW[aw][bw] (;B[ab];W[ac])
		//             ^     ^
		pbuf.flush(sd.curnode)
		cn := sd.curnode
		sd.curnode = game.NewNode()
		cn.AddChild(sd.curnode)
		return nil
	} else if sd.curchar == rparen {
		// AW[aw][bw] (;B[ab])
		//                   ^
		pbuf.flush(sd.curnode)
		cn, err := sd.popBranch()
		if err != nil {
			return err
		}
		sd.curnode = cn
		return nil
	}
	return sd.parseError("unexpected character between")
}

// handleProperty is a simple state for handling parsing property
// keys (AW, C, etc)
//
// Transitions:
//
//     property => propData
func handleProperty(sd *stateData, pbuf *propBuffer) error {
	if unicode.IsUpper(sd.curchar) {
		// AW[aw][bw]
		//  ^
		sd.addToBuf(sd.curchar)
		return nil
	} else if sd.curchar == lbrace {
		// AW[aw][bw]
		//   ^
		prop := sd.flushBuf()
		// TODO(kashomon): Add validation to ignore invalid properties
		pbuf.prop = prop
		sd.curstate = propDataState
		return nil
	}
	return sd.parseError("unexpected character during property parsing")
}

// handlePropData is a state for handling parsing property
// data
//
// Transitions:
//
//     propData => propData
//     propData => between
func handlePropData(sd *stateData, pbuf *propBuffer) error {
	if sd.curchar == rbrace &&
		// C[foo 1[k\] bar]
		//           ^
		// There was a backslash used for escaping a r-brace.
		sd.holdChar == backslash {
		sd.holdChar = rune(0)
		sd.addToBuf(sd.curchar)
		return nil
	} else if sd.holdChar == backslash {
		// C[foo 1[k\z bar]
		//           ^
		// Turns out the '\\' char wasn't meant to be used for escaping
		sd.addToBuf(sd.prevchar)
		sd.addToBuf(sd.curchar)
		sd.holdChar = rune(0)
		return nil
	} else if sd.curchar == backslash {
		// C[foo 1[k\] bar]
		//          ^
		// Ignore for now, wait for next character
		sd.holdChar = sd.curchar
		return nil
	} else if sd.curchar == rbrace {
		// C[foo 1[k\] bar]
		//                ^
		pbuf.addToData(sd.flushBuf())
		sd.curstate = betweenState
		return nil
	}
	// C[foo 1[k\] bar]
	//    ^
	sd.addToBuf(sd.curchar)
	return nil
}
