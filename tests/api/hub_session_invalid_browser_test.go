package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSessionInvalidBrowser_RejectsUnknownBrowser(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "POST /wd/hub/session rejects unknown browser",
		Package:   "tests.api.HubSessionInvalidBrowserTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "WebDriver session API",
		Story:     "Hub session invalid browser",
		Suite:     "Hub session invalid browser",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("POST session with unknown browser", func() {
			require.NoError(t, hubapi.CreateSessionExpectStatus(cfg, "unknown-browser", "1.0", http.StatusBadRequest))
		})
	})
}
