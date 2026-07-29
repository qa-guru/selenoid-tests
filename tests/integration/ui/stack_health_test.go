package ui_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestStackHealth_HubAndUiRespondOnHealthEndpoints(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Hub and UI respond on health endpoints",
		Package:   "tests.integration.StackHealthTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "Stack health",
		Story:     "Stack health integration",
		Suite:     "Stack health integration",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		a.Step("GET hub /status", func() {
			_, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		var ping *uiapi.Ping
		a.Step("GET UI /ping", func() {
			var err error
			ping, err = uiapi.FetchPing(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify UI ping uptime", func() {
			require.NotEmpty(t, ping.Uptime)
		})
	})
}
