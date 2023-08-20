package common

type EventType int

const (
	START_MAPPING EventType = iota
	EMIT_KEY
	EMIT_VALUE
	END_MAPPING
	START_ARRAY
	EMIT_ELEMENT
	END_ARRAY
)

type Event interface {
	GetKind() EventType
}

type PayLoadType int

const (
	STRING PayLoadType = iota
	NUMBER
	BOOLEAN
	NULL
)

type HasPayload interface {
	GetPayLoadType() PayLoadType
	GetPayload() string
}

type EventWithPayload struct {
	Kind        EventType
	PayloadType PayLoadType
	Payload     string
}

type eventWithoutPayload struct {
	Kind EventType
}

func (e *EventWithPayload) GetKind() EventType {
	return e.Kind
}

func (e *eventWithoutPayload) GetKind() EventType {
	return e.Kind
}

func (e *EventWithPayload) GetPayLoadType() PayLoadType {
	return e.PayloadType
}

func (e *EventWithPayload) GetPayload() string {
	return e.Payload
}

func NewStringEvent(payload string) Event {
	return &EventWithPayload{
		Kind:        EMIT_VALUE,
		PayloadType: STRING,
		Payload:     payload,
	}
}

func NewNumberEvent(payload string) Event {
	return &EventWithPayload{
		Kind:        EMIT_VALUE,
		PayloadType: NUMBER,
		Payload:     payload,
	}
}

func NewBooleanEvent(payload string) Event {
	return &EventWithPayload{
		Kind:        EMIT_VALUE,
		PayloadType: BOOLEAN,
		Payload:     payload,
	}
}

func NewNullEvent() Event {
	return &EventWithPayload{
		Kind:        EMIT_VALUE,
		PayloadType: NULL,
		Payload:     "null",
	}
}

func NewStartMappingEvent() Event {
	return &eventWithoutPayload{
		Kind: START_MAPPING,
	}
}

func NewKeyEvent(payload string) Event {
	return &EventWithPayload{
		Kind:        EMIT_KEY,
		PayloadType: STRING,
		Payload:     payload,
	}
}

func NewEndMappingEvent() Event {
	return &eventWithoutPayload{
		Kind: END_MAPPING,
	}
}

func NewStartArrayEvent() Event {
	return &eventWithoutPayload{
		Kind: START_ARRAY,
	}
}

func NewEndArrayEvent() Event {
	return &eventWithoutPayload{
		Kind: END_ARRAY,
	}
}

func NewEmitElementEvent() Event {
	return &eventWithoutPayload{
		Kind: EMIT_ELEMENT,
	}
}
