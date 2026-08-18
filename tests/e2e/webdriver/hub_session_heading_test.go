package wd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSessionHeading_RemoteSessionRendersLoginForm(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote session renders default stack login form",
		Package:   "tests.HubSessionHeadingTests",
		Layer:     "e2e",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "Hub session heading",
		Suite:     "Hub session heading",
		Browser:   cfg.Browser,
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runRemoteSmokeSession(t, a, cfg, func(sessionID string) {
			a.Step("Verify login form", func() {
				text, err := hubapi.WaitElementTextBySelector(cfg, sessionID, config.DefaultSmokeHeadingSelector, 15*time.Second)
				require.NoError(t, err)
				require.Equal(t, config.DefaultSmokeHeading, text)
			})
		})
	})
}
