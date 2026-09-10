package playwrightapi

import (
	"net/http"
	"strings"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// AppendQuery adds query parameters to a Playwright hub WS URL.
func AppendQuery(wsEndpoint, query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return wsEndpoint
	}
	sep := "?"
	if strings.Contains(wsEndpoint, "?") {
		sep = "&"
	}
	return wsEndpoint + sep + query
}

// HTTPURL rewrites a Playwright hub WS endpoint to http(s) for handshake probes.
func HTTPURL(wsEndpoint string) string {
	s := strings.TrimSpace(wsEndpoint)
	switch {
	case strings.HasPrefix(s, "wss://"):
		return "https://" + strings.TrimPrefix(s, "wss://")
	case strings.HasPrefix(s, "ws://"):
		return "http://" + strings.TrimPrefix(s, "ws://")
	default:
		return s
	}
}

// Handshake GETs the Playwright hub URL with WebSocket upgrade headers.
// Hub v3.0.16 answers 400 on *-min + enableHAR before upgrading the socket.
func Handshake(cfg *config.Config, wsEndpoint string) (*httpx.Response, error) {
	if strings.TrimSpace(wsEndpoint) == "" {
		var err error
		wsEndpoint, err = cfg.ResolvePlaywrightWsEndpoint()
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{
		"Connection":            "Upgrade",
		"Upgrade":               "websocket",
		"Sec-WebSocket-Version": "13",
		"Sec-WebSocket-Key":     "dGhlIHNhbXBsZSBub25jZQ==",
	}
	return httpx.New("").Do(http.MethodGet, HTTPURL(wsEndpoint), nil, headers)
}

// Connect opens a remote Playwright browser via hub WS (default chromium endpoint from config).
// firefox/webkit use their engine; chrome/msedge/chromium share Chromium.Connect.
func Connect(pw *playwright.Playwright, cfg *config.Config, wsEndpoint string) (playwright.Browser, error) {
	if strings.TrimSpace(wsEndpoint) == "" {
		var err error
		wsEndpoint, err = cfg.ResolvePlaywrightWsEndpoint()
		if err != nil {
			return nil, err
		}
	}
	if strings.Contains(wsEndpoint, "firefox") {
		return pw.Firefox.Connect(wsEndpoint)
	}
	if strings.Contains(wsEndpoint, "webkit") {
		return pw.WebKit.Connect(wsEndpoint)
	}
	return pw.Chromium.Connect(wsEndpoint)
}

// Close shuts down a remote Playwright browser session.
func Close(browser playwright.Browser) error {
	if browser == nil {
		return nil
	}
	return browser.Close()
}
