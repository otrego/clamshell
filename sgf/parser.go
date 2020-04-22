package sgf

import (
	"errors"
	"io"
)

// A node is comprised of key-value pairs (might hold other info someday)
type Node struct {
	Fields map[string]string
}

func NewNode() *Node {
	n := &Node{}
	n.Fields = make(map[string]string)
	return n
}

// A game record (sgf) struct for now is just a collection of nodes
// TODO: this struct should retain the branch structure of game variations
type Sgf struct {
	Index int
	Nodes map[int]*Node
}

func NewSgf() *Sgf {
	s := &Sgf{Index: 0}
	s.Nodes = make(map[int]*Node)

	return s
}

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

// Parses the root branch
func (p *Parser) Parse() (*Sgf, error) {
	if tok := p.scanSkipWhitespace(); tok.Type != LeftParen {
		return nil, errors.New("Corrupted sgf: must start with '('")
	}
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

// parseBranch only gets called right after we consumed a "("
func (p *Parser) parseBranch() error {
	// loop through looking for nodes and branches
	for {
		// scan a token
		tok := p.scanSkipWhitespace()
		switch tok.Type {
		// if it's a semicolon, parse a node
		case Semicolon:
			node := p.parseNode()
			p.Game.Nodes[p.Game.Index] = node
			p.Game.Index++
		// if it's a left paren, parse a branch (recursive)
		case LeftParen:
			p.parseBranch()
		// if it's a right paren, return
		case RightParen:
			return nil
		// otherwise throw an error
		default:
			return errors.New("Corrupted sgf: error parsing branch")
		}
	}
	return nil
}

// parseNode only gets called right after we consumed a ";"
func (p *Parser) parseNode() *Node {
	node := NewNode()

	// the data in a node is pairs of fields and values
	for {
		tok := p.scanSkipWhitespace()
		p.unscan()
		// a semicolon, rightparen, or leftparen all end the node
		if tok.Type == Semicolon || tok.Type == RightParen || tok.Type == LeftParen {
			break
			// otherwise, parse a field and value
		} else {
			field, value, err := p.parseFieldValue()
			if err == nil {
				node.Fields[field] = value
			}
		}
	}

	return node
}
