package wd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSession_RemoteSessionOpensDefaultStack(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Chrome session opens default stack login",
		Package:   "tests.HubSessionTests",
		Layer:     "e2e",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "WebDriver session",
		Suite:     "WebDriver session",
		Browser:   allurex.BrowserChrome,
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runRemoteSmokeSession(t, a, cfg, func(sessionID string) {
			a.Step("Verify session id is assigned", func() {
				require.NotEmpty(t, sessionID)
			})
			a.Step("Verify page title", func() {
				title, err := hubapi.GetSessionTitle(cfg, sessionID)
				require.NoError(t, err)
				require.Equal(t, config.DefaultSmokeTitle, title)
			})
			a.Step("Verify login form", func() {
				text, err := hubapi.WaitElementTextBySelector(cfg, sessionID, config.DefaultSmokeHeadingSelector, 15*time.Second)
				require.NoError(t, err)
				require.Equal(t, config.DefaultSmokeHeading, text)
			})
		})
	})
}
