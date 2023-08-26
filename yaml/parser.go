package yaml

import (
	"hbibel/yaml-to-json/common"
	"strings"
)

// This parser is entirely hacked together and needs to be rewritten from
// scratch in order to implement the YAML spec properly. In particular, I
// completely ignored the "node style" stuff and implemented it only to get the
// most basic use cases working.

type nestingType int

const (
	IN_DOCUMENT nestingType = iota
	IN_ARRAY
	IN_MAPPING
)

type breadcrumb struct {
	nt       nestingType
	position uint32
}

func TokensToEvents(tokens <-chan Token) <-chan common.Event {
	events := make(chan common.Event)

	go func() {
		// to keep track of the nesting level within the document
		breadcrumbs := []breadcrumb{{IN_DOCUMENT, 0}}
		lineTokens := make([]Token, 0, 5)
		var lineEvents []common.Event

		for token := range tokens {
			if token.Kind() != NEWLINE {
				lineTokens = append(lineTokens, token)
				continue
			}

			lineEvents, breadcrumbs = toEvents(lineTokens, breadcrumbs)
			for _, event := range lineEvents {
				events <- event
			}

			lineTokens = lineTokens[:0]
		}
		lineEvents, breadcrumbs = toEvents(lineTokens, breadcrumbs)
		for _, event := range lineEvents {
			events <- event
		}

		// close remaining breadcrumb elements
		for _, breadcrumb := range breadcrumbs {
			switch breadcrumb.nt {
			case IN_ARRAY:
				events <- common.NewEndArrayEvent()
			case IN_MAPPING:
				events <- common.NewEndMappingEvent()
			}
		}

		close(events)
	}()

	return events
}

func isNumeric(s string) bool {
	hasDot := false
	for _, r := range s {
		if r == '.' {
			if hasDot {
				return false
			}
			hasDot = true
		} else if r > '9' || r < '0' {
			return false
		}
	}
	return true
}

func toEvents(tokens []Token, breadcrumbs []breadcrumb) ([]common.Event, []breadcrumb) {
	events := make([]common.Event, 0, 3)
	var position uint32 = 0

	for len(tokens) > 0 {
		if tokens[0].Kind() == INDENT {
			position = position + tokens[0].(*indentToken).spaceCount
			tokens = tokens[1:]
			continue
		}

		if len(tokens) == 1 {
			wt := tokens[0].(*wordToken)
			if wt.content == "true" {
				events = append(events, common.NewBooleanEvent("true"))
			} else if wt.content == "false" {
				events = append(events, common.NewBooleanEvent("false"))
			} else if wt.content == "null" {
				events = append(events, common.NewNullEvent())
			} else if isNumeric(wt.content) {
				events = append(events, common.NewNumberEvent(wt.content))
			} else {
				events = append(events, common.NewStringEvent(wt.content))
			}
			break
		}

		if len(tokens) > 1 && tokens[0].Kind() == DASH && tokens[1].Kind() == SPACE {
			notInArray := breadcrumbs[len(breadcrumbs)-1].nt != IN_ARRAY
			lessIndented := breadcrumbs[len(breadcrumbs)-1].position != position
			if notInArray || lessIndented {
				events = append(events, common.NewStartArrayEvent())
				breadcrumbs = append(breadcrumbs, breadcrumb{IN_ARRAY, position})
			}
			events = append(events, common.NewEmitElementEvent())
			position = position + 2
			tokens = tokens[2:]
			continue
		}

		// Try if it's a mapping key
		wb := strings.Builder{}
		for i, t := range tokens {
			if t.Kind() == COLON && (len(tokens) == i+1 || tokens[i+1].Kind() == SPACE) {
				events = append(events, common.NewStartMappingEvent())
				events = append(events, common.NewKeyEvent(wb.String()))
				breadcrumbs = append(breadcrumbs, breadcrumb{IN_MAPPING, position})
				position++
				if len(tokens) == i+1 {
					tokens = tokens[i+1:]
				} else {
					tokens = tokens[i+2:]
				}
				break
			}

			position = position + uint32(len(t.String()))
			wb.WriteString(t.String())
		}

		if len(tokens) == 0 {
			break
		}
		wb.Reset()

		// it's not an array, map, etc, so it's a word
		for _, t := range tokens {
			position = position + uint32(len(t.String()))
			wb.WriteString(t.String())
		}
		events = append(events, common.NewStringEvent(wb.String()))
		break
	}

	return events, breadcrumbs
}
