package resilience_test

import (
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/stack"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiStatusRecovery_SelenoidRecoversAfterHubRestart(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "SELENOID recovers after hub restart",
		Package:   "tests.UiStatusRecoveryTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status recovery",
		Story:     "UI status recovery",
		Suite:     "UI status recovery",
		Tags:      []string{"resilience", "positive"},
	}, func(a *allurex.A) {
		baseURL, err := cfg.ResolveUiLocalBaseURL()
		require.NoError(t, err)

		a.Step("Ensure hub and UI are up", func() {
			require.NoError(t, stack.EnsureControllable(cfg))
		})

		pw, err := playwright.Run()
		require.NoError(t, err)
		defer func() { require.NoError(t, pw.Stop()) }()

		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		})
		require.NoError(t, err)

		page, err := browser.NewPage(playwright.BrowserNewPageOptions{
			Viewport: &playwright.Size{Width: 1920, Height: 1280},
		})
		require.NoError(t, err)

		a.Step("Baseline CONNECTED", func() {
			openDashboard(t, page, baseURL)
		})

		a.Step("Stop hub", func() {
			require.NoError(t, stack.KillHub(cfg))
			require.NoError(t, stack.WaitForHubDown(cfg, 15*time.Second))
		})

		a.Step("Close browser after hub stop", func() {
			require.NoError(t, browser.Close())
		})

		a.Step("UI /status reports hub error", func() {
			resp, err := uiapi.FetchStatusWhenHubUnavailable(cfg)
			require.NoError(t, err)
			require.NotEmpty(t, resp.Errors)
		})

		a.Step("Restart hub", func() {
			require.NoError(t, stack.StartHubDetached())
			require.NoError(t, stack.WaitForHubReady(cfg, 30*time.Second))
		})

		a.Step("Both indicators recover within UI retry window", func() {
			browser2, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(true),
			})
			require.NoError(t, err)
			defer func() { require.NoError(t, browser2.Close()) }()

			page2, err := browser2.NewPage(playwright.BrowserNewPageOptions{
				Viewport: &playwright.Size{Width: 1920, Height: 1280},
			})
			require.NoError(t, err)

			openDashboard(t, page2, baseURL)
			stayConnected(t, page2, 5*time.Second)
		})
	})
}
