package wd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatusSession_StatusBrowsersReflectActiveSession(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Hub /status browsers map reflects active session browser family",
		Package:   "tests.integration.HubStatusSessionIntegrationTests",
		Layer:     "integration",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub status with session",
		Story:     "Hub status with session",
		Suite:     "Hub status browsers with active session",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		defer func() {
			a.Step("Delete hub session", func() {
				require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			})
		}()

		a.Step("GET hub /status", func() {
			status, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
			require.Equal(t, 1, status.Used)
			require.NotNil(t, status.Browsers)
			assertBrowserVersionListed(t, status.Browsers, "chrome", cfg.BrowserVersion)
		})
	})
}
