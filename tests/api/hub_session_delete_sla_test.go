package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSession_DeleteReturnsWithinSLA(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "DELETE /wd/hub/session/{id} returns within 5s",
		Package:   "tests.api.HubSessionDeleteSlaTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "WebDriver session API",
		Story:     "Session delete SLA",
		Suite:     "Hub session delete SLA",
		Tags:      []string{"api", "positive", "smoke"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create remote session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
			require.NotEmpty(t, sessionID)
		})
		a.Step("Delete session within SLA", func() {
			require.NoError(t, hubapi.DeleteSessionWithin(cfg, sessionID, hubapi.SessionDeleteSLA))
		})
	})
}
