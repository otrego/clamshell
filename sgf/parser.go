package sgf

import (
	"errors"
	"io"
)

type Parser struct {
    Game    *Sgf
	scanner *Scanner
    token   *Token
    save    bool
}

func NewParser(r io.Reader) *Parser {
    game := NewSgf()
    return &Parser{
        Game: game,
        scanner: NewScanner(r),
    }
}

func (p *Parser) scan() *Token {
	if p.save {
		p.save = false
		return p.token
	}
	p.token = p.scanner.Scan()
	return p.token
}

func (p *Parser) unscan() {
	p.save = true
}

func (p *Parser) scanSkipWhitespace() *Token {
	tok := p.scan()
	if tok.Type == Whitespace {
		tok = p.scan()
	}
	return tok
}

func (p *Parser) Parse() (*Sgf, error) {
    err := p.parseBranch()
    if err != nil {
        return nil, err
    }
    return p.Game, nil
}

func (p *Parser) parseFieldValue() (string, string, error) {
    var tok *Token
    var field string
    var value string
	if tok = p.scanSkipWhitespace(); tok.Type != String {
        return "", "", errors.New("Corrupted sgf: 3")
	} else {
		field = tok.Raw
	}

	if tok = p.scanSkipWhitespace(); tok.Type != LeftBracket {
        return "", "", errors.New("Corrupted sgf: 4")
	}

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
    Fields      map[string]string
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

