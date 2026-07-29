package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubClipboard_UnknownSessionReturns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /clipboard/{sessionId} for unknown session returns 404 JSON",
		Package:   "tests.api.HubClipboardApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub clipboard proxy",
		Story:     "Hub clipboard proxy API",
		Suite:     "Hub clipboard proxy API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET /clipboard/unknown-session", func() {
			require.NoError(t, hubapi.GetProxyExpectStatus(cfg, "/clipboard", "unknown-session", http.StatusNotFound))
		})
	})
}
