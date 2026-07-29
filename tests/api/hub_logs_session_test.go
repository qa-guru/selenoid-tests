package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubLogsSession_SessionLogsRequireWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /logs/{sessionId} without WebSocket upgrade returns 400 for active session",
		Package:   "tests.api.HubLogsSessionApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub logs WebSocket",
		Story:     "Hub logs WebSocket",
		Suite:     "Hub session logs WebSocket API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		a.Step("GET /logs/{sessionId} without WebSocket headers", func() {
			require.NoError(t, getExpectStatusBodyContains(hubHTTP(cfg), "/logs/"+sessionID, http.StatusBadRequest, "websocket"))
		})
		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
	})
}
