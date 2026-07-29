package uiapi

import (
	"encoding/json"

	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

// SseEvent is a UI SSE hub event payload (state and/or errors).
type SseEvent struct {
	State  *hubapi.Status `json:"state"`
	Errors []string       `json:"errors"`
}

// HasState reports whether the event carries hub state.
func (e *SseEvent) HasState() bool {
	return e != nil && e.State != nil
}

// HasErrors reports whether the event carries a non-empty errors list.
func (e *SseEvent) HasErrors() bool {
	return e != nil && len(e.Errors) > 0
}

// ParseSseEvent unmarshals an SSE JSON data payload.
func ParseSseEvent(body []byte) (*SseEvent, error) {
	var ev SseEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}
