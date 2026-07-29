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

type wdElementRef struct {
	Value map[string]string `json:"value"`
}

// FindElement POST /wd/hub/session/{id}/element (W3C element reference).
func FindElement(cfg *config.Config, sessionID, using, value string) (string, error) {
	body := map[string]string{"using": using, "value": value}
	resp, err := hubClient(cfg).PostJSON(
		fmt.Sprintf("/wd/hub/session/%s/element", sessionID),
		body,
		http.StatusOK,
	)
	if err != nil {
		return "", err
	}
	return parseElementID(resp.Body)
}

// GetElementText GET /wd/hub/session/{id}/element/{elementId}/text.
func GetElementText(cfg *config.Config, sessionID, elementID string) (string, error) {
	path := fmt.Sprintf("/wd/hub/session/%s/element/%s/text", sessionID, elementID)
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

// GetElementTextBySelector finds an element and reads its visible text.
func GetElementTextBySelector(cfg *config.Config, sessionID, cssSelector string) (string, error) {
	elementID, err := FindElement(cfg, sessionID, "css selector", cssSelector)
	if err != nil {
		return "", err
	}
	return GetElementText(cfg, sessionID, elementID)
}

func parseElementID(body []byte) (string, error) {
	var payload wdElementRef
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	for _, id := range payload.Value {
		if id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("no element id in session element response")
}
