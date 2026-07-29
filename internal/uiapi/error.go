package uiapi

import "encoding/json"

// ErrorMessage is one UI /error entry.
type ErrorMessage struct {
	Msg string `json:"msg"`
}

// ErrorResponse is UI GET /error payload.
type ErrorResponse struct {
	Errors  []ErrorMessage `json:"errors"`
	Version string         `json:"version"`
}

// ParseError unmarshals a UI error JSON body.
func ParseError(body []byte) (*ErrorResponse, error) {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}
	return &errResp, nil
}
