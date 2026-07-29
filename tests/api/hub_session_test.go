package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSession_CreateAndDelete(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "POST /wd/hub/session creates and DELETE removes session",
		Package:   "tests.api.HubSessionApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "WebDriver session API",
		Story:     "WebDriver session API",
		Suite:     "Hub session API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create remote session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify session id", func() {
			require.NotEmpty(t, sessionID)
		})
		a.Step("Delete session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
	})
}
