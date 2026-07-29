package wd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func assertBrowserVersionListed(t *testing.T, browsers map[string]any, family, version string) {
	t.Helper()
	raw, ok := browsers[family]
	require.True(t, ok, "expected %q in hub /status browsers map", family)
	versions, ok := raw.(map[string]any)
	require.True(t, ok, "expected %q browsers entry to be a version map", family)
	require.NotNil(t, versions[version], "expected version %q listed for %q", version, family)
}

func runRemoteSessionLifecycle(
	t *testing.T,
	a *allurex.A,
	cfg *config.Config,
	browserName, browserVersion string,
	createStep, deleteStep string,
) {
	t.Helper()
	var usedBefore int
	a.Step("Snapshot hub /status used counter", func() {
		st, err := hubapi.Fetch(cfg)
		require.NoError(t, err)
		usedBefore = st.Used
	})

	var sessionID string
	a.Step(createStep, func() {
		var err error
		sessionID, err = hubapi.CreateSessionWithBrowser(cfg, browserName, browserVersion)
		require.NoError(t, err)
	})

	a.Step("Verify hub reports active session", func() {
		st, err := hubapi.WaitUntilUsed(cfg, usedBefore+1, 30*time.Second)
		require.NoError(t, err)
		require.Equal(t, usedBefore+1, st.Used)
	})

	a.Step("Verify session id is assigned", func() {
		require.NotEmpty(t, sessionID)
	})

	a.Step(deleteStep, func() {
		require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
	})

	a.Step("Verify hub released session", func() {
		st, err := hubapi.WaitUntilUsed(cfg, usedBefore, 30*time.Second)
		require.NoError(t, err)
		require.Equal(t, usedBefore, st.Used)
	})
}
