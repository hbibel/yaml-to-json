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

			var numSpaces uint32
			remaining, numSpaces = countLeadingSpaces(remaining)
			if numSpaces > 0 {
				tokens <- &indentToken{numSpaces}
			}

			for len(remaining) > 0 {

				var space string
				remaining, space = getLeadingSpaces(remaining)
				if len(space) > 0 {
					tokens <- &spaceToken{space}
					// it's strictly not necessary to continue here, but the code is more
					// consistent this way
					continue
				}

				remaining, ok = tryParseSymbol([]rune{'-'}, remaining)
				if ok {
					tokens <- dashToken
					continue
				}

				remaining, ok = tryParseSymbol([]rune{':'}, remaining)
				if ok {
					tokens <- colonToken
					continue
				}

				remaining, ok = tryParseSymbol([]rune{'"'}, remaining)
				if ok {
					tokens <- doubleQuoteToken
					continue
				}

				remaining, ok = tryParseSymbol([]rune{'\''}, remaining)
				if ok {
					tokens <- singleQuoteToken
					continue
				}

				var word string
				remaining, word = getNextWord(remaining)
				if len(word) > 0 {
					tokens <- &wordToken{word}
				}
			}

			tokens <- newlineToken
		}
		close(tokens)
	}()

}

func countLeadingSpaces(runes []rune) ([]rune, uint32) {
	var numSpaces uint32 = 0
	for _, char := range runes {
		if char == ' ' {
			numSpaces++
		} else {
			break
		}
	}

	return runes[numSpaces:], numSpaces
}

func getLeadingSpaces(runes []rune) ([]rune, string) {
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
	return runes[firstNonSpacePos:], literal
}

func tryParseSymbol(symbol []rune, runes []rune) ([]rune, bool) {
	symbolLength := len(symbol)
	if len(runes) < symbolLength {
		return runes, false
	}

	for i, symbolRune := range symbol {
		if symbolRune != runes[i] {
			return runes, false
		}
	}

	return runes[symbolLength:], true
}

func getNextWord(runes []rune) ([]rune, string) {
	literalBuilder := strings.Builder{}
	var pos int
	var c rune
	for pos, c = range runes {
		if isSpace(c) || isSpecial(c) {
			break
		}
		literalBuilder.WriteRune(c)
	}
	word := literalBuilder.String()
	return runes[pos:], word
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
