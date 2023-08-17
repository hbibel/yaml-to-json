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

func failIfUnexpected(t *testing.T, expected []kindAndContent, tokens <-chan yaml.Token, done chan<- bool) {
	go func() {
		i := 0
		for token := range tokens {
			if i >= len(expected) {
				t.Errorf("Too many tokens: {%v, '%v'}", token.Kind(), token.String())
				break
			}

			if expected[i].kind != token.Kind() {
				t.Errorf("Unexpected token kind: %v", token)
			}
			if expected[i].content != token.String() {
				t.Errorf("Unexpected token content: %v", token)
			}
			i++
		}
		done <- true
	}()
}
