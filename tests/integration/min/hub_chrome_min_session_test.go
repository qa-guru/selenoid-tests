package min_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubChromeMinSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Chrome min WebDriver session starts browser node and completes",
		Package:   "tests.integration.HubChromeMinSessionTests",
		Layer:     "integration",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session (min)",
		Story:     "WebDriver session (min)",
		Suite:     "WebDriver hub session (chrome-min)",
		Tags:      []string{"integration", "min", "positive"},
	}, func(a *allurex.A) {
		skipUnlessWDMinReady(t, cfg, "chrome", cfg.ChromeMinVersionForSession())
		runRemoteSessionLifecycle(
			t, a, cfg,
			"chrome", cfg.ChromeMinVersionForSession(),
			"Create Chrome min hub session",
			"Delete Chrome min session",
		)
	})
}
