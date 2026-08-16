package pw_test

import (
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestHubPlaywrightNavigate_ExampleDomain(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Playwright session navigates to example.com",
		Package:   "tests.integration.HubPlaywrightNavigateTests",
		Layer:     "integration",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS session",
		Story:     "Playwright WS session",
		Suite:     "Playwright navigate integration",
		Browser:   allurex.BrowserChromium,
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		var usedBefore int
		a.Step("Snapshot hub used counter", func() {
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
		a.Step("Connect Playwright via hub WS endpoint", func() {
			var err error
			browser, err = playwrightapi.Connect(pw, cfg, "")
			require.NoError(t, err)
		})

		a.Step("Navigate to example.com", func() {
			page, err := browser.NewPage()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, page.Close())
			}()
			_, err = page.Goto("https://example.com/")
			require.NoError(t, err)
			title, err := page.Title()
			require.NoError(t, err)
			require.Equal(t, "Example Domain", title)
		})

		a.Step("Close remote session", func() {
			require.NoError(t, playwrightapi.Close(browser))
		})

		a.Step("Verify hub released session", func() {
			st, err := hubapi.WaitUntilUsed(cfg, usedBefore, 30*time.Second)
			require.NoError(t, err)
			require.Equal(t, usedBefore, st.Used)
		})
	})
}
