package ui_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiStatusWithSession_UsedCounterIncrementsWithHubSession(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "UI /status used counter increments while hub session is active",
		Package:   "tests.integration.UiStatusWithSessionTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI hub proxy",
		Story:     "UI hub proxy",
		Suite:     "UI status with active session",
		Tags:      []string{"integration", "positive"},
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

		a.Step("Verify hub and UI report incremented used counter", func() {
			hubStatus, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
			uiStatus, err := uiapi.FetchStatus(cfg)
			require.NoError(t, err)
			require.Equal(t, usedBefore+1, hubStatus.Used)
			require.Equal(t, hubStatus.Used, uiStatus.Used)
		})

		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})

		a.Step("Verify used counter returned to baseline", func() {
			st, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
			require.Equal(t, usedBefore, st.Used)
		})
	})
}
