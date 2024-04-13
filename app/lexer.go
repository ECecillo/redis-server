package main

type Token int

// Ressource : https://blog.gopheracademy.com/advent-2014/parsers-lexers/
// Ref : https://redis.io/docs/latest/develop/reference/protocol-spec/
const (
	ILLEGAL Token = iota

	CR   // \r
	LF   // \n
	CRLF // \r\n

	WS  // " "
	EOF // \0

	IDENT // document key

	// TYPES

	PLUS       // + : simple strings
	MINUS      // - : Simple errors
	COLON      // : =>
	DOLLAR     // $ : Bulk Strings
	ASTERISK   // * : Arrays
	UNDERSCORE // _ : Nulls
	HASH       // # : Booleans
	COMMA      // , : doubles
	LPAREN     // ( : Big numbers
	EXCL       // ! : Bulk errors
	EQ         // = : Verbatim strings
	PERCENT    // % : Maps
	TILD       // ~ : Sets
	GREATER    // > : Pushes

	// COMMANDS
	HELLO // Used for Handshake Client-Server
	PING
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isCariageReturnLineFeed(ch rune) bool {
	return ch == '\r'
}

func isLineFeed(ch rune) bool {
	return ch == '\n'
}

var eof = rune(0)

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}
