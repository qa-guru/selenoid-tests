package ui_test

import (
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestUiManualHarPlaywright_CapabilitiesCreateSessionShowsHarInArchive(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Manual Playwright Capabilities enableHar → Finished sessions HAR icon + HarViewer",
		Package:   "tests.UiManualHarTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "Capabilities manual session",
		Story:     "Manual Playwright session with hub HAR",
		Suite:     "UI manual HAR",
		Tags:      []string{"ui-manual-har", "playwright", "smoke", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open Capabilities and select Playwright chrome", func() {
				openDashboard(t, page, baseURL)
				openNewSession(t, page, baseURL)
				selectPlaywrightChrome(t, page, cfg)
			})
			a.Step("Enable enableHar + enableVideo and set session name", func() {
				setSegTrue(t, page, "caps-playwright-enable-har")
				setSegTrue(t, page, "caps-playwright-enable-video")
				setManualHarSessionName(t, page, "caps-playwright-session-name")
			})
			a.Step("Create Playwright session from Capabilities (UI bare WS path)", func() {
				sessionID = clickCreateSession(t, page)
				require.NotEmpty(t, sessionID)
				// Hub bootstraps about:blank and attaches CDP HAR for manual UI sessions.
				time.Sleep(3 * time.Second)
			})
			a.Step("Kill session from session panel", func() {
				killSessionFromUI(t, page)
			})
			a.Step("Wait until hub archives HAR artifact", func() {
				row, err := hubapi.WaitArchivedSessionHar(cfg, sessionID, 45*time.Second)
				require.NoError(t, err)
				require.NotNil(t, row)
				require.NotEmpty(t, row.HAR)
			})
			a.Step("Open Finished sessions archive", func() {
				openFinishedSessionsArchive(t, page, baseURL)
			})
			a.Step("Verify HAR artifact icon on the finished row", func() {
				waitArchiveHarIcon(t, page, baseURL, sessionID, 30*time.Second)
			})
			a.Step("Open session detail and verify HarViewer", func() {
				openArchiveSessionDetail(t, page, sessionID)
				require.NoError(t, page.Locator("[data-testid=session-har-viewer]").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
				count, err := page.Locator("[data-testid=session-no-har]").Count()
				require.NoError(t, err)
				require.Zero(t, count)
			})
		})

		a.Step("Hub /sessions/?json lists har artifact", func() {
			row, err := hubapi.WaitArchivedSessionHar(cfg, sessionID, 45*time.Second)
			require.NoError(t, err)
			require.NotNil(t, row)
			require.NotEmpty(t, row.HAR)
		})
		a.Step("Hub /har file is downloadable", func() {
			file, err := hubapi.WaitForSessionHar(cfg, sessionID, 45*time.Second)
			require.NoError(t, err)
			require.NotEmpty(t, file)
		})
	})
}
