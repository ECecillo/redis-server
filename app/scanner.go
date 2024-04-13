package main // Scanner represents a lexical scanner.

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) scanCariageReturnLineFeed() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent CRLF character into the buffer.
	// Non-CRLF characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isCariageReturnLineFeed(ch) {
			s.unread()
			break
		} else if !isLineFeed(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return CRLF, buf.String()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isCariageReturnLineFeed(ch) {
		s.unread()
		return s.scanCariageReturnLineFeed()
	} else if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '+':
		return PLUS, string(ch)
	case '-':
		return MINUS, string(ch)
	case ':':
		return COLON, string(ch)
	case '$':
		return DOLLAR, string(ch)
	case '*':
		return ASTERISK, string(ch)
	case '_':
		return UNDERSCORE, string(ch)
	case '#':
		return HASH, string(ch)
	case ',':
		return COMMA, string(ch)
	case '(':
		return LPAREN, string(ch)
	case '!':
		return EXCL, string(ch)
	case '=':
		return EQ, string(ch)
	case '%':
		return PERCENT, string(ch)
	case '~':
		return TILD, string(ch)
	case '>':
		return GREATER, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "HELLO":
		return HELLO, buf.String()
	case "PING":
		return PING, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}
