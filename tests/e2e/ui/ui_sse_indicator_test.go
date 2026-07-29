package ui_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestUiSseIndicator_SseUsesConnectedStyling(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "SSE indicator uses connected StatusTile styling when connected",
		Package:   "tests.UiSseIndicatorTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status bar",
		Story:     "UI SSE indicator",
		Suite:     "UI SSE indicator",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Verify SSE indicator uses connected StatusTile styling", func() {
				assertStatusTileConnected(t, page, "#sse-status")
			})
		})
	})
}
