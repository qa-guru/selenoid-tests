package ui_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

const (
	manualHarSessionName = "ui-manual-har-e2e"
	createSessionTimeout = 120 * time.Second
)

func openNewSession(t *testing.T, page playwright.Page, baseURL string) {
	t.Helper()
	_, err := page.Goto(baseURL+"/#/new-session", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)
	require.NoError(t, page.Locator("[data-testid=capabilities-setup]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}))
}

func webdriverChipValue(cfg *config.Config) string {
	return fmt.Sprintf("chrome_%s", cfg.ChromeVersionForSession())
}

func playwrightChipValue(cfg *config.Config) (string, error) {
	endpoint, err := cfg.ResolvePlaywrightWsEndpoint()
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("cannot parse playwright browser from %q", endpoint)
	}
	version := parts[len(parts)-1]
	if idx := strings.Index(version, "?"); idx >= 0 {
		version = version[:idx]
	}
	browser := parts[len(parts)-2]
	return fmt.Sprintf("%s_%s", browser, version), nil
}

func selectWebDriverChrome(t *testing.T, page playwright.Page, cfg *config.Config) {
	t.Helper()
	chip := page.Locator("[data-testid=capabilities-browser-select] button[data-value='" + webdriverChipValue(cfg) + "']")
	require.NoError(t, chip.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, chip.Click())
	require.NoError(t, page.Locator("[data-testid=capabilities-caps]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}))
}

func selectPlaywrightChrome(t *testing.T, page playwright.Page, cfg *config.Config) {
	t.Helper()
	value, err := playwrightChipValue(cfg)
	require.NoError(t, err)
	chip := page.Locator("[data-testid=capabilities-browser-select-playwright] button[data-value='" + value + "']")
	require.NoError(t, chip.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, chip.Click())
	require.NoError(t, page.Locator("[data-testid=capabilities-playwright-panel]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}))
}

func setSegTrue(t *testing.T, page playwright.Page, fieldTestID string) {
	t.Helper()
	btn := page.Locator(fmt.Sprintf("[data-testid=%s] button[data-value=true]", fieldTestID))
	require.NoError(t, btn.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, btn.Click())
}

func setManualHarSessionName(t *testing.T, page playwright.Page, fieldTestID string) {
	t.Helper()
	field := page.Locator(fmt.Sprintf("[data-testid=%s]", fieldTestID))
	require.NoError(t, field.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, field.Fill(manualHarSessionName))
}

// fillHubAuthFromConfig overrides baked UI defaults with credentials from the
// active test profile (remoteUrl / apiBaseUrl / uiUrl). Prod UI may ship a
// stale VITE_HUB_AUTH_* pin while nginx htpasswd has rotated.
func fillHubAuthFromConfig(t *testing.T, page playwright.Page, cfg *config.Config) {
	t.Helper()
	user, pass := cfg.ResolveHubBasicAuth()
	require.NotEmpty(t, user, "hub basic auth user missing from test config")
	require.NotEmpty(t, pass, "hub basic auth password missing from test config")

	userField := page.Locator("[data-testid=capabilities-caps-auth-user]")
	passField := page.Locator("[data-testid=capabilities-caps-auth-pass]")
	accessKey := page.Locator("[data-testid=capabilities-caps-access-key-field]")

	if n, err := userField.Count(); err == nil && n > 0 {
		require.NoError(t, userField.Fill(user))
		require.NoError(t, passField.Fill(pass))
		return
	}
	if n, err := accessKey.Count(); err == nil && n > 0 {
		require.NoError(t, accessKey.Fill(user+":"+pass))
		return
	}
	t.Fatal("no WebDriver auth duo or Playwright accessKey field on Capabilities")
}

func clickCreateSession(t *testing.T, page playwright.Page) string {
	t.Helper()
	create := page.Locator("[data-testid=capabilities-create-session]")
	require.NoError(t, create.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	enabled, err := create.IsEnabled()
	require.NoError(t, err)
	require.True(t, enabled, "Create Session button is disabled — select a browser first")
	require.NoError(t, create.Click())
	deadline := time.Now().Add(createSessionTimeout)
	for time.Now().Before(deadline) {
		if id := sessionIDFromURL(page.URL()); id != "" {
			return id
		}
		time.Sleep(250 * time.Millisecond)
	}
	t.Fatalf("session page did not open within %s, url=%s", createSessionTimeout, page.URL())
	return ""
}

func killSessionFromUI(t *testing.T, page playwright.Page) string {
	t.Helper()
	cfg := config.MustLoad()
	sessionID := sessionIDFromURL(page.URL())
	require.NotEmpty(t, sessionID, "killSessionFromUI: expected /#/sessions/<id> URL, got %s", page.URL())

	kill := page.Locator("[data-testid=session-kill]")
	require.NoError(t, kill.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, kill.Click())

	timeout := createSessionTimeout
	if strings.Contains(cfg.Env, "qa_guru") {
		timeout = 3 * time.Minute
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		url := page.URL()
		require.Equal(t, sessionID, sessionIDFromURL(url),
			"kill must keep URL on /#/sessions/%s, got %s", sessionID, url)

		count, err := kill.Count()
		if err == nil && count == 0 {
			return sessionID
		}
		finished, err := page.Locator("text=FINISHED").Count()
		if err == nil && finished > 0 {
			return sessionID
		}
		time.Sleep(250 * time.Millisecond)
	}
	t.Fatalf("session kill did not finish within %s, url=%s", timeout, page.URL())
	return sessionID
}

func openFinishedSessionsArchive(t *testing.T, page playwright.Page, baseURL string) {
	t.Helper()
	_, err := page.Goto(baseURL+"/#/sessions", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)
	require.NoError(t, page.Locator("[data-testid=archive-panel]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}))
}

func filterArchiveSession(t *testing.T, page playwright.Page, sessionID string) {
	t.Helper()
	filter := page.Locator("[data-testid=session-filter-input]")
	require.NoError(t, filter.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, filter.Fill(sessionID))
	time.Sleep(500 * time.Millisecond)
}

func waitHubArchivedHar(t *testing.T, cfg *config.Config, sessionID string) {
	t.Helper()
	timeout := 45 * time.Second
	if strings.Contains(cfg.Env, "qa_guru") {
		timeout = 5 * time.Minute
	}
	row, err := hubapi.WaitArchivedSessionHar(cfg, sessionID, timeout)
	require.NoError(t, err)
	require.NotNil(t, row)
	require.NotEmpty(t, row.HAR)
}

func waitArchiveHarIcon(t *testing.T, page playwright.Page, baseURL, sessionID string, timeout time.Duration) {
	t.Helper()
	filterArchiveSession(t, page, sessionID)
	row := page.Locator(fmt.Sprintf("[data-session='%s']", sessionID))
	icon := row.Locator("[data-testid=artifact-har]")
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		count, err := icon.Count()
		if err == nil && count > 0 {
			require.NoError(t, icon.First().WaitFor(playwright.LocatorWaitForOptions{
				State: playwright.WaitForSelectorStateVisible,
			}))
			return
		}
		_, err = page.Goto(baseURL+"/#/sessions", playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		})
		require.NoError(t, err)
		require.NoError(t, page.Locator("[data-testid=archive-panel]").WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		}))
		time.Sleep(400 * time.Millisecond)
	}
	t.Fatalf("HAR artifact icon did not appear for session %s within %s", sessionID, timeout)
}

func openArchiveSessionDetail(t *testing.T, page playwright.Page, sessionID string) {
	t.Helper()
	link := page.Locator(fmt.Sprintf("[data-session='%s'] [data-testid=session-detail-link]", sessionID))
	require.NoError(t, link.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	require.NoError(t, link.Click())
	require.NoError(t, page.Locator("[data-testid=session-page]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}))
}
