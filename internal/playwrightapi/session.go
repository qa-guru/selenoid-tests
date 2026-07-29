package playwrightapi

import (
	"strings"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// Connect opens a remote Playwright browser via hub WS (default chromium endpoint from config).
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
