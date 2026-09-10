package min_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestHubPlaywrightMinHarSession_EnableHarRejectedWith400(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Playwright-min enableHAR is HTTP 400 (no session)",
		Package:   "tests.integration.HubPlaywrightMinHarRejectedTests",
		Layer:     "integration",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Playwright min session create",
		Story:     "hub v3.0.16 rejects enableHAR on *-min before create — DELETE SLA does not apply",
		Suite:     "Playwright hub min HAR reject",
		Browser:   allurex.BrowserChromium,
		Tags:      []string{"integration", "min", "negative"},
	}, func(a *allurex.A) {
		wsEndpoint, err := cfg.ResolvePlaywrightWsEndpoint()
		require.NoError(t, err)
		ws := playwrightapi.AppendQuery(wsEndpoint, "enableHAR=true")

		var before map[string]struct{}
		a.Step("Snapshot hub session ids", func() {
			before, err = hubapi.CollectSessionIDs(cfg)
			require.NoError(t, err)
		})

		a.Step("WS upgrade with enableHAR is HTTP 400", func() {
			resp, herr := playwrightapi.Handshake(cfg, ws)
			require.NoError(t, herr)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
			require.Contains(t, string(resp.Body), "enableHAR")
			require.Contains(t, string(resp.Body), "-min")
		})

		a.Step("Hub did not create a session", func() {
			after, serr := hubapi.CollectSessionIDs(cfg)
			require.NoError(t, serr)
			require.Equal(t, before, after)
		})
	})
}
