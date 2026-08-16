package wd_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubChromeWarmSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Chrome warm WebDriver session starts browser node and completes",
		Package:   "tests.integration.HubChromeWarmSessionIntegrationTests",
		Layer:     "integration",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "WebDriver session",
		Suite:     "WebDriver hub session (chrome warm)",
		Browser:   allurex.BrowserChrome,
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		runRemoteSessionLifecycle(
			t, a, cfg,
			"chrome", cfg.ChromeVersionForSession(),
			"Create Chrome warm hub session",
			"Delete Chrome warm session",
		)
	})
}
