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
	// Prod Create Session with VNC+video(+HAR) often exceeds 2m (sidecar + browser boot).
	createSessionTimeoutProd = 5 * time.Minute
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
	fillReactInput(t, field, manualHarSessionName)
}

// fillHubAuthFromConfig overrides baked UI defaults with credentials from the
// active test profile (remoteUrl / apiBaseUrl / uiUrl). Prod UI may ship a
// stale VITE_HUB_AUTH_* pin while nginx htpasswd has rotated.
// fillReactInput clears then types so React controlled fields pick up onChange
// (playwright-go Fill alone is occasionally a no-op on Vite/React builds).
func fillReactInput(t *testing.T, loc playwright.Locator, value string) {
	t.Helper()
	require.NoError(t, loc.Click())
	require.NoError(t, loc.Fill(""))
	require.NoError(t, loc.Type(value, playwright.LocatorTypeOptions{
		Delay: playwright.Float(5),
	}))
}

func fillHubAuthFromConfig(t *testing.T, page playwright.Page, cfg *config.Config) {
	t.Helper()
	user, pass := cfg.ResolveHubBasicAuth()
	if user == "" || pass == "" {
		// selenoid_github_e2e CI stack has no hub basic auth — UI defaults are enough.
		return
	}

	userField := page.Locator("[data-testid=capabilities-caps-auth-user]")
	passField := page.Locator("[data-testid=capabilities-caps-auth-pass]")
	accessKey := page.Locator("[data-testid=capabilities-caps-access-key-field]")
	pwPanel := page.Locator("[data-testid=capabilities-playwright-panel]")

	// Playwright panel uses accessKey only — prefer it when visible so we do not
	// fill a stale/hidden WebDriver auth duo left in the DOM.
	if n, err := pwPanel.Count(); err == nil && n > 0 {
		if vis, err := pwPanel.IsVisible(); err == nil && vis {
			require.NoError(t, accessKey.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
			fillReactInput(t, accessKey, user+":"+pass)
			return
		}
	}
	if n, err := accessKey.Count(); err == nil && n > 0 {
		if vis, err := accessKey.IsVisible(); err == nil && vis {
			fillReactInput(t, accessKey, user+":"+pass)
			return
		}
	}
	if n, err := userField.Count(); err == nil && n > 0 {
		fillReactInput(t, userField, user)
		fillReactInput(t, passField, pass)
		return
	}
	t.Fatal("no WebDriver auth duo or Playwright accessKey field on Capabilities")
}

func createSessionWait(cfg *config.Config) time.Duration {
	if cfg != nil && strings.Contains(cfg.Env, "qa_guru") {
		return createSessionTimeoutProd
	}
	return createSessionTimeout
}

func clickCreateSession(t *testing.T, page playwright.Page) string {
	t.Helper()
	cfg := config.MustLoad()
	timeout := createSessionWait(cfg)
	create := page.Locator("[data-testid=capabilities-create-session]")
	require.NoError(t, create.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}))
	enabled, err := create.IsEnabled()
	require.NoError(t, err)
	require.True(t, enabled, "Create Session button is disabled — select a browser first")
	require.NoError(t, create.Click())
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if id := sessionIDFromURL(page.URL()); id != "" {
			return id
		}
		time.Sleep(250 * time.Millisecond)
	}
	cls, _ := create.GetAttribute("class")
	body, _ := page.Locator("[data-testid=capabilities-setup]").InnerText()
	if len(body) > 600 {
		body = body[:600]
	}
	t.Fatalf("session page did not open within %s, url=%s btn=%s setup=%q", timeout, page.URL(), cls, body)
	return ""
}

func killSessionFromUI(t *testing.T, page playwright.Page) string {
	t.Helper()
	cfg := config.MustLoad()
	sessionID := sessionIDFromURL(page.URL())
	require.NotEmpty(t, sessionID, "killSessionFromUI: expected /#/sessions/<id> URL, got %s", page.URL())

	kill := page.Locator("[data-testid=session-kill]")
	// Prod Playwright Create Session can take >30s before the live panel mounts.
	killWait := createSessionTimeout
	if strings.Contains(cfg.Env, "qa_guru") {
		killWait = 3 * time.Minute
	}
	require.NoError(t, kill.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(killWait.Milliseconds())),
	}))
	require.NoError(t, kill.Click())

	timeout := killWait
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
