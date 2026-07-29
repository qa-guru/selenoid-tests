package ui_test

import (
	"strings"
	"testing"

	"github.com/mxschmitt/playwright-go"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestUiDashboardLoad_DashboardOpensRootUrl(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "Dashboard opens root URL",
		Package:   "tests.UiDashboardLoadTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI dashboard",
		Story:     "UI dashboard load",
		Suite:     "UI dashboard load",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Verify Statistics route (or root) is loaded", func() {
				current := page.URL()
				ok := strings.Contains(current, "/#/statistics") ||
					strings.HasSuffix(current, "/") ||
					strings.HasSuffix(current, ":8080/")
				if !ok {
					t.Fatalf("Expected Statistics hash route or root, got: %s", current)
				}
			})
		})
	})
}
