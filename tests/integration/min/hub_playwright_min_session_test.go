package min_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightMinSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright min WS session starts browser node and completes",
		Package:   "tests.integration.HubPlaywrightMinSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS session (min)",
		Story:     "Playwright WS session (min)",
		Suite:     "Playwright hub WS session (min)",
		Browser:   allurex.BrowserChromium,
		Tags:      []string{"integration", "min", "positive"},
	}, func(a *allurex.A) {
		runPlaywrightSessionLifecycle(
			t, a, cfg,
			"Connect Playwright via hub WS endpoint (min)",
			"Close remote session",
		)
	})
}
