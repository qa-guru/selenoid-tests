package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiPing_ReturnsUptimeAndVersion(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /ping returns uptime and version",
		Package:   "tests.api.UiPingTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI ping",
		Story:     "UI ping API",
		Suite:     "UI ping API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var ping *uiapi.Ping
		a.Step("GET /ping", func() {
			var err error
			ping, err = uiapi.FetchPing(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify uptime and version", func() {
			require.NotEmpty(t, ping.Uptime)
			require.NotEmpty(t, ping.Version)
		})
	})
}
