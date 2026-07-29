package cme2e_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func skipUnlessCmSessionReady(t *testing.T, cfg *config.Config) {
	t.Helper()
	sessionID, err := hubapi.CreateSessionWithBrowser(cfg, cfg.Browser, cfg.BrowserVersion)
	if err != nil {
		t.Skipf("CM e2e WD session %s %s unavailable (CI linux/amd64 canon): %v",
			cfg.Browser, cfg.BrowserVersion, err)
	}
	require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
}
