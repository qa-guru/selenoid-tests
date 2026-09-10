package min_test

import (
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestHubPlaywrightMinHarSession_HubDeleteReturnsWithinSLA(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Hub DELETE of Playwright-min enableHAR session returns within 8s",
		Package:   "tests.integration.HubPlaywrightMinHarDeleteSlaTests",
		Layer:     "integration",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Playwright min session delete",
		Story:     "min image has no DevTools :7070 — DELETE must not wait 35s for CDP",
		Suite:     "Playwright hub min HAR delete SLA",
		Browser:   allurex.BrowserChromium,
		Tags:      []string{"integration", "min", "positive"},
	}, func(a *allurex.A) {
		pw, err := playwright.Run()
		require.NoError(t, err)
		defer func() { require.NoError(t, pw.Stop()) }()

		wsEndpoint, err := cfg.ResolvePlaywrightWsEndpoint()
		require.NoError(t, err)
		ws := playwrightapi.AppendQuery(wsEndpoint, "enableHAR=true&enableVideo=false")

		var before map[string]struct{}
		a.Step("Snapshot hub session ids", func() {
			before, err = hubapi.CollectSessionIDs(cfg)
			require.NoError(t, err)
		})

		var browser playwright.Browser
		a.Step("Connect Playwright-min with enableHAR", func() {
			browser, err = playwrightapi.Connect(pw, cfg, ws)
			require.NoError(t, err)
			require.True(t, browser.IsConnected())
		})
		defer func() { _ = playwrightapi.Close(browser) }()

		var sessionID string
		a.Step("Read new hub session id", func() {
			sessionID, err = hubapi.WaitNewSessionID(cfg, before, 30*time.Second)
			require.NoError(t, err)
			require.NotEmpty(t, sessionID, "hub status must expose the Playwright session id")
		})

		a.Step("Hub DELETE within HAR SLA (must not block 35s on missing CDP)", func() {
			require.NoError(t, hubapi.DeleteSessionWithin(cfg, sessionID, hubapi.PlaywrightHarSessionDeleteSLA))
		})
	})
}
