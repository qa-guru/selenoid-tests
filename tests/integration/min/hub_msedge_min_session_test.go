package min_test

import (
	"testing"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubMsedgeMinSession_RemoteSessionStartsAndCompletes(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Edge min WebDriver session starts browser node and completes",
		Package:   "tests.integration.HubMsedgeMinSessionTests",
		Layer:     "integration",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session (min)",
		Story:     "WebDriver session (min)",
		Suite:     "WebDriver hub session (msedge-min)",
		Browser:   allurex.BrowserMsedge,
		Tags:      []string{"integration", "min", "positive"},
	}, func(a *allurex.A) {
		skipUnlessWDMinReady(t, cfg, "msedge", cfg.MsedgeMinVersionForSession())
		runRemoteSessionLifecycle(
			t, a, cfg,
			"msedge", cfg.MsedgeMinVersionForSession(),
			"Create Edge min hub session",
			"Delete Edge min session",
		)
	})
}
