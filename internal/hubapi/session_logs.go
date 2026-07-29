package hubapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

type wdLogEntry struct {
	Message string `json:"message"`
}

// GetPerformanceLogMessages POST /wd/hub/session/{id}/se/log type=performance.
func GetPerformanceLogMessages(cfg *config.Config, sessionID string) ([]string, error) {
	path := fmt.Sprintf("/wd/hub/session/%s/se/log", sessionID)
	resp, err := hubClient(cfg).PostJSON(path, map[string]string{"type": "performance"}, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Value []wdLogEntry `json:"value"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, err
	}
	msgs := make([]string, 0, len(payload.Value))
	for _, entry := range payload.Value {
		if entry.Message != "" {
			msgs = append(msgs, entry.Message)
		}
	}
	return msgs, nil
}

// ExecuteCdpCommand POST /wd/hub/session/{id}/goog/cdp/execute (HasCdp parity).
func ExecuteCdpCommand(cfg *config.Config, sessionID, cmd string, params map[string]any) (map[string]any, error) {
	path := fmt.Sprintf("/wd/hub/session/%s/goog/cdp/execute", sessionID)
	body := map[string]any{"cmd": cmd, "params": params}
	resp, err := hubClient(cfg).PostJSON(path, body, http.StatusOK)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Value map[string]any `json:"value"`
	}
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return nil, err
	}
	return payload.Value, nil
}
