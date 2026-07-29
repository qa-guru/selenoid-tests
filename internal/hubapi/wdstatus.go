package hubapi

import "encoding/json"

// WebDriverStatus is W3C GET /status payload (HubWebDriverStatus).
type WebDriverStatus struct {
	Value WebDriverStatusValue `json:"value"`
}

// WebDriverStatusValue is the inner ready/message object.
type WebDriverStatusValue struct {
	Ready   bool   `json:"ready"`
	Message string `json:"message"`
}

// ParseWebDriverStatus unmarshals a WebDriver status JSON body.
func ParseWebDriverStatus(body []byte) (*WebDriverStatus, error) {
	var st WebDriverStatus
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, err
	}
	return &st, nil
}
