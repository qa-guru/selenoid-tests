package hubapi

import (
	"encoding/json"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// Status is hub capacity JSON (flat /status or UI .state).
type Status struct {
	Total    int            `json:"total"`
	Used     int            `json:"used"`
	Queued   int            `json:"queued"`
	Pending  int            `json:"pending"`
	Browsers map[string]any `json:"browsers"`
	State    *Status        `json:"state"`
}

// Fetch GET hub status using config hubStatusPath (prod: /hub/status).
func Fetch(cfg *config.Config) (*Status, error) {
	client := httpx.New(cfg.APIBase())
	body, err := client.GetBytes(cfg.HubStatusPath)
	if err != nil {
		return nil, err
	}
	return Parse(body)
}

// Parse accepts flat hub JSON or UI-shaped {"state":{...}}.
func Parse(body []byte) (*Status, error) {
	var envelope struct {
		State *Status `json:"state"`
		Status
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	if envelope.State != nil {
		return envelope.State, nil
	}
	st := envelope.Status
	return &st, nil
}
