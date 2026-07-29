package pw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightWebkitSession_RemoteSessionConnects(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright WebKit WS session connects via hub",
		Package:   "tests.integration.HubPlaywrightWebkitSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WebKit WS session",
		Story:     "Playwright WebKit WS session",
		Suite:     "Playwright WebKit hub WS session",
		Tags:      []string{"integration", "playwright"},
	}, func(a *allurex.A) {
		var endpoint string
		a.Step("Resolve Playwright WebKit WS endpoint", func() {
			var err error
			endpoint, err = cfg.ResolvePlaywrightWsEndpointForBrowser("playwright-webkit")
			require.NoError(t, err)
		})
		runPlaywrightSessionLifecycle(
			t, a, cfg, endpoint,
			"Connect Playwright WebKit via hub WS",
			"Close remote session",
		)
	})
}
