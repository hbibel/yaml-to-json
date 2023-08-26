package json

import (
	"hbibel/yaml-to-json/common"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	events := []common.Event{}
	expected := []string{}
	runTest(t, events, expected)
}

func TestStringEvent(t *testing.T) {
	events := []common.Event{
		common.NewStringEvent("foo"),
	}
	expected := []string{
		"\"foo\"",
	}
	runTest(t, events, expected)
}

func TestNumberEvent(t *testing.T) {
	events := []common.Event{
		common.NewNumberEvent("42"),
	}
	expected := []string{
		"42",
	}
	runTest(t, events, expected)
}

func TestBooleanEvent(t *testing.T) {
	events := []common.Event{
		common.NewBooleanEvent("true"),
	}
	expected := []string{
		"true",
	}
	runTest(t, events, expected)
}

func TestNullEvent(t *testing.T) {
	events := []common.Event{
		common.NewNullEvent(),
	}
	expected := []string{
		"null",
	}
	runTest(t, events, expected)
}

func TestMapping(t *testing.T) {
	events := []common.Event{
		common.NewStartMappingEvent(),
		common.NewKeyEvent("foo"),
		common.NewStringEvent("bar"),
		common.NewKeyEvent("baz"),
		common.NewStringEvent("qux"),
		common.NewEndMappingEvent(),
	}
	expected := []string{
		"{",
		"\"foo\"",
		":",
		"\"bar\"",
		",",
		"\"baz\"",
		":",
		"\"qux\"",
		"}",
	}
	runTest(t, events, expected)
}

func TestArray(t *testing.T) {
	events := []common.Event{
		common.NewStartArrayEvent(),
		common.NewEmitElementEvent(),
		common.NewStringEvent("foo"),
		common.NewEmitElementEvent(),
		common.NewStringEvent("bar"),
		common.NewEndArrayEvent(),
	}
	expected := []string{
		"[",
		"\"foo\"",
		",",
		"\"bar\"",
		"]",
	}
	runTest(t, events, expected)
}

func TestNestedMapping(t *testing.T) {
	events := []common.Event{
		common.NewStartMappingEvent(),
		common.NewKeyEvent("foo"),
		common.NewStartMappingEvent(),
		common.NewKeyEvent("bar"),
		common.NewStringEvent("baz"),
		common.NewEndMappingEvent(),
		common.NewEndMappingEvent(),
	}
	expected := []string{
		"{",
		"\"foo\"",
		":",
		"{",
		"\"bar\"",
		":",
		"\"baz\"",
		"}",
		"}",
	}
	runTest(t, events, expected)
}

func TestArrayOfMappings(t *testing.T) {
	events := []common.Event{
		common.NewStartArrayEvent(),
		common.NewEmitElementEvent(),
		common.NewStartMappingEvent(),
		common.NewKeyEvent("foo"),
		common.NewStringEvent("bar"),
		common.NewEndMappingEvent(),
		common.NewEndArrayEvent(),
	}
	expected := []string{
		"[",
		"{",
		"\"foo\"",
		":",
		"\"bar\"",
		"}",
		"]",
	}
	runTest(t, events, expected)
}

func runTest(t *testing.T, events []common.Event, expectedChunks []string) {
	eventsChannel := make(chan common.Event)
	done := make(chan bool)

	chunkChannel := RenderEvents(eventsChannel)
	var chunks = make([]string, 0)
	go func() {
		for chunk := range chunkChannel {
			chunks = append(chunks, chunk)
		}
		done <- true
	}()

	for _, event := range events {
		eventsChannel <- event
	}
	close(eventsChannel)

	<-done
	if !reflect.DeepEqual(chunks, expectedChunks) {
		t.Error("Expected", expectedChunks, "got", chunks)
	}
}
