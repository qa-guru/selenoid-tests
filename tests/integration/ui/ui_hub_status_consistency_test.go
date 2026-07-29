package ui_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiHubStatusConsistency_UiStatusMirrorsHubCounters(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "UI /status state mirrors hub /status counters",
		Package:   "tests.integration.UiHubStatusConsistencyTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI hub proxy",
		Story:     "UI hub status consistency",
		Suite:     "UI hub status consistency",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		var hubStatus, uiStatus *hubapi.Status
		a.Step("GET hub /status", func() {
			var err error
			hubStatus, err = hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		a.Step("GET UI /status", func() {
			var err error
			uiStatus, err = uiapi.FetchStatus(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify proxied counters match hub", func() {
			require.Equal(t, hubStatus.Total, uiStatus.Total)
			require.Equal(t, hubStatus.Used, uiStatus.Used)
			require.Equal(t, hubStatus.Queued, uiStatus.Queued)
			require.Equal(t, hubStatus.Pending, uiStatus.Pending)
		})
	})
}
