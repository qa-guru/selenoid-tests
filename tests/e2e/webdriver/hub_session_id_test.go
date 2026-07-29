package wd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubSessionId_RemoteSessionAssignsSessionId(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote session assigns session id",
		Package:   "tests.HubSessionIdTests",
		Layer:     "e2e",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "Hub session id",
		Suite:     "Hub session id",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runRemoteSmokeSession(t, a, cfg, func(sessionID string) {
			a.Step("Verify session id", func() {
				require.NotEmpty(t, sessionID)
			})
		})
	})
}
