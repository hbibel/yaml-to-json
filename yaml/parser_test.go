package yaml

import (
	"hbibel/yaml-to-json/common"
	"reflect"
	"testing"
)

// TODO error cases

// TODO cases to test
//   List of lists
//   "a number 42 within a string"
//   "42 starts a string"
//   "a normal string - but with a dash"

func TestTokensToEventsNoTokens(t *testing.T) {
	tokens := []Token{}
	expectedEvents := []common.Event{}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsStringToken(t *testing.T) {
	tokens := []Token{
		&wordToken{"foo"},
	}
	expectedEvents := []common.Event{
		common.NewStringEvent("foo"),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsIntegerToken(t *testing.T) {
	tokens := []Token{
		&wordToken{"42"},
	}
	expectedEvents := []common.Event{
		common.NewNumberEvent("42"),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsFloatToken(t *testing.T) {
	tokens := []Token{
		&wordToken{"42.0"},
	}
	expectedEvents := []common.Event{
		common.NewNumberEvent("42.0"),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsBooleanToken(t *testing.T) {
	tokens := []Token{
		&wordToken{"true"},
	}
	expectedEvents := []common.Event{
		common.NewBooleanEvent("true"),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsNullToken(t *testing.T) {
	tokens := []Token{
		&wordToken{"null"},
	}
	expectedEvents := []common.Event{
		common.NewNullEvent(),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsMultiWordString(t *testing.T) {
	tokens := []Token{
		&wordToken{"foo"},
		&spaceToken{" "},
		&wordToken{"bar"},
	}
	expectedEvents := []common.Event{
		common.NewStringEvent("foo bar"),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsSingleItemList(t *testing.T) {
	tokens := []Token{
		dashToken,
		&spaceToken{" "},
		&wordToken{"foo"},
	}
	expectedEvents := []common.Event{
		common.NewStartArrayEvent(),
		common.NewEmitElementEvent(),
		common.NewStringEvent("foo"),
		common.NewEndArrayEvent(),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsListOfLists(t *testing.T) {
	tokens := []Token{
		dashToken,
		&spaceToken{" "},
		dashToken,
		&spaceToken{" "},
		&wordToken{"foo"},
		newlineToken,
		&indentToken{2},
		dashToken,
		&spaceToken{" "},
		&wordToken{"bar"},
	}
	expectedEvents := []common.Event{
		common.NewStartArrayEvent(),
		common.NewEmitElementEvent(),
		common.NewStartArrayEvent(),
		common.NewEmitElementEvent(),
		common.NewStringEvent("foo"),
		common.NewEmitElementEvent(),
		common.NewStringEvent("bar"),
		common.NewEndArrayEvent(),
		common.NewEndArrayEvent(),
	}
	runTest(t, tokens, expectedEvents)
}

func TestTokensToEventsSimpleMap(t *testing.T) {
	tokens := []Token{
		&wordToken{"foo"},
		colonToken,
		newlineToken,
		&indentToken{2},
		&wordToken{"bar"},
	}
	expectedEvents := []common.Event{
		common.NewStartMappingEvent(),
		common.NewKeyEvent("foo"),
		common.NewStringEvent("bar"),
		common.NewEndMappingEvent(),
	}
	runTest(t, tokens, expectedEvents)
}

func runTest(t *testing.T, tokens []Token, expectedEvents []common.Event) {
	tokenChannel := make(chan Token)
	done := make(chan bool)

	eventChannel := TokensToEvents(tokenChannel)
	var events = make([]common.Event, 0)
	go func() {
		for event := range eventChannel {
			events = append(events, event)
		}
		done <- true
	}()

	for _, token := range tokens {
		tokenChannel <- token
	}
	close(tokenChannel)

	<-done
	if !reflect.DeepEqual(events, expectedEvents) {
		t.Error("Expected", expectedEvents, "got", events)
	}
}
