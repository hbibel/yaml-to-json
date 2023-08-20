package yaml

import "strings"

// TODO more symbols like '>-', '|' etc

type TokenKind int

const (
	INDENT TokenKind = iota
	WORD
	SPACE
	DASH
	NEWLINE
	COLON
	DOUBLE_QUOTE
	SINGLE_QUOTE
)

type Token interface {
	Kind() TokenKind
	String() string
}

type symbolicToken struct {
	kind    TokenKind
	content string
}
type indentToken struct {
	spaceCount uint32
}
type wordToken struct {
	content string
}
type spaceToken struct {
	content string
}

func (t *symbolicToken) Kind() TokenKind {
	return t.kind
}

func (t *indentToken) Kind() TokenKind {
	return INDENT
}

func (t *wordToken) Kind() TokenKind {
	return WORD
}

func (t *spaceToken) Kind() TokenKind {
	return SPACE
}

func (t *symbolicToken) String() string {
	return t.content
}

func (t *indentToken) String() string {
	sb := strings.Builder{}
	var i uint32
	for i = 0; i < t.spaceCount; i++ {
		sb.WriteByte(' ')
	}
	return sb.String()
}

func (t *wordToken) String() string {
	return t.content
}

func (t *spaceToken) String() string {
	return t.content
}

var newlineToken = &symbolicToken{NEWLINE, "\n"}
var dashToken Token = &symbolicToken{DASH, "-"}
var colonToken = &symbolicToken{COLON, ":"}
var doubleQuoteToken = &symbolicToken{DOUBLE_QUOTE, "\""}
var singleQuoteToken = &symbolicToken{SINGLE_QUOTE, "'"}
