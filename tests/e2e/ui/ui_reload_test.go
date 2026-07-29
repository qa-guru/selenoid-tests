package ui_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestUiReload_DashboardReloadKeepsConnected(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "Dashboard reload keeps CONNECTED indicators",
		Package:   "tests.UiReloadTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI dashboard",
		Story:     "UI dashboard reload",
		Suite:     "UI dashboard reload",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard and wait for CONNECTED", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Reload dashboard", func() {
				reloadDashboard(t, page)
				waitConnected(t, page)
			})
		})
	})
}
