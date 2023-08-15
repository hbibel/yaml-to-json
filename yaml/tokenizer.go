package yaml

import (
	"strings"
)

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
	for i = 0; i <= t.spaceCount; i++ {
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

func Tokenize(lines <-chan string, tokens chan<- Token) {
	go func() {
		for line := range lines {

			var ok bool
			var remaining []rune = []rune(line)
			remaining = parseIndent(remaining, tokens)

			for len(remaining) > 0 {

				remaining, ok = tryParseSpace(remaining, tokens)
				if ok {
					// it's strictly not necessary to continue here, but the code is more
					// consistent this way
					continue
				}

				remaining, ok = tryParseDash(remaining, tokens)
				if ok {
					continue
				}

				remaining, ok = tryParseColon(remaining, tokens)
				if ok {
					continue
				}

				remaining, ok = tryParseDoubleQuote(remaining, tokens)
				if ok {
					continue
				}

				remaining, ok = tryParseSingleQuote(remaining, tokens)
				if ok {
					continue
				}

				remaining = parseWord(remaining, tokens)
			}

			tokens <- newlineToken
		}
		close(tokens)
	}()

}

func parseIndent(runes []rune, tokens chan<- Token) []rune {
	var numSpaces uint32 = 0
	for _, char := range runes {
		if char == ' ' {
			numSpaces++
		} else {
			break
		}
	}

	if numSpaces > 0 {
		tokens <- &indentToken{numSpaces}
	}

	return runes[numSpaces:]
}

func tryParseSpace(runes []rune, tokens chan<- Token) ([]rune, bool) {
	literalBuilder := strings.Builder{}
	firstNonSpacePos := 0
	for pos, c := range runes {
		if !isSpace(c) {
			firstNonSpacePos = pos
			break
		}

		literalBuilder.WriteRune(c)
	}

	literal := literalBuilder.String()
	if literal == "" {
		return runes, false
	}

	tokens <- &spaceToken{literal}
	return runes[firstNonSpacePos:], true
}

func tryParseDash(runes []rune, tokens chan<- Token) ([]rune, bool) {
	if !(runes[0] == '-') {
		return runes, false
	}
	if len(runes) == 1 || !isSpace(runes[1]) {
		return runes, false
	}

	tokens <- dashToken
	return runes[2:], true
}

func tryParseColon(runes []rune, tokens chan<- Token) ([]rune, bool) {
	if !(runes[0] == ':') {
		return runes, false
	}
	if len(runes) > 1 && !isSpace(runes[1]) {
		return runes, false
	}

	tokens <- colonToken
	if len(runes) == 1 {
		return []rune{}, true
	}
	return runes[2:], true
}

func tryParseDoubleQuote(runes []rune, tokens chan<- Token) ([]rune, bool) {
	if !(runes[0] == '"') {
		return runes, false
	}

	tokens <- doubleQuoteToken
	return runes[1:], true
}

func tryParseSingleQuote(runes []rune, tokens chan<- Token) ([]rune, bool) {
	if !(runes[0] == '\'') {
		return runes, false
	}

	tokens <- singleQuoteToken
	return runes[1:], true
}

func parseWord(runes []rune, tokens chan<- Token) []rune {
	literalBuilder := strings.Builder{}
	for pos, c := range runes {
		if isSpace(c) || isSpecial(c) {
			tokens <- &wordToken{literalBuilder.String()}
			return runes[pos:]
		}
		literalBuilder.WriteRune(c)
	}
	return []rune{}
}

func isSpecial(c rune) bool {
	return (c == '\'' ||
		c == '"' ||
		c == ':' ||
		c == '-')
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t'
}
