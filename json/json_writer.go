package json

import (
	"hbibel/yaml-to-json/common"
)

func RenderEvents(events <-chan common.Event) <-chan string {
	output := make(chan string)
	go func() {
		firstElement := true
		for op := range events {
			switch op.GetKind() {
			case common.START_MAPPING:
				firstElement = true
				output <- "{"
			case common.EMIT_KEY:
				if !firstElement {
					output <- ","
				}
				firstElement = false
				output <- renderAsKey(op)
				output <- ":"
			case common.EMIT_VALUE:
				output <- renderAsValue(op)
			case common.END_MAPPING:
				firstElement = false
				output <- "}"
			case common.START_ARRAY:
				firstElement = true
				output <- "["
			case common.EMIT_ELEMENT:
				if !firstElement {
					output <- ","
				}
				firstElement = false
			case common.END_ARRAY:
				firstElement = false
				output <- "]"
			}
		}
		close(output)
	}()
	return output
}

// TODO escape special characters
func renderAsKey(op common.Event) string {
	return "\"" + op.(common.HasPayload).GetPayload() + "\""
}

func renderAsValue(op common.Event) string {
	withPayload := op.(common.HasPayload)
	switch withPayload.GetPayLoadType() {
	case common.STRING:
		return "\"" + withPayload.GetPayload() + "\""
	case common.NUMBER:
		return string(withPayload.GetPayload())
	case common.BOOLEAN:
		return withPayload.GetPayload()
	case common.NULL:
		return "null"
	}
	panic("unknown payload type")
}
