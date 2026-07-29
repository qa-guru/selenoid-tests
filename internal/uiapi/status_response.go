package uiapi

import (
	"encoding/json"

	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

// StatusResponse is the UI /status wrapper with proxied hub .state.
type StatusResponse struct {
	State *hubapi.Status `json:"state"`
}

// ParseStatusResponse unmarshals UI-shaped /status JSON (keeps the wrapper).
func ParseStatusResponse(body []byte) (*StatusResponse, error) {
	var resp StatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
