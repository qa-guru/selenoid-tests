package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubError_ReturnsInvalidSessionJSON(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /error returns invalid session id JSON",
		Package:   "tests.api.HubErrorApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub error",
		Story:     "Hub /error API",
		Suite:     "Hub /error API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET /error", func() {
			require.NoError(t, hubapi.FetchErrorExpectInvalidSession(cfg))
		})
	})
}
