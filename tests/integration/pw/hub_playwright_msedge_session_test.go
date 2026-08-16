package pw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightMsedgeSession_RemoteSessionConnects(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright Microsoft Edge WS session connects via hub",
		Package:   "tests.integration.HubPlaywrightMsedgeSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright Microsoft Edge WS session",
		Story:     "Playwright Microsoft Edge WS session",
		Suite:     "Playwright Microsoft Edge hub WS session",
		Browser:   allurex.BrowserMsedge,
		Tags:      []string{"integration", "playwright"},
	}, func(a *allurex.A) {
		var endpoint string
		a.Step("Resolve Playwright Microsoft Edge WS endpoint", func() {
			var err error
			endpoint, err = cfg.ResolvePlaywrightWsEndpointForBrowser("playwright-msedge")
			require.NoError(t, err)
		})
		runPlaywrightSessionLifecycle(
			t, a, cfg, endpoint,
			"Connect Playwright Microsoft Edge via hub WS",
			"Close remote session",
		)
	})
}
