package ui_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestUiStatusBar_SseAndSelenoidStayConnected(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "SSE and SELENOID stay CONNECTED on stable stack",
		Package:   "tests.UiStatusBarTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status bar",
		Story:     "UI status bar",
		Suite:     "UI status bar",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard and wait for CONNECTED", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Verify status indicators show connected StatusTile styling", func() {
				assertStatusTileConnected(t, page, "#sse-status")
				assertStatusTileConnected(t, page, "#selenoid-status")
			})
			a.Step("Keep CONNECTED stable for 5 seconds", func() {
				stayConnected(t, page, 5_000)
			})
		})
	})
}
