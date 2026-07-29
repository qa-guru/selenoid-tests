package hubapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// SessionCreateResult holds session id and echoed browserName.
type SessionCreateResult struct {
	SessionID   string
	BrowserName string
}

// CreateSession POST /wd/hub/session with config browser/version.
func CreateSession(cfg *config.Config) (string, error) {
	return CreateSessionWithBrowser(cfg, cfg.Browser, cfg.BrowserVersion)
}

// CreateSessionWithBrowser POST /wd/hub/session for explicit browser/version.
func CreateSessionWithBrowser(cfg *config.Config, browserName, browserVersion string) (string, error) {
	resp, err := hubClient(cfg).PostJSON("/wd/hub/session", CreateSessionBody(browserName, browserVersion), http.StatusOK)
	if err != nil {
		return "", err
	}
	return ParseSessionID(resp.Body)
}

// CreateSessionWithSelenoidOptions POST session with selenoid:options map.
func CreateSessionWithSelenoidOptions(cfg *config.Config, browserName, browserVersion string, selenoidOptions map[string]any) (string, error) {
	body := createSessionBodyWithOptions(browserName, browserVersion, selenoidOptions)
	resp, err := hubClient(cfg).PostJSON("/wd/hub/session", body, http.StatusOK)
	if err != nil {
		return "", err
	}
	return ParseSessionID(resp.Body)
}

// CreateSessionWithCapabilities POST session and read browserName from capabilities.
func CreateSessionWithCapabilities(cfg *config.Config) (*SessionCreateResult, error) {
	resp, err := hubClient(cfg).PostJSON("/wd/hub/session", CreateSessionBody(cfg.Browser, cfg.BrowserVersion), http.StatusOK)
	if err != nil {
		return nil, err
	}
	sessionID, err := ParseSessionID(resp.Body)
	if err != nil {
		return nil, err
	}
	browserName, err := parseBrowserName(resp.Body)
	if err != nil {
		return nil, err
	}
	return &SessionCreateResult{SessionID: sessionID, BrowserName: browserName}, nil
}

// CreateSessionExpectStatus POST session and require HTTP status (negative tests).
func CreateSessionExpectStatus(cfg *config.Config, browserName, browserVersion string, expected int) error {
	_, err := hubClient(cfg).PostJSON("/wd/hub/session", CreateSessionBody(browserName, browserVersion), expected)
	return err
}

// DeleteSession DELETE /wd/hub/session/{id} (expects 200).
func DeleteSession(cfg *config.Config, sessionID string) error {
	_, err := hubClient(cfg).Delete("/wd/hub/session/"+sessionID, http.StatusOK)
	return err
}

// DeleteSessionExpectStatus DELETE and require HTTP status.
func DeleteSessionExpectStatus(cfg *config.Config, sessionID string, expected int) error {
	_, err := hubClient(cfg).Delete("/wd/hub/session/"+sessionID, expected)
	return err
}

func createSessionBodyWithOptions(browserName, browserVersion string, selenoidOptions map[string]any) map[string]any {
	alwaysMatch := createAlwaysMatch(browserName, browserVersion)
	if len(selenoidOptions) > 0 {
		alwaysMatch["selenoid:options"] = selenoidOptions
	}
	return map[string]any{
		"capabilities": map[string]any{
			"alwaysMatch": alwaysMatch,
		},
	}
}

func parseBrowserName(body []byte) (string, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	value, _ := payload["value"].(map[string]any)
	if value == nil {
		return "", fmt.Errorf("missing value in session response")
	}
	caps, _ := value["capabilities"].(map[string]any)
	if caps == nil {
		return "", nil
	}
	if name, ok := caps["browserName"].(string); ok && name != "" {
		return name, nil
	}
	if always, ok := caps["alwaysMatch"].(map[string]any); ok {
		if name, ok := always["browserName"].(string); ok {
			return name, nil
		}
	}
	return "", nil
}
