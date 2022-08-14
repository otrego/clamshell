package sgf

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/prop"
)

var ErrParse = errors.New("error parsing SGF")

// Parse is a convenience helper to parse sgf strings.
func Parse(s string) (*movetree.MoveTree, error) {
	return FromString(s).Parse()
}

// Parser parses SGFs into MoveTree objects.
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

	branches []*movetree.Node
	curnode  *movetree.Node
}

func (sd *stateData) addBranch(n *movetree.Node) {
	sd.branches = append(sd.branches, n)
}

func (sd *stateData) popBranch() (*movetree.Node, error) {
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
		"during state %v, at index %v, line %v, column %v, char '%v'. %v: %w",
		sd.curstate, sd.idx, sd.row, sd.col, string(sd.curchar), msg, ErrParse)
}

// propBuffer contains a buffer of property data that has yet to be flushed.
type propBuffer struct {
	prop     string
	propdata []string
}

func (b *propBuffer) flush(n *movetree.Node) error {
	if b.prop != "" && len(b.propdata) != 0 {
		if err := prop.ProcessPropertyData(n, b.prop, b.propdata); err != nil {
			return err
		}
	}
	b.prop = ""
	b.propdata = []string{}
	return nil
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

// Parse parses a movetree into a tree of moves, return a movetree or a parsing
// error.
func (p *Parser) Parse() (*movetree.MoveTree, error) {
	g := movetree.New()
	stateData := &stateData{}
	pbuf := &propBuffer{}

	// the parser uses a finite state machine to perform parsing, having the
	// following states & actions
	//
	//               addBranch '(', popBranch ')'
	//               V ^
	//               | |
	//               | |  -------------------']'--
	//               V ^  V                      |
	// BEGINNING => BETWEEN => PROPERTY => PROPERTY DATA
	//               |    V                      ^
	//               |    --------------'['-------
	//               V
	//              END
	var c rune
	var err error
	for c, _, err = p.rdr.ReadRune(); err == nil; c, _, err = p.rdr.ReadRune() {
		stateData.idx++
		stateData.col++
		stateData.curchar = c
		if stateData.curchar == newline {
			stateData.row++
			stateData.col = 0
		}

		switch stateData.curstate {
		case beginningState:
			if err = handleBeginning(stateData, pbuf, g); err != nil {
				return nil, err
			}

		case betweenState:
			if err = handleBetween(stateData, pbuf); err != nil {
				return nil, err
			}

		case propertyState:
			if err = handleProperty(stateData, pbuf); err != nil {
				return nil, err
			}

		case propDataState:
			if err = handlePropData(stateData, pbuf); err != nil {
				return nil, err
			}

		default:
			// This is unlkely to happen unless we messed up our parser correctness.
			return nil, stateData.parseError("unexpected parsing state")
		}
		stateData.prevchar = c
	}

	// We should **always** end with an EOF error
	if err == nil || !errors.Is(err, io.EOF) {
		return nil, stateData.parseError(fmt.Sprintf("expected to end on EOF; got %v", err))
	} else if len(stateData.branches) != 0 {
		return nil, stateData.parseError("expected to end on root branch, but ended in nested condition")
	}

	return g, nil
}

// handleBeginning handles the beginning state, initializing the first (root)
// node.
//
// Transitions:
//
//     beginning => between
func handleBeginning(stateData *stateData, pbuf *propBuffer, g *movetree.MoveTree) error {
	if unicode.IsSpace(stateData.curchar) {
		return nil // We can safely ignore whitespace here.
	} else if stateData.curchar == lparen {
		// (;AW[aw][bw]
		// ^
		stateData.branches = append(stateData.branches, g.Root)
		return nil
	} else if stateData.curchar == scolon {
		// (;AW[aw][bw]
		//  ^
		stateData.curstate = betweenState
		stateData.curnode = g.Root
		return nil
	}
	return stateData.parseError("unexpected char")
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
func handleBetween(stateData *stateData, pbuf *propBuffer) error {
	if unicode.IsSpace(stateData.curchar) {
		// We can safely ignore whitespace here.
		return nil
	} else if unicode.IsUpper(stateData.curchar) {
		// AW[aw][bw]
		// ^
		if err := pbuf.flush(stateData.curnode); err != nil {
			return stateData.parseError(err.Error())
		}
		stateData.addToBuf(stateData.curchar)
		stateData.curstate = propertyState
		return nil
	} else if stateData.curchar == lbrace {
		// AW[aw][bw]
		//   ^   ^
		stateData.curstate = propDataState
		return nil
	} else if stateData.curchar == lparen {
		// AW[aw][bw] (;B[ab]
		//            ^
		pbuf.flush(stateData.curnode)
		if err := pbuf.flush(stateData.curnode); err != nil {
			return stateData.parseError(err.Error())
		}
		stateData.addBranch(stateData.curnode)
		return nil
	} else if stateData.curchar == scolon {
		// AW[aw][bw] (;B[ab];W[ac])
		//             ^     ^
		pbuf.flush(stateData.curnode)
		if err := pbuf.flush(stateData.curnode); err != nil {
			return stateData.parseError(err.Error())
		}
		cn := stateData.curnode
		stateData.curnode = movetree.NewNode()
		cn.AddChild(stateData.curnode)
		stateData.curnode.Parent = cn
		return nil
	} else if stateData.curchar == rparen {
		// AW[aw][bw] (;B[ab])
		//                   ^
		pbuf.flush(stateData.curnode)
		if err := pbuf.flush(stateData.curnode); err != nil {
			return stateData.parseError(err.Error())
		}
		cn, err := stateData.popBranch()
		if err != nil {
			return err
		}
		stateData.curnode = cn
		return nil
	}
	return stateData.parseError("unexpected character between")
}

// handleProperty is a simple state for handling parsing property
// keys (AW, C, etc)
//
// Transitions:
//
//     property => propData
func handleProperty(stateData *stateData, pbuf *propBuffer) error {
	if unicode.IsUpper(stateData.curchar) {
		// AW[aw][bw]
		//  ^
		stateData.addToBuf(stateData.curchar)
		return nil
	} else if stateData.curchar == lbrace {
		// AW[aw][bw]
		//   ^
		prop := stateData.flushBuf()
		pbuf.prop = prop
		stateData.curstate = propDataState
		return nil
	}
	return stateData.parseError("unexpected character during property parsing")
}

// handlePropData is a state for handling parsing property
// data
//
// Transitions:
//
//     propData => propData
//     propData => between
func handlePropData(stateData *stateData, pbuf *propBuffer) error {
	if stateData.curchar == rbrace && stateData.holdChar == backslash {
		// C[foo 1[k\] bar]
		//           ^
		// There was a backslash used for escaping a r-brace.
		stateData.holdChar = rune(0)
		stateData.addToBuf(stateData.curchar)
		return nil
	} else if stateData.holdChar == backslash {
		// C[foo 1[k\z bar]
		//           ^
		// Turns out the '\\' char wasn't meant to be used for escaping
		stateData.addToBuf(stateData.prevchar)
		stateData.addToBuf(stateData.curchar)
		stateData.holdChar = rune(0)
		return nil
	} else if stateData.curchar == backslash {
		// C[foo 1[k\] bar]
		//          ^
		// Ignore for now, wait for next character
		stateData.holdChar = stateData.curchar
		return nil
	} else if stateData.curchar == rbrace {
		// C[foo 1[k\] bar]
		//                ^
		pbuf.addToData(stateData.flushBuf())
		stateData.curstate = betweenState
		return nil
	}
	// C[foo 1[k\] bar]
	//    ^
	stateData.addToBuf(stateData.curchar)
	return nil
}
