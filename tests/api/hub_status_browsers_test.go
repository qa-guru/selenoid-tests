package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatusBrowsers_ListsChromeFamily(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /status lists configured browser families",
		Package:   "tests.api.HubStatusBrowsersTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub status",
		Story:     "Hub status browsers",
		Suite:     "Hub status browsers",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET /status", func() {
			var err error
			status, err = hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify chrome is configured", func() {
			require.NotNil(t, status.Browsers["chrome"])
		})
	})
}
