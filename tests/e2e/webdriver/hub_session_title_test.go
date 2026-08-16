package wd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSessionTitle_RemoteSessionHasExampleDomainTitle(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote session page has Example Domain title",
		Package:   "tests.HubSessionTitleTests",
		Layer:     "e2e",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session",
		Story:     "Hub session title",
		Suite:     "Hub session title",
		Browser:   cfg.Browser,
		Tags:      []string{"smoke", "positive"},
	}, func(a *allurex.A) {
		runRemoteSmokeSession(t, a, cfg, func(sessionID string) {
			a.Step("Verify document title", func() {
				title, err := hubapi.GetSessionTitle(cfg, sessionID)
				require.NoError(t, err)
				require.Equal(t, "Example Domain", title)
			})
		})
	})
}
