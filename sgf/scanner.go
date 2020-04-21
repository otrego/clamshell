package sgf

import (
	"bufio"
	"bytes"
	"io"
)

const (
	Eof             = iota
	Whitespace
	LeftParen
    RightParen
    LeftBracket
    RightBracket
    Semicolon
    Backslash
    String
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isSpecial(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == ';'
}

type Token struct {
    Type    int
    Raw     string
}

type Scanner struct {
	*bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	s.UnreadRune()
}

func (s *Scanner) Scan() *Token {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if !isSpecial(ch) && ch != eof {
		s.unread()
		return s.scanString()
	}

	switch ch {
	case '(':
		return &Token{
            Type: LeftParen,
            Raw: "(",
        }
	case ')':
		return &Token{
            Type: RightParen,
            Raw: ")",
        }
	case '[':
		return &Token{
            Type: LeftBracket,
            Raw: "[",
        }
	case ']':
		return &Token{
            Type: RightBracket,
            Raw: "]",
        }
	case ';':
		return &Token{
            Type: Semicolon,
            Raw: ";",
        }
    default:
        return &Token{
            Type: Eof,
            Raw: "",
        }
	}
}

func (s *Scanner) scanWhitespace() *Token {
	buf := new(bytes.Buffer)
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return &Token{
        Type: Whitespace,
        Raw: buf.String(),
    }
}

func (s *Scanner) scanString() *Token {
	buf := new(bytes.Buffer)
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isSpecial(ch) {
			s.unread()
			break
		} else if ch == '\\' {
            buf.WriteRune(ch)
            buf.WriteRune(s.read())
        } else {
			buf.WriteRune(ch)
		}
	}
    return &Token{
        Type: String,
        Raw: buf.String(),
    }
}
