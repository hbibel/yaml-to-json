package yaml_test

import (
	"hbibel/yaml-to-json/yaml"
	"testing"
)

type kindAndContent struct {
	kind    yaml.TokenKind
	content string
}

func TestTokenizeNoLines(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	expected := []kindAndContent{}
	failIfUnexpected(t, expected, tokens, done)
	close(lines)
}

func TestTokenizeEmptyLine(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	expected := []kindAndContent{
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	lines <- ""
	close(lines)
}

func TestTokenizeSingleToken(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	lines <- "-"
	expected := []kindAndContent{
		{yaml.DASH, "-"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeConsecutiveTokens(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	lines <- "-:"
	expected := []kindAndContent{
		{yaml.DASH, "-"},
		{yaml.COLON, ":"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeLeadingSpaces(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	lines <- "  -"
	expected := []kindAndContent{
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeAlphaWord(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	lines <- "key"
	expected := []kindAndContent{
		{yaml.WORD, "key"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	close(lines)
}

func TestTokenizeMultipleLines(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	input := []string{
		"  -",
		"  -",
	}
	expected := []kindAndContent{
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.NEWLINE, "\n"},
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func TestRepeatSymbolicToken(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	input := []string{
		"--",
	}
	expected := []kindAndContent{
		{yaml.DASH, "-"},
		{yaml.DASH, "-"},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func TestSmallYamlFile(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)
	done := make(chan bool)
	defer func() { <-done }()

	yaml.Tokenize(lines, tokens)

	input := []string{
		"key: value",
		"key2: ",
		"  - 'x'",
		"  - y",
		"  - \"z\"",
	}
	expected := []kindAndContent{
		{yaml.WORD, "key"},
		{yaml.COLON, ":"},
		{yaml.SPACE, " "},
		{yaml.WORD, "value"},
		{yaml.NEWLINE, "\n"},
		{yaml.WORD, "key2"},
		{yaml.COLON, ":"},
		{yaml.SPACE, " "},
		{yaml.NEWLINE, "\n"},
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.SPACE, " "},
		{yaml.SINGLE_QUOTE, "'"},
		{yaml.WORD, "x"},
		{yaml.SINGLE_QUOTE, "'"},
		{yaml.NEWLINE, "\n"},
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.SPACE, " "},
		{yaml.WORD, "y"},
		{yaml.NEWLINE, "\n"},
		{yaml.INDENT, "  "},
		{yaml.DASH, "-"},
		{yaml.SPACE, " "},
		{yaml.DOUBLE_QUOTE, "\""},
		{yaml.WORD, "z"},
		{yaml.DOUBLE_QUOTE, "\""},
		{yaml.NEWLINE, "\n"},
	}
	failIfUnexpected(t, expected, tokens, done)

	for _, line := range input {
		lines <- line
	}

	close(lines)
}

func failIfUnexpected(t *testing.T, expected []kindAndContent, tokens <-chan yaml.Token, done chan<- bool) {
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
