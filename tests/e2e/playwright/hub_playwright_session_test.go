package pw_test

import (
	"testing"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubPlaywrightSession_RemotePlaywrightSessionOpensExampleDomain(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Remote Playwright WS session opens example.com",
		Package:   "tests.HubPlaywrightSessionTests",
		Layer:     "e2e",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS session",
		Story:     "Playwright WS session",
		Suite:     "Playwright WS session",
		Tags:      []string{"playwright", "smoke", "positive"},
	}, func(a *allurex.A) {
		runRemotePlaywrightSmoke(t, a, cfg, func(page playwright.Page) {
			a.Step("Verify page title", func() {
				title, err := page.Title()
				require.NoError(t, err)
				require.Equal(t, "Example Domain", title)
			})
			a.Step("Verify heading text", func() {
				text, err := page.Locator("h1").TextContent()
				require.NoError(t, err)
				require.Equal(t, "Example Domain", text)
			})
			a.Step("Verify browser is connected", func() {
				require.False(t, page.IsClosed())
			})
		})
	})
}
