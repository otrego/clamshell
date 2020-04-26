package parser

import (
	"errors"
	"github.com/otrego/clamshell/game"
	"github.com/otrego/clamshell/scanner"
	"io"
	"strings"
)

// Parser uses a scanner to construct an Sgf object
// By including the token and save attributes, we can "cache" one token
// and "unscan" once
type Parser struct {
	scanner *scanner.Scanner
	token   *scanner.Token
	save    bool
}

// FromString creates a Parser using a string as input
func FromString(s string) *Parser {
	r := strings.NewReader(s)
	return FromReader(r)
}

// FromReader creates a Parser struct using io.reader as input
func FromReader(r io.Reader) *Parser {
	return &Parser{
		scanner: scanner.New(r),
	}
}

// scan one token using the scanner
func (p *Parser) scan() *scanner.Token {
	// the save bool is used to "rewind" once
	if p.save {
		p.save = false
		return p.token
	}
	// otherwise just scan normally
	p.token = p.scanner.Scan()
	return p.token
}

// unscan lets the Parser know to use the saved token on the next scan
// instead of scanning for a new token
func (p *Parser) unscan() {
	p.save = true
}

// scanSkipWhitespace is a wrapper around scan() that skips whitespace
func (p *Parser) scanSkipWhitespace() *scanner.Token {
	tok := p.scan()
	if tok.Type == scanner.Whitespace {
		tok = p.scan()
	}
	return tok
}

// Parse parses the root branch
func (p *Parser) Parse() (*game.Game, error) {
	if tok := p.scanSkipWhitespace(); tok.Type != scanner.LeftParen {
		return nil, errors.New("Corrupted sgf: must start with '('")
	}
	g := game.New()
	err := p.parseBranch(g.Root)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// parse, for example "PW[Player White]" into
// field="PW"
// value="Player White"
func (p *Parser) parseProperties() (string, []string, error) {
	var tok *scanner.Token
	var field string
	values := []string{}
	//var value string
	// parse a string
	// TODO: better erroring
	if tok = p.scanSkipWhitespace(); tok.Type != scanner.String {
		return "", []string{}, errors.New("Error parsing properties")
	}
	field = tok.Raw

	// can have multiple values
	for {

		// parse a left bracket
		// TODO: better erroring
		if tok = p.scanSkipWhitespace(); tok.Type != scanner.LeftBracket {
			p.unscan()
			break
		}
		value := ""

		// parse anything until a right bracket
		for {
			tok = p.scan()
			if tok.Type == scanner.EOF {
				return "", []string{}, errors.New("EOF")
			}
			if tok.Type == scanner.RightBracket {
				break
			}
			value += tok.Raw
		}
		values = append(values, value)
	}

	return field, values, nil
}

// parseBranch only gets called right after we consumed a "("
func (p *Parser) parseBranch(cur *game.Node) error {
	// loop through looking for nodes and branches
	for {
		// scan a token
		tok := p.scanSkipWhitespace()
		switch tok.Type {
		// if it's a semicolon, parse a node
		case scanner.Semicolon:
			node := p.parseNode()
			cur.Children = append(cur.Children, node)
			node.Parent = cur
			cur = node
		// if it's a left paren, parse a branch (recursive)
		case scanner.LeftParen:
			if err := p.parseBranch(cur); err != nil {
				return err
			}
		// if it's a right paren, return
		case scanner.RightParen:
			return nil
		// otherwise throw an error
		default:
			return errors.New("Corrupted sgf: error parsing branch")
		}
	}
	return nil
}

// parseNode only gets called right after we consumed a ";"
func (p *Parser) parseNode() *game.Node {
	node := game.NewNode()

	// the data in a node is pairs of fields and values
	for {
		tok := p.scanSkipWhitespace()
		p.unscan()
		// a semicolon, rightparen, or leftparen all end the node
		if tok.Type == scanner.Semicolon || tok.Type == scanner.RightParen || tok.Type == scanner.LeftParen {
			break
			// otherwise, parse a field and value
		} else {
			field, values, err := p.parseProperties()
			if err == nil {
				if v := node.Properties[field]; v == nil {
					node.Properties[field] = []string{}
				}
				node.Properties[field] = append(node.Properties[field], values...)
			}
		}
	}

	return node
}
