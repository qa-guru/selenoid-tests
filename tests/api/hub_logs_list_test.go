package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestHubLogsList_SessionLogsRequireWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /logs/{sessionId} without WebSocket upgrade returns 400",
		Package:   "tests.api.HubLogsListApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub logs",
		Story:     "Hub logs endpoint API",
		Suite:     "Hub logs endpoint API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET /logs/unknown-session without WebSocket headers", func() {
			require.NoError(t, getExpectStatusBodyContains(hubHTTP(cfg), "/logs/unknown-session", http.StatusBadRequest, "websocket"))
		})
	})
}
