package pw_test

import (
	"fmt"
	"os"
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
	// Remote WS connect only — install playwright driver shell, skip local browsers (Java testPlaywright parity).
	if err := playwright.Install(&playwright.RunOptions{SkipInstallBrowsers: true}); err != nil {
		fmt.Fprintf(os.Stderr, "playwright driver install: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func runPlaywrightSessionLifecycle(
	t *testing.T,
	a *allurex.A,
	cfg *config.Config,
	wsEndpoint string,
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
		browser, err = playwrightapi.Connect(pw, cfg, wsEndpoint)
		require.NoError(t, err)
	})

	a.Step("Verify remote browser is connected", func() {
		require.True(t, browser.IsConnected(), "Playwright browser node should be connected")
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
