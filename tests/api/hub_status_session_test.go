package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatusSession_UsedCounterTracksLifecycle(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /status used counter tracks session lifecycle",
		Package:   "tests.api.HubStatusSessionApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub status with session",
		Story:     "Hub status with session",
		Suite:     "Hub status session counters API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var usedBefore int
		a.Step("Snapshot hub used counter", func() {
			st, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
			usedBefore = st.Used
		})
		var sessionID string
		a.Step("Create hub session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify used counter incremented", func() {
			st, err := hubapi.WaitUntilUsed(cfg, usedBefore+1, 30*time.Second)
			require.NoError(t, err)
			require.Equal(t, usedBefore+1, st.Used)
		})
		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
		a.Step("Verify used counter returned to baseline", func() {
			st, err := hubapi.WaitUntilUsed(cfg, usedBefore, 30*time.Second)
			require.NoError(t, err)
			require.Equal(t, usedBefore, st.Used)
		})
	})
}
