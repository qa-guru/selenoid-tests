package min_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestMain(m *testing.M) {
	if err := playwright.Install(&playwright.RunOptions{SkipInstallBrowsers: true}); err != nil {
		fmt.Fprintf(os.Stderr, "playwright driver install: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// skipUnlessWDMinReady skips when the local stack cannot start a WD min session
// (e.g. arm64 msedge-min manifest missing, or dev browsers.json out of sync).
func skipUnlessWDMinReady(t *testing.T, cfg *config.Config, browser, version string) {
	t.Helper()
	if browser == "msedge" && runtime.GOARCH == "arm64" {
		t.Skip("qaguru/webdriver-msedge:*-min has no arm64 manifest; WD min msedge runs on CI linux/amd64")
	}
	sessionID, err := hubapi.CreateSessionWithBrowser(cfg, browser, version)
	if err != nil {
		t.Skipf("WD min %s %s unavailable on this stack (Java testMin parity): %v", browser, version, err)
	}
	require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
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

func runPlaywrightSessionLifecycle(
	t *testing.T,
	a *allurex.A,
	cfg *config.Config,
	connectStep, closeStep string,
) {
	t.Helper()

	var usedBefore int
	a.Step("Snapshot hub /status used counter", func() {
		st, err := hubapi.Fetch(cfg)
		require.NoError(t, err)
		usedBefore = st.Used
	})

	pw, err := playwright.Run()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pw.Stop())
	}()

	var browser playwright.Browser
	a.Step(connectStep, func() {
		var err error
		browser, err = playwrightapi.Connect(pw, cfg, "")
		require.NoError(t, err)
	})

	a.Step("Verify remote browser is connected", func() {
		require.True(t, browser.IsConnected(), "Playwright min browser node should be connected")
	})

	a.Step("Verify hub reports active session", func() {
		st, err := hubapi.WaitUntilUsed(cfg, usedBefore+1, 30*time.Second)
		require.NoError(t, err)
		require.Equal(t, usedBefore+1, st.Used)
	})

	a.Step(closeStep, func() {
		require.NoError(t, playwrightapi.Close(browser))
	})

	a.Step("Verify browser is disconnected after close", func() {
		require.False(t, browser.IsConnected(), "Browser should disconnect after close")
	})

	a.Step("Verify hub released session", func() {
		st, err := hubapi.WaitUntilUsed(cfg, usedBefore, 30*time.Second)
		require.NoError(t, err)
		require.Equal(t, usedBefore, st.Used)
	})
}
