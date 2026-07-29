package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatus_ReturnsStatisticsJSON(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /status returns hub statistics JSON",
		Package:   "tests.api.HubStatusTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub status",
		Story:     "Hub status",
		Suite:     "Hub status",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET /status", func() {
			var err error
			status, err = hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify root counters", func() {
			require.GreaterOrEqual(t, status.Total, 0)
			require.GreaterOrEqual(t, status.Used, 0)
			require.GreaterOrEqual(t, status.Queued, 0)
			require.GreaterOrEqual(t, status.Pending, 0)
		})
		a.Step("Verify browsers map is present", func() {
			require.NotNil(t, status.Browsers)
		})
	})
}
