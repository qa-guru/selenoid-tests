package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiStatus_ReturnsProxiedHubJSON(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET UI /status returns proxied hub statistics JSON",
		Package:   "tests.api.UiStatusTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status proxy",
		Story:     "UI status proxy",
		Suite:     "UI status proxy",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET UI /status", func() {
			var err error
			status, err = uiapi.FetchStatus(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify proxied state counters", func() {
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
