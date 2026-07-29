package ui_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiSseWithSession_SsePayloadParseableWithActiveSession(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "SSE payload remains parseable while hub session is active",
		Package:   "tests.integration.UiSseWithSessionTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI SSE",
		Story:     "UI SSE",
		Suite:     "UI SSE with hub session",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		defer func() {
			a.Step("Delete hub session", func() {
				require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			})
		}()

		a.Step("Verify SSE payload has state or errors", func() {
			event, err := uiapi.ReadFirstSSEEvent(cfg)
			require.NoError(t, err)
			require.True(t, event.HasState() || event.HasErrors(),
				"SSE event should carry hub state or errors")
		})
	})
}
