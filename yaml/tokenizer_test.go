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

	yaml.Tokenize(lines, tokens)

	expected := []kindAndContent{}
	failIfUnexpected(t, expected, tokens)
	close(lines)
}

func TestTokenizeEmptyLine(t *testing.T) {
	lines := make(chan string)
	tokens := make(chan yaml.Token)

	yaml.Tokenize(lines, tokens)

	expected := []kindAndContent{}
	failIfUnexpected(t, expected, tokens)

	lines <- ""
	close(lines)
}

func failIfUnexpected(t *testing.T, expected []kindAndContent, tokens <-chan yaml.Token) {
	go func() {
		i := 0
		for token := range tokens {
			if i >= len(expected) {
				t.Errorf("Too many tokens: %v", token)
			}

			if expected[i].kind != token.Kind() {
				t.Errorf("Unexpected token kind: %v", token)
			}
			if expected[i].content != token.String() {
				t.Errorf("Unexpected token content: %v", token)
			}
			i++
		}
	}()
}

// func TestTokenizeEmptyChunk(t *testing.T) {
// 	chunks := make(chan string)
// 	tokens := make(chan yaml.Token)

// 	yaml.Tokenize(chunks, tokens)
// 	go func() {
// 		for token := range tokens {
// 			t.Errorf("Unexpected token: %v", token)
// 		}
// 	}()

// 	close(chunks)
// }

// func TestTokenizeEmptyString(t *testing.T) {
// 	chunks := make(chan string)
// 	tokens := make(chan yaml.Token)

// 	yaml.Tokenize(chunks, tokens)
// 	go func() {
// 		for token := range tokens {
// 			t.Errorf("Unexpected token: %v", token)
// 		}
// 	}()

// 	chunks <- ""
// 	close(chunks)
// }

// func TestTokenizeSingleToken(t *testing.T) {
// 	chunks := make(chan string)
// 	tokens := make(chan yaml.Token)
// 	actualTokens := make([]yaml.Token, 0)

// 	yaml.Tokenize(chunks, tokens)
// 	go func() {
// 		for token := range tokens {
// 			actualTokens = append(actualTokens, token)
// 		}
// 	}()

// 	chunks <- "-"
// 	close(chunks)

// 	expectedTokens := []yaml.Token{
// 		{yaml.DASH, ""},
// 	}
// 	for i, expected := range expectedTokens {
// 		actual := actualTokens[i]
// 		if actual != expected {
// 			t.Errorf("Expected token %v but got %v", expected, actual)
// 		}
// 	}
// }

// func TestMultipleTokens(t *testing.T) {
// 	chunks := make(chan string)
// 	tokens := make(chan yaml.Token)
// 	actualTokens := make([]yaml.Token, 0)

// 	yaml.Tokenize(chunks, tokens)
// 	go func() {
// 		for token := range tokens {
// 			actualTokens = append(actualTokens, token)
// 		}
// 	}()

// 	chunks <- "key: value"
// 	close(chunks)

// 	expectedTokens := []yaml.Token{
// 		{yaml.WORD, "key"},
// 		{yaml.COLON, ""},
// 		{yaml.WORD, "value"},
// 	}
// 	for i, expected := range expectedTokens {
// 		actual := actualTokens[i]
// 		if actual != expected {
// 			t.Errorf("Expected token %v but got %v", expected, actual)
// 		}
// 	}
// }
