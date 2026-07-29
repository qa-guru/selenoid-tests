package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubCapabilities_EchoesBrowserName(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "POST /wd/hub/session echoes browserName in capabilities",
		Package:   "tests.api.HubCapabilitiesApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "WebDriver session API",
		Story:     "WebDriver session API",
		Suite:     "Hub session capabilities",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var created *hubapi.SessionCreateResult
		a.Step("Create remote session", func() {
			var err error
			created, err = hubapi.CreateSessionWithCapabilities(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify capabilities.browserName", func() {
			require.NotEmpty(t, created.BrowserName)
			require.Equal(t, cfg.Browser, created.BrowserName)
		})
		a.Step("Delete session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, created.SessionID))
		})
	})
}
