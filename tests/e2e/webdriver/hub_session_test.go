package wd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSession_RemoteSessionOpensExampleDomain(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Chrome session opens example.com",
		Package:   "tests.HubSessionTests",
		Layer:     "e2e",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "WebDriver session",
		Suite:     "WebDriver session",
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runRemoteSmokeSession(t, a, cfg, func(sessionID string) {
			a.Step("Verify session id is assigned", func() {
				require.NotEmpty(t, sessionID)
			})
			a.Step("Verify page title", func() {
				title, err := hubapi.GetSessionTitle(cfg, sessionID)
				require.NoError(t, err)
				require.Equal(t, "Example Domain", title)
			})
			a.Step("Verify heading text", func() {
				text, err := hubapi.GetElementTextBySelector(cfg, sessionID, "h1")
				require.NoError(t, err)
				require.Equal(t, "Example Domain", text)
			})
		})
	})
}
