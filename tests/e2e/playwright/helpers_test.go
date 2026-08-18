package pw_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
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

func runRemotePlaywrightSmoke(
	t *testing.T,
	a *allurex.A,
	cfg *config.Config,
	fn func(page playwright.Page),
) {
	t.Helper()

	pw, err := playwright.Run()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pw.Stop())
	}()

	var browser playwright.Browser
	a.Step("Connect Playwright via hub WS endpoint", func() {
		var err error
		browser, err = playwrightapi.Connect(pw, cfg, "")
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		if browser != nil {
			_ = playwrightapi.Close(browser)
		}
	})

	var page playwright.Page
	a.Step("Open browser context and page", func() {
		var err error
		page, err = browser.NewPage()
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		if page != nil {
			_ = page.Close()
		}
	})

	a.Step("Navigate to smoke URL", func() {
		_, err := page.Goto(cfg.SmokeURL)
		require.NoError(t, err)
		require.NoError(t, page.Locator(config.DefaultSmokeHeadingSelector).WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}))
	})

	fn(page)

	a.Step("Close remote session", func() {
		require.NoError(t, playwrightapi.Close(browser))
		browser = nil
	})
}
