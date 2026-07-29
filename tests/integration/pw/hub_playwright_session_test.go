package pw_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright WS session starts browser node and completes",
		Package:   "tests.integration.HubPlaywrightSessionTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS session",
		Story:     "Playwright WS session",
		Suite:     "Playwright hub WS session",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		runPlaywrightSessionLifecycle(
			t, a, cfg, "",
			"Connect Playwright via hub WS endpoint",
			"Close remote session",
		)
	})
}
