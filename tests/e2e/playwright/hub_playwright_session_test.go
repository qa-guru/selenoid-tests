package pw_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightSession_RemotePlaywrightSessionOpensDefaultStack(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright WS session opens default stack login",
		Package:   "tests.HubPlaywrightSessionTests",
		Layer:     "e2e",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS session",
		Story:     "Playwright WS session",
		Suite:     "Playwright WS session",
		Browser:   allurex.BrowserChromium,
		Tags:      []string{"playwright", "smoke", "positive"},
	}, func(a *allurex.A) {
		runRemotePlaywrightSmoke(t, a, cfg, func(page playwright.Page) {
			a.Step("Verify page title", func() {
				title, err := page.Title()
				require.NoError(t, err)
				require.Equal(t, config.DefaultSmokeTitle, title)
			})
			a.Step("Verify login form", func() {
				text, err := page.Locator(config.DefaultSmokeHeadingSelector).TextContent()
				require.NoError(t, err)
				require.Equal(t, config.DefaultSmokeHeading, text)
			})
			a.Step("Verify browser is connected", func() {
				require.False(t, page.IsClosed())
			})
		})
	})
}
