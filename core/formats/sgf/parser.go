package sgf

import (
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
	curstate      parseState

	lastCharInBuf rune
	buf           strings.Builder

	branches []*game.Node
	curnode  *game.Node
}

func (sd *stateData) addBranch(n *game.Node) {
	sd.branches = append(sd.branches, n)
}

func (sd *stateData) popBranch() *game.Node {
	parent := sd.branches[len(sd.branches)-1]
	sd.branches = sd.branches[:len(sd.branches)-1]
	return parent
}

func (sd *stateData) addToBuf(c rune) {
	sd.lastCharInBuf = c
	sd.buf.WriteRune(c)
}

func (sd *stateData) flushBuf() string {
	sd.lastCharInBuf = '~'
	o := sd.buf.String()
	sd.buf = strings.Builder{}
	return o
}

// parseError creates a parsing error, which fives index, line, col and
// charector context
func (sd *stateData) parseError(msg string) error {
	return fmt.Errorf(
		"error during state %v, at index %v, line %v, column %v, curchar %v. %v",
		sd.curstate, sd.idx, sd.row, sd.col, sd.curchar, msg)
}

// propBuffer contains a buffer of property data that has yet to be flushed.
type propBuffer struct {
	prop     string
	propdata []string
}

func (b *propBuffer) flush(n *game.Node) {
	n.Properties[b.prop] = b.propdata
	// Any post-processing should happen here
	b.prop = ""
	b.propdata = []string{}
}

// special chars, used to delimit sections of the SGF.
const (
	lparen  rune = '('
	rparen  rune = ')'
	lbrace  rune = '['
	rbrace  rune = ']'
	scolon  rune = ';'
	newline rune = '\n'
)

// parseState is used to denate parseState
type parseState int

const (
	// In beginning, we have not yet begun parsing the SGF.
	beginningState parseState = iota

	// In property , we are looking for data to complete a property token, such as
	// AW or B.  A property is considered complete when we see a left brace '['.
	propertyState

	// In propData , we are looking for all data associated with a property
	// token. We are finished with the Property data when we see a right brace ']'
	propDataState

	// In between, we are not accumulating data, we are just trying to figure out
	// where to go next.  Thus, we could find a new property, we could find more
	// property data, or we colud find a new variation.
	betweenState
)

// Parse parses a game into a tree of moves.
func (p *Parser) Parse() (*game.Game, error) {
	g := game.New()
	sd := &stateData{}
	pbuf := &propBuffer{}

	for c, _, err := p.rdr.ReadRune(); err == nil; c, _, err = p.rdr.ReadRune() {
		sd.idx++
		sd.col++
		sd.curchar = c
		if sd.curchar == newline {
			sd.row++
			sd.col = 0
		}

		switch sd.curstate {
		case beginningState:
			if unicode.IsSpace(sd.curchar) {
				// We can safely ignore whitespace here.
			} else if sd.curchar == lparen {
				// (;AW[aw][bw]
				// ^
				sd.branches = append(sd.branches, g.Root)
			} else if sd.curchar == scolon {
				// (;AW[aw][bw]
				//  ^
				sd.curstate = betweenState
				sd.curnode = g.Root
			} else {
				return nil, sd.parseError("unexpected char")
			}

		case betweenState:
			if unicode.IsSpace(sd.curchar) {
				// We can safely ignore whitespace here.
			} else if unicode.IsUpper(sd.curchar) {
				// AW[aw][bw]
				// ^
				pbuf.flush(sd.curnode)
				sd.addToBuf(sd.curchar)
				sd.curstate = propertyState
			} else if sd.curchar == lbrace {
				// AW[aw][bw]
				//   ^   ^
				sd.curstate = propDataState
			} else if sd.curchar == lparen {
				// AW[aw][bw] (;B[ab]
				//            ^
				pbuf.flush(sd.curnode)
				sd.addBranch(sd.curnode)
			} else if sd.curchar == scolon {
				// AW[aw][bw] (;B[ab];W[ac])
				//             ^     ^
				pbuf.flush(sd.curnode)
				cn := sd.curnode
				sd.curnode = game.NewNode()
				cn.AddChild(sd.curnode)
			} else if sd.curchar == rparen {
				// AW[aw][bw] (;B[ab])
				//                   ^
				pbuf.flush(sd.curnode)
				sd.curnode = sd.popBranch()
			} else {
				return nil, sd.parseError("unknown token")
			}

		case propertyState:
			if unicode.IsUpper(sd.curchar) {
				// AW[aw][bw]
				//  ^
				sd.addToBuf(sd.curchar)
			} else if sd.curchar == lbrace {
				// AW[aw][bw]
				//   ^
				prop := sd.flushBuf()
				// TODO(kashomon): Add validation to ignore invalid properties
				pbuf.prop = prop
				sd.curstate = propDataState
			} else {
				return nil, sd.parseError("unexpected character")
			}

		case propDataState:
			// C[foo 1[k\] bar]
			//           ^
			if sd.curchar == rbrace &&
				sd.lastCharInBuf == '\\' {
				sd.addToBuf(sd.curchar)
			} else if sd.curchar == rbrace {
				// C[foo 1[k\] bar]
				//                ^
				pbuf.propdata = append(pbuf.propdata, sd.flushBuf())
				sd.curstate = betweenState
			} else {
				// C[foo 1[k\] bar]
				//    ^
				sd.addToBuf(sd.curchar)
			}
			break

		default:
			break
			// return nil, sd.parseError("unexpected state -- couldn't match state")
		}
	}
	return g, nil
}
