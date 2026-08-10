package ui_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestUiVideosPage_ArchiveLoadsFinishedSessionsShell(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "Sessions archive loads Finished sessions shell without requiring full catalog",
		Package:   "tests.UiVideosPageTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI sessions archive",
		Story:     "UI finished sessions archive",
		Suite:     "UI sessions archive page",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open Sessions page (archive)", func() {
				_, err := page.Goto(baseURL+"/#/sessions", playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
				require.NoError(t, page.Locator("[data-testid='archive-panel']").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
			})
			a.Step("Verify Finished sessions table shell", func() {
				require.NoError(t, page.Locator("[data-testid='archive-panel']").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
				count, err := page.Locator("[data-testid='archive-table']").Count()
				require.NoError(t, err)
				require.GreaterOrEqual(t, count, 1)
				headCount, err := page.Locator("[data-testid='archive-head']").Count()
				require.NoError(t, err)
				require.GreaterOrEqual(t, headCount, 1)
			})
			a.Step("Verify empty state, rows, or pager", func() {
				pager := page.Locator("[data-testid='archive-pager']")
				if count, err := pager.Count(); err == nil && count > 0 {
					if visible, _ := pager.IsVisible(); visible {
						return
					}
				}
				empty := page.Locator(".archive-panel .no-any")
				if count, err := empty.Count(); err == nil && count > 0 {
					return
				}
				if visible, _ := empty.IsVisible(); visible {
					return
				}
				rows, err := page.Locator("[data-testid='session-card']").Count()
				require.NoError(t, err)
				if rows > 0 {
					return
				}
				// Empty archive on prod: table shell stays in DOM.
				tableCount, err := page.Locator("[data-testid='archive-table']").Count()
				require.NoError(t, err)
				require.GreaterOrEqual(t, tableCount, 1)
			})
		})
	})
}
