package min_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubFirefoxMinSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Firefox min WebDriver session starts browser node and completes",
		Package:   "tests.integration.HubFirefoxMinSessionTests",
		Layer:     "integration",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session (min)",
		Story:     "WebDriver session (min)",
		Suite:     "WebDriver hub session (firefox-min)",
		Browser:   allurex.BrowserFirefox,
		Tags:      []string{"integration", "min", "positive"},
	}, func(a *allurex.A) {
		skipUnlessWDMinReady(t, cfg, "firefox", cfg.FirefoxMinVersionForSession())
		runRemoteSessionLifecycle(
			t, a, cfg,
			"firefox", cfg.FirefoxMinVersionForSession(),
			"Create Firefox min hub session",
			"Delete Firefox min session",
		)
	})
}
