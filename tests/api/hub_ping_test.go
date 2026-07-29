package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubPing_ReturnsUptimeAndVersion(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /ping returns uptime and version",
		Package:   "tests.api.HubPingTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub ping",
		Story:     "Hub ping API",
		Suite:     "Hub ping API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var ping *hubapi.Ping
		a.Step("GET /ping", func() {
			var err error
			ping, err = hubapi.FetchPing(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify uptime and version", func() {
			require.NotEmpty(t, ping.Uptime)
			require.NotEmpty(t, ping.Version)
		})
	})
}
