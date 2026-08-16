package pw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightChromeSession_RemoteSessionConnects(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright Chrome WS session connects via hub",
		Package:   "tests.integration.HubPlaywrightChromeSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright Chrome WS session",
		Story:     "Playwright Chrome WS session",
		Suite:     "Playwright Chrome hub WS session",
		Browser:   allurex.BrowserChrome,
		Tags:      []string{"integration", "playwright"},
	}, func(a *allurex.A) {
		var endpoint string
		a.Step("Resolve Playwright Chrome WS endpoint", func() {
			var err error
			endpoint, err = cfg.ResolvePlaywrightWsEndpointForBrowser("playwright-chrome")
			require.NoError(t, err)
		})
		runPlaywrightSessionLifecycle(
			t, a, cfg, endpoint,
			"Connect Playwright Chrome via hub WS",
			"Close remote session",
		)
	})
}
