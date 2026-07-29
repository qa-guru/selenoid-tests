package wd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func runRemoteSmokeSession(
	t *testing.T,
	a *allurex.A,
	cfg *config.Config,
	fn func(sessionID string),
) {
	t.Helper()
	browser := cfg.Browser
	version := cfg.BrowserVersion

	var sessionID string
	a.Step("Create remote WebDriver session", func() {
		var err error
		sessionID, err = hubapi.CreateSessionWithBrowser(cfg, browser, version)
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		if sessionID != "" {
			_ = hubapi.DeleteSession(cfg, sessionID)
		}
	})

	a.Step("Open smoke URL", func() {
		require.NoError(t, hubapi.NavigateSession(cfg, sessionID, cfg.SmokeURL))
	})

	fn(sessionID)

	a.Step("Delete remote session", func() {
		require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		sessionID = ""
	})
}
