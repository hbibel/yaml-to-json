package yaml

import (
	"testing"
)

type kindAndContent struct {
	kind    TokenKind
	content string
}

func TestTokenizeNoLines(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	expected := []kindAndContent{}
	failIfUnexpected(t, expected, tokens, done)
	close(lines)
}

func TestTokenizeEmptyLine(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	expected := []kindAndContent{
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	lines <- ""
	close(lines)
}

func TestTokenizeSingleToken(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	lines <- "-"
	expected := []kindAndContent{
		{DASH, "-"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeConsecutiveTokens(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	lines <- "-:"
	expected := []kindAndContent{
		{DASH, "-"},
		{COLON, ":"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeLeadingSpaces(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	lines <- "  -"
	expected := []kindAndContent{
		{INDENT, "  "},
		{DASH, "-"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeAlphaWord(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	lines <- "key"
	expected := []kindAndContent{
		{WORD, "key"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeMultipleLines(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	input := []string{
		"  -",
		"  -",
	}
	expected := []kindAndContent{
		{INDENT, "  "},
		{DASH, "-"},
		{NEWLINE, "\n"},
		{INDENT, "  "},
		{DASH, "-"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func TestRepeatSymbolicToken(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	input := []string{
		"--",
	}
	expected := []kindAndContent{
		{DASH, "-"},
		{DASH, "-"},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func TestSmallYamlFile(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan Token)
	done := make(chan bool)
	defer func() { <-done }()

	Tokenize(lines, tokens)

	input := []string{
		"key: value",
		"key2: ",
		"  - 'x'",
		"  - y",
		"  - \"z\"",
	}
	expected := []kindAndContent{
		{WORD, "key"},
		{COLON, ":"},
		{SPACE, " "},
		{WORD, "value"},
		{NEWLINE, "\n"},
		{WORD, "key2"},
		{COLON, ":"},
		{SPACE, " "},
		{NEWLINE, "\n"},
		{INDENT, "  "},
		{DASH, "-"},
		{SPACE, " "},
		{SINGLE_QUOTE, "'"},
		{WORD, "x"},
		{SINGLE_QUOTE, "'"},
		{NEWLINE, "\n"},
		{INDENT, "  "},
		{DASH, "-"},
		{SPACE, " "},
		{WORD, "y"},
		{NEWLINE, "\n"},
		{INDENT, "  "},
		{DASH, "-"},
		{SPACE, " "},
		{DOUBLE_QUOTE, "\""},
		{WORD, "z"},
		{DOUBLE_QUOTE, "\""},
		{NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func failIfUnexpected(t *testing.T, expected []kindAndContent, tokens <-chan Token, done chan<- bool) {
	go func() {
		actual := []kindAndContent{}
		for token := range tokens {
			actual = append(actual, kindAndContent{token.Kind(), token.String()})
		}

		if len(actual) != len(expected) {
			t.Errorf("\nActual: %v\nExpected: %v", actual, expected)
		} else {
			for i, actualToken := range actual {
				if expected[i].kind != actualToken.kind {
					t.Errorf("Unexpected token kind: %v", actualToken)
				}
				if expected[i].content != actualToken.content {
					t.Errorf("Unexpected token content: '%v'", actualToken)
				}
			}
		}

		done <- true
	}()
}
