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

func TestUiSessionsList_ActiveSessionAppearsInList(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Active hub session appears in dashboard sessions list",
		Package:   "tests.UiSessionsListTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI sessions list",
		Story:     "UI sessions list",
		Suite:     "UI sessions list",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session via API", func() {
			var err error
			// Java uses config.browserVersion (149.0-min on github_e2e); warm chrome when min image unavailable.
			sessionID, err = hubapi.CreateSessionWithBrowser(cfg, cfg.Browser, cfg.ChromeVersionForSession())
			require.NoError(t, err)
		})
		defer func() {
			a.Step("Delete hub session", func() {
				require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			})
		}()

		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard and wait for stack Connected", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Open Sessions page", func() {
				_, err := page.Goto(baseURL+"/#/sessions", playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
			})
			a.Step("Verify session row shows browser name", func() {
				name := page.Locator(".sessions__list .session .browser .name").First()
				require.NoError(t, name.WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(float64(sessionAppearTimeout.Milliseconds())),
				}))
				text, err := name.InnerText()
				require.NoError(t, err)
				require.True(t, strings.Contains(strings.ToLower(text), cfg.Browser),
					"expected browser %q in session row, got %q", cfg.Browser, text)
			})
			a.Step("Verify live-sessions empty-state is hidden", func() {
				empty := page.Locator(".sessions-panel .no-any")
				visible, err := empty.IsVisible()
				require.NoError(t, err)
				require.False(t, visible, "empty state should be hidden when session is live")
			})
		})
	})
}
