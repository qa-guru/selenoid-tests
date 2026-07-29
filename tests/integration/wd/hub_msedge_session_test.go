package wd_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubMsedgeSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Edge WebDriver session starts browser node and completes",
		Package:   "tests.integration.HubMsedgeSessionIntegrationTests",
		Layer:     "integration",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "WebDriver session",
		Suite:     "WebDriver hub session (msedge)",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		runRemoteSessionLifecycle(
			t, a, cfg,
			"msedge", cfg.MsedgeVersionForSession(),
			"Create Edge hub session",
			"Delete Edge session",
		)
	})
}
