package scanner

import (
	"bufio"
	"bytes"
	"io"
)

type ControlChar int

// constants
const (
	EOF ControlChar = iota
	Whitespace
	LeftParen
	RightParen
	LeftBracket
	RightBracket
	Semicolon
	Backslash
	String
	Illegal
)

var eof = rune(0)

// isWhitespace is a check to see if a rune is whitespace
// (space, tab, newline, carriage return)
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// isSpecial is a check for special characters used in SGF files
func isSpecial(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == ';'
}

// Token contains attributes for the Type of token and the raw string
type Token struct {
	Type ControlChar
	Raw  string
}

// Scanner wraps bufio.Reader but has the public Scan() function to scan tokens
type Scanner struct {
	*bufio.Reader
}

// New creates a scanner object out of an io.Reader interface
func New(r io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(r)}
}

// Either return the current rune or EOF (null)
func (s *Scanner) read() rune {
	ch, _, err := s.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// Go back one rune
func (s *Scanner) unread() {
	s.UnreadRune()
}

// Scan is the function that tokenizes the input and outputs tokens
func (s *Scanner) Scan() *Token {
	// read in a character
	ch := s.read()

	// if it's whitespace
	if isWhitespace(ch) {
		// then go back and scan as much whitespace as possible
		s.unread()
		return s.scanWhitespace()
		// if it's not a special character or EOF, then it's a "string"
	} else if !isSpecial(ch) && ch != eof {
		// go back and scan as much of the string as possible
		s.unread()
		return s.scanString()
	}

	// otherwise, handle the special characters
	switch ch {
	case '(':
		return &Token{
			Type: LeftParen,
			Raw:  "(",
		}
	case ')':
		return &Token{
			Type: RightParen,
			Raw:  ")",
		}
	case '[':
		return &Token{
			Type: LeftBracket,
			Raw:  "[",
		}
	case ']':
		return &Token{
			Type: RightBracket,
			Raw:  "]",
		}
	case ';':
		return &Token{
			Type: Semicolon,
			Raw:  ";",
		}
	case eof:
		return &Token{
			Type: EOF,
			Raw:  "",
		}
	default:
		return &Token{
			Type: Illegal,
			Raw:  "",
		}
	}
}

// scanWhiteSpace consumes as much as whitespace as possible
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
		Raw:  buf.String(),
	}
}

// scanString consumes as much (not whitespace, not special, not eof) as possible
func (s *Scanner) scanString() *Token {
	buf := new(bytes.Buffer)
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isSpecial(ch) {
			s.unread()
			break
			// backslash is sort of a special case in that we will consume
			// the next character *no matter what* (i.e., even if it's a special character)
		} else if ch == '\\' {
			buf.WriteRune(ch)
			buf.WriteRune(s.read())
		} else {
			buf.WriteRune(ch)
		}
	}
	return &Token{
		Type: String,
		Raw:  buf.String(),
	}
}
