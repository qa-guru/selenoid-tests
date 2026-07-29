package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubSessionDeleteUnknown_Returns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "DELETE unknown session returns client error",
		Package:   "tests.api.HubSessionDeleteUnknownTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "WebDriver session API",
		Story:     "Hub session delete unknown",
		Suite:     "Hub session delete unknown",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("DELETE unknown session id", func() {
			require.NoError(t, hubapi.DeleteSessionExpectStatus(cfg, "missing-session-id", http.StatusNotFound))
		})
	})
}
