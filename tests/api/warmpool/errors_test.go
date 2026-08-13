package warmpool_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestWarmPoolRelease_ValidationErrors(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/release returns 400 without slotId and 404 for unknown id",
		Package:   "tests.api.warmpool.WarmPoolReleaseErrorApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Release errors",
		Suite:     "Warm-pool errors",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("missing slotId → 400", func() {
			status, body, err := cli.Post("/pool/release", map[string]any{})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusBadRequest, body, "slotId is required")
		})
		a.Step("unknown slotId → 404", func() {
			status, body, err := cli.Post("/pool/release", map[string]any{"slotId": "no-such-slot"})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusNotFound, body, "slot not found")
		})
	})
}

func TestWarmPoolMethodNotAllowed(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "GET on POST-only pool routes returns 405",
		Package:   "tests.api.warmpool.WarmPoolMethodNotAllowedApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "HTTP methods",
		Suite:     "Warm-pool errors",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		for _, path := range []string{"/pool/reserve", "/pool/release", "/pool/preopen", "/pool/video/start", "/pool/video/stop"} {
			p := path
			a.Step("GET "+p, func() {
				status, _, err := cli.GetStatus(p)
				require.NoError(t, err)
				require.Equal(t, http.StatusMethodNotAllowed, status)
			})
		}
	})
}
