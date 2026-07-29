package hubapi

import "encoding/json"

// SessionCreateResponse is a W3C New Session success body (value.sessionId).
type SessionCreateResponse struct {
	Value struct {
		SessionID string `json:"sessionId"`
	} `json:"value"`
}

// ParseSessionID extracts value.sessionId from a create-session JSON body.
func ParseSessionID(body []byte) (string, error) {
	var resp SessionCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.Value.SessionID, nil
}
