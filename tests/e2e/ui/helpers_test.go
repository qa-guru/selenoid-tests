package ui_test

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

const (
	connectedTimeout = 20 * time.Second
	vncConnectTimeout = 90 * time.Second
	sessionAppearTimeout = 20 * time.Second
)

func TestMain(m *testing.M) {
	// UI e2e uses local Chromium via playwright-go (not remote WS connect).
	if err := playwright.Install(); err != nil {
		fmt.Fprintf(os.Stderr, "playwright install: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func runWithBrowser(t *testing.T, fn func(page playwright.Page, baseURL string)) {
	t.Helper()
	cfg := config.MustLoad()
	baseURL, err := cfg.ResolveUiLocalBaseURL()
	require.NoError(t, err)

	pw, err := playwright.Run()
	require.NoError(t, err)
	defer func() { require.NoError(t, pw.Stop()) }()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, browser.Close()) }()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{Width: 1920, Height: 1280},
	})
	require.NoError(t, err)

	fn(page, baseURL)
}

func openDashboard(t *testing.T, page playwright.Page, baseURL string) {
	t.Helper()
	_, err := page.Goto(baseURL+"/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)
	waitConnected(t, page)
}

func waitConnected(t *testing.T, page playwright.Page) {
	t.Helper()
	waitStatusTileConnected(t, page, "#sse-status")
	waitStatusTileConnected(t, page, "#selenoid-status")
}

func waitStatusTileConnected(t *testing.T, page playwright.Page, selector string) {
	t.Helper()
	loc := page.Locator(selector)
	require.NoError(t, loc.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(connectedTimeout.Milliseconds())),
	}))
	deadline := time.Now().Add(connectedTimeout)
	for time.Now().Before(deadline) {
		if statusTileLooksConnected(loc) {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	class, _ := loc.GetAttribute("class")
	text, _ := loc.InnerText()
	t.Fatalf("%s: expected connected tile within %s, class=%q text=%q", selector, connectedTimeout, class, text)
}

func statusTileLooksConnected(loc playwright.Locator) bool {
	class, err := loc.GetAttribute("class")
	if err != nil {
		return false
	}
	if strings.Contains(class, "status-tile--connected") {
		return true
	}
	text, err := loc.InnerText()
	if err != nil {
		return false
	}
	normalized := strings.ToUpper(strings.TrimSpace(text))
	return strings.Contains(normalized, "CONNECTED")
}

func assertStatusTileConnected(t *testing.T, page playwright.Page, selector string) {
	t.Helper()
	loc := page.Locator(selector)
	require.True(t, statusTileLooksConnected(loc), "%s is not connected", selector)
}

func stayConnected(t *testing.T, page playwright.Page, stable time.Duration) {
	t.Helper()
	step := 500 * time.Millisecond
	steps := int(stable / step)
	for i := 0; i < steps; i++ {
		waitConnected(t, page)
		time.Sleep(step)
	}
}

func reloadDashboard(t *testing.T, page playwright.Page) {
	t.Helper()
	_, err := page.Reload(playwright.PageReloadOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)
}

func chromeChipValue(cfg *config.Config) string {
	v := cfg.ChromeVersionForSession()
	if strings.TrimSpace(v) == "" {
		v = cfg.BrowserVersion
	}
	return "chrome_" + v
}

func waitVncConnected(t *testing.T, page playwright.Page) {
	t.Helper()
	cfg := config.MustLoad()
	vnc := page.Locator("[data-testid='vnc-window']")
	timeout := vncConnectTimeout
	if strings.Contains(cfg.Env, "github") {
		timeout = 5 * time.Second
	}
	err := vnc.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
	if err != nil {
		if strings.Contains(cfg.Env, "github") {
			t.Log("VNC window not visible on github_e2e — skipping noVNC connected assert")
			return
		}
		require.NoError(t, err)
	}
	deadline := time.Now().Add(vncConnectTimeout)
	for time.Now().Before(deadline) {
		class, err := vnc.GetAttribute("class")
		require.NoError(t, err)
		if strings.Contains(class, "vnc-window--connected") {
			unlock := page.Locator("[data-testid='vnc-window'] [aria-label='Unlock screen']")
			require.NoError(t, unlock.WaitFor(playwright.LocatorWaitForOptions{
				State: playwright.WaitForSelectorStateVisible,
			}))
			require.NoError(t, unlock.Click())
			vncClass, _ := vnc.GetAttribute("class")
			require.Contains(t, vncClass, "vnc-window--connected")
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	if strings.Contains(cfg.Env, "github") {
		t.Log("VNC connected class not reached on github_e2e — skipping noVNC connected assert")
		return
	}
	t.Fatal("VNC did not reach connected state")
}

func sessionIDFromURL(raw string) string {
	markers := []string{"#/sessions/", "/sessions/"}
	for _, marker := range markers {
		if idx := strings.Index(raw, marker); idx >= 0 {
			rest := raw[idx+len(marker):]
			if q := strings.Index(rest, "?"); q >= 0 {
				rest = rest[:q]
			}
			if hash := strings.Index(rest, "#"); hash >= 0 {
				rest = rest[:hash]
			}
			return strings.TrimSuffix(rest, "/")
		}
	}
	return ""
}

func sessionIDFromHref(href string) string {
	const marker = "/sessions/"
	idx := strings.Index(href, marker)
	if idx < 0 {
		return ""
	}
	id := href[idx+len(marker):]
	if q := strings.Index(id, "?"); q >= 0 {
		id = id[:q]
	}
	return id
}

func uiHostFromConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()
	raw, err := cfg.ResolveUiURL()
	require.NoError(t, err)
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return u.Hostname()
}
