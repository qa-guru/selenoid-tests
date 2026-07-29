package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubDownload_UnknownSessionReturns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /download/{sessionId} for unknown session returns 404 JSON",
		Package:   "tests.api.HubDownloadApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub download proxy",
		Story:     "Hub download proxy API",
		Suite:     "Hub download proxy API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET /download/unknown-session", func() {
			require.NoError(t, hubapi.GetProxyExpectStatus(cfg, "/download", "unknown-session", http.StatusNotFound))
		})
	})
}
