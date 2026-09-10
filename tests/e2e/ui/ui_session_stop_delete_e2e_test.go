package ui_test

import (
	"strings"
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestUiSessionStop_ShowsFinishedImmediatelyThenDelete(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Stop session shows FINISHED immediately, then Delete wipes artifacts",
		Package:   "tests.UiSessionStopDeleteTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "Session stop",
		Story:     "Optimistic Stop then Delete",
		Suite:     "UI session stop delete",
		Tags:      []string{"smoke", "positive", "ui-session-stop"},
	}, func(a *allurex.A) {
		var sessionID string
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			console := attachConsoleErrorTracker(page)

			a.Step("Open dashboard", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Create WebDriver session with VNC+video via hub API", func() {
				var err error
				sessionID, err = hubapi.CreateSessionWithSelenoidOptions(
					cfg, cfg.Browser, cfg.ChromeVersionForSession(),
					map[string]any{
						"enableVNC":      true,
						"enableVideo":    true,
						"name":           "selenoid-ui-stop-e2e",
						"sessionTimeout": "2m",
					},
				)
				require.NoError(t, err)
				require.NotEmpty(t, sessionID)
			})
			a.Step("Open session page", func() {
				_, err := page.Goto(baseURL+"/#/sessions/"+sessionID, playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
			})
			if !strings.Contains(cfg.Env, "github") {
				a.Step("Wait for VNC connected", func() {
					waitVncConnected(t, page)
				})
			}
			a.Step("Stop — FINISHED and VNC gone without waiting for SSE", func() {
				killedID := killSessionFromUI(t, page)
				require.Equal(t, sessionID, killedID)
				console.assertNoStopTeardownErrors(t)
			})
			a.Step("Video appears and Delete becomes enabled", func() {
				require.NoError(t, page.Locator("[data-testid=session-detail-video]").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(float64(createSessionWait(cfg).Milliseconds())),
				}))
				enabled, err := page.Locator("[data-testid=session-delete]").IsEnabled()
				require.NoError(t, err)
				require.True(t, enabled, "Delete session must enable once artifacts exist and the page is no longer live")
			})
			a.Step("Delete session from UI", func() {
				deleteSessionFromUI(t, page)
				sessionID = ""
			})
		})
		if sessionID != "" {
			a.Step("Cleanup hub session", func() {
				_ = hubapi.DeleteSession(cfg, sessionID)
			})
		}
	})
}
