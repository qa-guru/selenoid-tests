package hubapi

import (
	"fmt"
	"net/http"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/helpers"
)

// CreateHarCaptureSession opens warm Chrome with performance logging and enableHAR=false.
func CreateHarCaptureSession(cfg *config.Config, label string) (string, error) {
	alwaysMatch := map[string]any{
		"browserName":    "chrome",
		"browserVersion": cfg.ChromeVersionForSession(),
		"goog:chromeOptions": map[string]any{
			"args": []string{
				"--headless=new",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--disable-features=ServiceWorker",
			},
		},
		"goog:loggingPrefs": map[string]any{
			"performance": "ALL",
		},
		"selenoid:options": map[string]any{
			"enableVNC":   false,
			"enableVideo": false,
			"enableHAR":   false,
			"headless":    true,
			"name":        label,
		},
	}
	body := map[string]any{
		"capabilities": map[string]any{
			"alwaysMatch": alwaysMatch,
		},
	}
	resp, err := hubClient(cfg).PostJSON("/wd/hub/session", body, http.StatusOK)
	if err != nil {
		return "", err
	}
	return ParseSessionID(resp.Body)
}

// CollectHarFromSession builds HAR JSON from performance logs (+ CDP bodies when mode=BODIES).
func CollectHarFromSession(cfg *config.Config, sessionID string, mode helpers.HarContentMode) ([]byte, error) {
	msgs, err := GetPerformanceLogMessages(cfg, sessionID)
	if err != nil {
		return nil, err
	}
	bodies := map[string]helpers.CapturedBody{}
	if mode == helpers.HarBodies {
		for _, requestID := range helpers.FinishedRequestIds(msgs) {
			result, err := ExecuteCdpCommand(cfg, sessionID, "Network.getResponseBody", map[string]any{
				"requestId": requestID,
			})
			if err != nil || result == nil {
				continue
			}
			text := stringVal(result["body"])
			if text == "" {
				continue
			}
			base64 := boolVal(result["base64Encoded"])
			bodies[requestID] = helpers.CapturedBody{Text: text, Base64Encoded: base64}
		}
	}
	har := helpers.ToHar(msgs, mode, bodies)
	if har == "" {
		return nil, fmt.Errorf("HarCapture produced empty HAR")
	}
	return []byte(har), nil
}

func stringVal(o any) string {
	if s, ok := o.(string); ok {
		return s
	}
	return ""
}

func boolVal(o any) bool {
	switch v := o.(type) {
	case bool:
		return v
	default:
		return false
	}
}
