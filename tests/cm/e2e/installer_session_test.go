package cme2e_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/cm"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestCmInstallerSession_RemoteSessionOpensDefaultStack(t *testing.T) {
	cfg := config.MustLoad()
	installer, err := cm.WithTempConfigDir(cfg)
	require.NoError(t, err)
	installer.StopAll()
	t.Cleanup(func() {
		installer.StopAll()
		installer.DeleteConfigDir()
	})

	result, err := installer.Configure()
	require.NoError(t, err)
	result.RequireSuccess("configure")
	result, err = installer.StartHub()
	require.NoError(t, err)
	result.RequireSuccess("start hub")
	require.NoError(t, installer.WaitForHubReady(60*time.Second))
	skipUnlessCmSessionReady(t, cfg)

	allurex.Run(t, allurex.Meta{
		Name:      "Remote Chrome session opens default stack login",
		Package:   "tests.CmInstallerSessionTests",
		Layer:     "e2e",
		Component: "cm",
		Epic:      "cm",
		Feature:   "Installed stack session",
		Story:     "Installed stack session",
		Suite:     "CM installer session",
		Tags:      []string{"smoke", "cm", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create remote WebDriver session", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithBrowser(cfg, cfg.Browser, cfg.BrowserVersion)
			require.NoError(t, err)
		})
		t.Cleanup(func() {
			if sessionID != "" {
				_ = hubapi.DeleteSession(cfg, sessionID)
			}
		})

		a.Step("Open smoke URL via remote WebDriver", func() {
			require.NoError(t, hubapi.NavigateSession(cfg, sessionID, cfg.SmokeURL))
		})
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
}
