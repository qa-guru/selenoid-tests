package pw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightFirefoxSession_RemoteSessionConnects(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright Firefox WS session connects via hub",
		Package:   "tests.integration.HubPlaywrightFirefoxSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright Firefox WS session",
		Story:     "Playwright Firefox WS session",
		Suite:     "Playwright Firefox hub WS session",
		Browser:   allurex.BrowserFirefox,
		Tags:      []string{"integration", "playwright"},
	}, func(a *allurex.A) {
		var endpoint string
		a.Step("Resolve Playwright Firefox WS endpoint", func() {
			var err error
			endpoint, err = cfg.ResolvePlaywrightWsEndpointForBrowser("playwright-firefox")
			require.NoError(t, err)
		})
		runPlaywrightSessionLifecycle(
			t, a, cfg, endpoint,
			"Connect Playwright Firefox via hub WS",
			"Close remote session",
		)
	})
}
