package sgf

import (
	"errors"
	"io"
)

// the Parser uses a scanner to construct an Sgf object (the Game attribute)
// By including the token and save attributes, we can "cache" one token
// and "unscan" once
type Parser struct {
	Game    *Sgf
	scanner *Scanner
	token   *Token
	save    bool
}

// create a new Parser struct, must initialize the Game attribute as well
func NewParser(r io.Reader) *Parser {
	game := NewSgf()
	return &Parser{
		Game:    game,
		scanner: NewScanner(r),
	}
}

// scan one token using the scanner
func (p *Parser) scan() *Token {
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
func (p *Parser) scanSkipWhitespace() *Token {
	tok := p.scan()
	if tok.Type == Whitespace {
		tok = p.scan()
	}
	return tok
}

// TODO: currently only parses a single branch
func (p *Parser) Parse() (*Sgf, error) {
	err := p.parseBranch()
	if err != nil {
		return nil, err
	}
	return p.Game, nil
}

// parse, for example "PW[Player White]" into
// field="PW"
// value="Player White"
func (p *Parser) parseFieldValue() (string, string, error) {
	var tok *Token
	var field string
	var value string
	// parse a string
	// TODO: better erroring
	if tok = p.scanSkipWhitespace(); tok.Type != String {
		return "", "", errors.New("Corrupted sgf: 3")
	} else {
		field = tok.Raw
	}

	// parse a left bracket
	// TODO: better erroring
	if tok = p.scanSkipWhitespace(); tok.Type != LeftBracket {
		return "", "", errors.New("Corrupted sgf: 4")
	}

	// parse anything until a right bracket
	for {
		tok = p.scan()
		if tok.Type == Eof {
			return "", "", errors.New("EOF")
		}
		if tok.Type == RightBracket {
			break
		}
		value += tok.Raw
	}

	return field, value, nil
}

// parse a branch starting with "("
func (p *Parser) parseBranch() error {
	if tok := p.scanSkipWhitespace(); tok.Type != LeftParen {
		return errors.New("Corrupted sgf")
	}

	for {
		node := p.parseNode()
		if node == nil {
			break
		}
		p.Game.Nodes[p.Game.Index] = node
		p.Game.Index++
	}

	/*
	   do stuff here
	*/

	if tok := p.scanSkipWhitespace(); tok.Type != RightParen {
		return errors.New("Corrupted sgf")
	}

	return nil
}

type Node struct {
	Fields map[string]string
}

func NewNode() *Node {
	n := &Node{}
	n.Fields = make(map[string]string)
	return n
}

func (p *Parser) parseNode() *Node {
	if tok := p.scanSkipWhitespace(); tok.Type != Semicolon {
		p.unscan()
		return nil
	}

	node := NewNode()

	for {
		if tok := p.scanSkipWhitespace(); tok.Type == Semicolon || tok.Type == RightParen || tok.Type == LeftParen {
			p.unscan()
			break
		} else {
			p.unscan()
			field, value, err := p.parseFieldValue()
			if err == nil {
				node.Fields[field] = value
			}
		}
	}

	return node
}
