package hubapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

type wdStringValue struct {
	Value string `json:"value"`
}

// NavigateSession POST /wd/hub/session/{id}/url (HubSessionApi.navigate parity).
func NavigateSession(cfg *config.Config, sessionID, targetURL string) error {
	body := map[string]string{"url": targetURL}
	_, err := hubClient(cfg).PostJSON(
		fmt.Sprintf("/wd/hub/session/%s/url", sessionID),
		body,
		http.StatusOK,
	)
	return err
}

// GetSessionTitle GET /wd/hub/session/{id}/title.
func GetSessionTitle(cfg *config.Config, sessionID string) (string, error) {
	path := fmt.Sprintf("/wd/hub/session/%s/title", sessionID)
	resp, err := hubClient(cfg).GetExpectStatus(path, http.StatusOK)
	if err != nil {
		return "", err
	}
	var payload wdStringValue
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return "", err
	}
	return payload.Value, nil
}

// GetSessionURL GET /wd/hub/session/{id}/url.
func GetSessionURL(cfg *config.Config, sessionID string) (string, error) {
	path := fmt.Sprintf("/wd/hub/session/%s/url", sessionID)
	resp, err := hubClient(cfg).GetExpectStatus(path, http.StatusOK)
	if err != nil {
		return "", err
	}
	var payload wdStringValue
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return "", err
	}
	return payload.Value, nil
}
