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

func TestUiHarContentControl_HiddenUntilEnableHar(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Capabilities harContent is hidden until enableHAR is on",
		Package:   "tests.UiHarContentTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "Capabilities HAR",
		Story:     "harContent control follows enableHAR",
		Suite:     "UI HAR content",
		Tags:      []string{"ui-har-content", "smoke", "positive"},
	}, func(a *allurex.A) {
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open Capabilities and select WebDriver chrome", func() {
				openDashboard(t, page, baseURL)
				openNewSession(t, page, baseURL)
				selectWebDriverChrome(t, page, cfg)
			})
			a.Step("harContent control absent while enableHAR is off", func() {
				require.NoError(t, page.Locator("[data-testid=caps-enable-har]").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
				n, err := page.Locator("[data-testid=caps-har-content]").Count()
				require.NoError(t, err)
				require.Zero(t, n)
			})
			a.Step("harContent control appears after enableHAR=true", func() {
				setSegTrue(t, page, "caps-enable-har")
				require.NoError(t, page.Locator("[data-testid=caps-har-content]").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
			})
		})
	})
}

func TestUiHarViewer_HubBodiesShowsResponseText(t *testing.T) {
	if tags := os.Getenv("TEST_TAGS"); strings.Contains(tags, "smoke") {
		t.Skip("excluded from TEST_TAGS=smoke gate; run via hub-prod / ui without smoke tags")
	}
	cfg := config.MustLoad()
	skipHeavyHarOnGithubCI(t, cfg)
	targetURL := strings.TrimRight(cfg.SmokeURL, "/")

	allurex.Run(t, allurex.Meta{
		Name:      "HarViewer Response tab shows body text for hub harContent=bodies",
		Package:   "tests.UiHarContentTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "HarViewer",
		Story:     "Finished bodies HAR shows content.text",
		Suite:     "UI HAR content",
		Tags:      []string{"ui-har-content", "webdriver", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub enableHAR bodies session and archive HAR", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithSelenoidOptions(cfg, "chrome", cfg.ChromeVersionForSession(), map[string]any{
				"enableVNC":      false,
				"enableVideo":    false,
				"enableHAR":      true,
				"harContent":     "bodies",
				"enableLog":      false,
				"sessionTimeout": "2m",
				"name":           "ui-har-bodies-e2e",
			})
			require.NoError(t, err)
			require.NoError(t, hubapi.NavigateSession(cfg, sessionID, targetURL))
			time.Sleep(2 * time.Second)
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			waitHubArchivedHar(t, cfg, sessionID)
		})

		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open finished session HarViewer", func() {
				_, err := page.Goto(baseURL+"/#/sessions/"+sessionID, playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
				require.NoError(t, page.Locator("[data-testid=session-har-viewer]").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(60_000),
				}))
				require.NoError(t, page.Locator("[data-testid=session-har-row-0]").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(60_000),
				}))
			})
			a.Step("Response tab shows captured body, not muted placeholder", func() {
				require.NoError(t, page.Locator("[data-testid=session-har-row-0]").Click())
				require.NoError(t, page.Locator("[data-testid=session-har-tab-response]").Click())
				panel := page.Locator("[data-testid=session-har-panel-response]")
				require.NoError(t, panel.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
				body, err := panel.InnerText()
				require.NoError(t, err)
				lower := strings.ToLower(body)
				require.NotContains(t, body, "Body not captured")
				hasBody := strings.Contains(body, config.DefaultSmokeHeading) || strings.Contains(body, config.DefaultSmokeTitle) || strings.Contains(lower, "<!doctype html>")
				require.True(t, hasBody, "expected captured HTML in HarViewer Response tab, got %q", trimForErr(body))
				sizeCell, err := page.Locator("[data-testid=session-har-row-0] td").Nth(4).InnerText()
				require.NoError(t, err)
				size := strings.TrimSpace(sizeCell)
				require.NotContains(t, []string{"", "—", "-"}, size, "table Size must not be empty/dash")
			})
		})
	})
}

func trimForErr(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 400 {
		return s[:400]
	}
	return s
}
