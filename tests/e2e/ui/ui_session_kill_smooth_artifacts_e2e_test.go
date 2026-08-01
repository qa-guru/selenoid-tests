package ui_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestUiSessionKillSmoothArtifacts_WebDriverKillKeepsLayoutAndUpdatesInPlace(t *testing.T) {
	// Heavy VNC+video+HAR Create Session is too slow/flaky for the release smoke gate.
	if tags := os.Getenv("TEST_TAGS"); strings.Contains(tags, "smoke") {
		t.Skip("excluded from TEST_TAGS=smoke gate; run via hub-prod / ui without smoke tags")
	}
	cfg := config.MustLoad()
	targetURL := strings.TrimRight(cfg.SmokeURL, "/") + "/"
	allurex.Run(t, allurex.Meta{
		Name:      "Kill session in-place: URL stable, VNC→video + HAR swap without layout jump",
		Package:   "tests.UiSessionKillSmoothTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "Session kill",
		Story:     "Smooth artifact transition after kill",
		Suite:     "UI session kill smooth",
		Tags:      []string{"ui-session-kill", "webdriver", "smoke", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			console := attachConsoleErrorTracker(page)

			a.Step("Open Capabilities and select WebDriver chrome with HAR + video + VNC", func() {
				openDashboard(t, page, baseURL)
				openNewSession(t, page, baseURL)
				selectWebDriverChrome(t, page, cfg)
				fillHubAuthFromConfig(t, page, cfg)
				setSegTrue(t, page, "caps-enable-har")
				setSegTrue(t, page, "caps-enable-video")
				setSegTrue(t, page, "caps-enable-vnc")
				setManualHarSessionName(t, page, "caps-session-name")
			})
			a.Step("Create session and wait for VNC", func() {
				sessionID = clickCreateSession(t, page)
				require.NotEmpty(t, sessionID)
				waitVncConnected(t, page)
			})
			a.Step("Navigate remote browser to smoke URL", func() {
				require.NoError(t, hubapi.NavigateSession(cfg, sessionID, targetURL))
				time.Sleep(2 * time.Second)
			})

			var mediaBefore, logBefore, harBefore layoutRect
			a.Step("Snapshot layout rects before kill", func() {
				mediaBefore = captureLayoutRect(t, page, "[data-testid=session-media-slot]")
				logBefore = captureLayoutRect(t, page, "[data-testid=session-log-panel]")
				harBefore = captureLayoutRect(t, page, "[data-testid=session-har-viewer]")
				count, err := page.Locator("[data-testid=session-har-viewer]").Count()
				require.NoError(t, err)
				require.Equal(t, 1, count, "expected exactly one HarViewer while live")
			})

			urlBefore := page.URL()
			a.Step("Kill session — URL must stay on /#/sessions/<id>", func() {
				killedID := killSessionFromUI(t, page)
				require.Equal(t, sessionID, killedID)
				require.Equal(t, urlBefore, page.URL())
				require.Equal(t, sessionID, sessionIDFromURL(page.URL()))
			})

			a.Step("Layout slots stable within 2px; single HarViewer; log panel unchanged", func() {
				assertLayoutStable(t, "session-media-slot", mediaBefore, captureLayoutRect(t, page, "[data-testid=session-media-slot]"), layoutTolerancePx)
				assertLayoutStable(t, "session-log-panel", logBefore, captureLayoutRect(t, page, "[data-testid=session-log-panel]"), layoutTolerancePx)
				assertLayoutStable(t, "session-har-viewer", harBefore, captureLayoutRect(t, page, "[data-testid=session-har-viewer]"), layoutTolerancePx)

				count, err := page.Locator("[data-testid=session-har-viewer]").Count()
				require.NoError(t, err)
				require.Equal(t, 1, count)

				require.NoError(t, page.Locator("text=FINISHED").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
			})

			a.Step("VNC slot transitions to video in-place", func() {
				require.NoError(t, page.Locator("[data-testid=session-detail-video]").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(float64(createSessionTimeout.Milliseconds())),
				}))
				assertLayoutStable(t, "session-media-slot-after-video",
					mediaBefore, captureLayoutRect(t, page, "[data-testid=session-media-slot]"), layoutTolerancePx)
			})

			a.Step("Console has no errors", func() {
				console.assertEmpty(t)
			})
		})

		if sessionID != "" {
			a.Step("Hub archived HAR artifact", func() {
				waitHubArchivedHar(t, cfg, sessionID)
			})
		}
	})
}
