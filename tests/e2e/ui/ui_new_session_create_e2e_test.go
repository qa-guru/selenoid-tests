package ui_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestUiNewSession_CreateSessionOpensDetail(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "New Session Create Session opens the session page",
		Package:   "tests.UiNewSessionCreateTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "New Session",
		Story:     "Create Session button",
		Suite:     "UI new session create",
		Tags:      []string{"smoke", "positive", "ui-new-session"},
	}, func(a *allurex.A) {
		var sessionID string
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open New Session and select WebDriver chrome", func() {
				openDashboard(t, page, baseURL)
				openNewSession(t, page, baseURL)
				selectWebDriverChrome(t, page, cfg)
				fillHubAuthFromConfig(t, page, cfg)
			})
			a.Step("Click Create Session — session page, not a skip", func() {
				sessionID = clickCreateSession(t, page)
				require.NotEmpty(t, sessionID)
				require.Contains(t, page.URL(), "/#/sessions/"+sessionID)
				require.NoError(t, page.Locator("[data-testid=session-stop]").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(float64(createSessionWait(cfg).Milliseconds())),
				}), "Create Session must land on a live session with Stop, not stay on New Session")
			})
		})
		if sessionID != "" {
			a.Step("Cleanup hub session", func() {
				_ = hubapi.DeleteSession(cfg, sessionID)
			})
		}
	})
}
