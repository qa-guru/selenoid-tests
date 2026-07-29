package cmapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

// skipUnlessCmSessionReady skips when CM hub cannot start a WD session (local arm64 + 149.0-min flake).
// CI linux/amd64 canon: selenoid_github_cm_integration + start-ci-cm-stack.sh.
func skipUnlessCmSessionReady(t *testing.T, cfg *config.Config) {
	t.Helper()
	sessionID, err := hubapi.CreateSessionWithBrowser(cfg, cfg.Browser, cfg.BrowserVersion)
	if err != nil {
		t.Skipf("CM WD session %s %s unavailable on this stack (Java testCmApi parity needs CI linux/amd64 or -DbrowserVersion=149.0): %v",
			cfg.Browser, cfg.BrowserVersion, err)
	}
	require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
}
