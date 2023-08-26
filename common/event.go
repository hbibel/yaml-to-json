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

func (e *eventWithoutPayload) GetKind() EventType {
	return e.Kind
}

func (e *eventWithoutPayload) String() string {
	switch e.Kind {
	case START_MAPPING:
		return "<START_MAPPING>"
	case EMIT_KEY:
		return "<EMIT_KEY>"
	case EMIT_VALUE:
		return "<EMIT_VALUE>"
	case END_MAPPING:
		return "<END_MAPPING>"
	case START_ARRAY:
		return "<START_ARRAY>"
	case EMIT_ELEMENT:
		return "<EMIT_ELEMENT>"
	case END_ARRAY:
		return "<END_ARRAY>"
	default:
		return "<UNKNOWN>"
	}
}

func (e *EventWithPayload) GetKind() EventType {
	return e.Kind
}

func (e *EventWithPayload) GetPayload() string {
	return e.Payload
}

func (e *EventWithPayload) GetPayLoadType() PayLoadType {
	return e.PayloadType
}

func (e *EventWithPayload) String() string {
	switch e.Kind {
	case START_MAPPING:
		return "<START_MAPPING '" + e.Payload + "'>"
	case EMIT_KEY:
		return "<EMIT_KEY '" + e.Payload + "'>"
	case EMIT_VALUE:
		return "<EMIT_VALUE [" + payloadTypeToString(e.PayloadType) + "] '" + e.Payload + "'>"
	case END_MAPPING:
		return "<END_MAPPING '" + e.Payload + "'>"
	case START_ARRAY:
		return "<START_ARRAY '" + e.Payload + "'>"
	case EMIT_ELEMENT:
		return "<EMIT_ELEMENT '" + e.Payload + "'>"
	case END_ARRAY:
		return "<END_ARRAY '" + e.Payload + "'>"
	default:
		return "<UNKNOWN '" + e.Payload + "'>"
	}
}

func payloadTypeToString(pt PayLoadType) string {
	switch pt {
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case NULL:
		return "NULL"
	default:
		return "UNKNOWN"
	}
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
