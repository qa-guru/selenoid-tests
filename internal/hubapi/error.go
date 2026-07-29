package hubapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// ErrorValue is hub GET /error JSON value block.
type ErrorValue struct {
	Error string `json:"error"`
}

// ErrorResponse is hub GET /error body.
type ErrorResponse struct {
	Value ErrorValue `json:"value"`
}

// FetchErrorExpectInvalidSession GET /error and verify invalid session id JSON (404).
func FetchErrorExpectInvalidSession(cfg *config.Config) error {
	resp, err := hubClient(cfg).GetExpectStatus("/error", http.StatusNotFound)
	if err != nil {
		return err
	}
	var payload ErrorResponse
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return err
	}
	if payload.Value.Error != "invalid session id" {
		return fmt.Errorf("expected value.error=invalid session id, got %q", payload.Value.Error)
	}
	return nil
}
