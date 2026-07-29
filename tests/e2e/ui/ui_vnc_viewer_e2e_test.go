package ui_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestUiVncViewerE2e_CapabilitiesChromeSessionShowsVnc(t *testing.T) {
	cfg := config.MustLoad()
	navigateURL, err := cfg.ResolveUiBrowserURL()
	require.NoError(t, err)
	if !strings.HasSuffix(navigateURL, "/") {
		navigateURL += "/"
	}
	navigateHost, err := url.Parse(navigateURL)
	require.NoError(t, err)
	expectedHost := navigateHost.Hostname()

	allurex.Run(t, allurex.Meta{
		Name:      "Capabilities chrome session shows connected VNC and navigates remote browser",
		Package:   "tests.UiVncViewerE2eTests",
		Layer:     "e2e",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "VNC viewer",
		Story:     "VNC viewer in UI",
		Suite:     "UI VNC viewer e2e",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		runWithBrowser(t, func(page playwright.Page, baseURL string) {
			a.Step("Open dashboard and wait for CONNECTED", func() {
				openDashboard(t, page, baseURL)
			})
			a.Step("Open New Session page and select chrome", func() {
				_, err := page.Goto(baseURL+"/#/new-session", playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
				require.NoError(t, page.Locator("[data-testid=capabilities-setup]").WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
				chip := page.Locator("[data-testid=capabilities-browser-select] button[data-value='" + chromeChipValue(cfg) + "']")
				require.NoError(t, chip.WaitFor(playwright.LocatorWaitForOptions{
					State: playwright.WaitForSelectorStateVisible,
				}))
			})
			a.Step("Create chrome session from Capabilities", func() {
				var err error
				sessionID, err = hubapi.CreateSessionWithSelenoidOptions(
					cfg, cfg.Browser, cfg.ChromeVersionForSession(),
					map[string]any{
						"enableVNC":      true,
						"enableVideo":    false,
						"name":           "selenoid-ui-e2e",
						"sessionTimeout": "2m",
					},
				)
				require.NoError(t, err)
				_, err = page.Goto(baseURL+"/#/sessions/"+sessionID, playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateDomcontentloaded,
				})
				require.NoError(t, err)
			})
			a.Step("Navigate remote browser via hub WebDriver API", func() {
				require.NoError(t, hubapi.NavigateSession(cfg, sessionID, navigateURL))
			})
			a.Step("Verify remote page loaded", func() {
				currentURL, err := hubapi.GetSessionURL(cfg, sessionID)
				require.NoError(t, err)
				require.Contains(t, currentURL, expectedHost,
					"Remote browser URL should point to UI host")
				title, err := hubapi.GetSessionTitle(cfg, sessionID)
				require.NoError(t, err)
				require.NotEmpty(t, strings.TrimSpace(title),
					"Remote browser title should not be empty")
			})
			a.Step("Wait for VNC connected and unlock screen", func() {
				waitVncConnected(t, page)
			})
			a.Step("Read session id from URL", func() {
				fromURL := sessionIDFromURL(page.URL())
				require.Equal(t, sessionID, fromURL)
			})
		})
		if sessionID != "" {
			a.Step("Delete hub session", func() {
				require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			})
		}
	})
}
