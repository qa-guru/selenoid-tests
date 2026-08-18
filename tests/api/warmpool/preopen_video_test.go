package warmpool_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/warmpool"
)

func TestWarmPoolPreopen_ValidationAndWarmAPI(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/preopen validates body and proxies to slot warm-api",
		Package:   "tests.api.warmpool.WarmPoolPreopenApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Preopen",
		Suite:     "Warm-pool preopen",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("missing url → 400", func() {
			status, body, err := cli.Post("/pool/preopen", map[string]any{"slotId": "pool-chrome-1"})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusBadRequest, body, "slotId and url are required")
		})
		a.Step("missing slotId → 400", func() {
			status, body, err := cli.Post("/pool/preopen", map[string]any{"url": config.DefaultSmokeURL})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusBadRequest, body, "slotId and url are required")
		})
		a.Step("unknown slot → 404", func() {
			status, body, err := cli.Post("/pool/preopen", map[string]any{
				"slotId": "no-such-slot",
				"url":    config.DefaultSmokeURL,
			})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusNotFound, body, "slot not found")
		})
		a.Step("known slot → 200 if warm-api up, else 500", func() {
			id := anySlotID(t, cli)
			status, body, err := cli.Post("/pool/preopen", map[string]any{
				"slotId": id,
				"url":    config.DefaultSmokeURL,
			})
			require.NoError(t, err)
			requireWarmProxy(t, status, body)
		})
	})
}

func TestWarmPoolVideo_ValidationAndWarmAPI(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/video/start and /stop validate body and proxy to warm-api",
		Package:   "tests.api.warmpool.WarmPoolVideoApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Video",
		Suite:     "Warm-pool video",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("video/start missing slotId → 400", func() {
			status, body, err := cli.Post("/pool/video/start", map[string]any{})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusBadRequest, body, "slotId is required")
		})
		a.Step("video/stop missing slotId → 400", func() {
			status, body, err := cli.Post("/pool/video/stop", map[string]any{})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusBadRequest, body, "slotId is required")
		})
		a.Step("video/start unknown slot → 404", func() {
			status, body, err := cli.Post("/pool/video/start", map[string]any{"slotId": "no-such-slot"})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusNotFound, body, "slot not found")
		})
		a.Step("video/stop unknown slot → 404", func() {
			status, body, err := cli.Post("/pool/video/stop", map[string]any{"slotId": "no-such-slot"})
			require.NoError(t, err)
			requireAPIError(t, status, http.StatusNotFound, body, "slot not found")
		})
		id := anySlotID(t, cli)
		a.Step("video/start known slot → 200 or 500", func() {
			status, body, err := cli.Post("/pool/video/start", map[string]any{
				"slotId":    id,
				"sessionId": "selenoid-tests-api",
			})
			require.NoError(t, err)
			requireWarmProxy(t, status, body)
		})
		a.Step("video/stop known slot → 200 or 500", func() {
			status, body, err := cli.Post("/pool/video/stop", map[string]any{"slotId": id})
			require.NoError(t, err)
			requireWarmProxy(t, status, body)
		})
	})
}

// requireWarmProxy accepts 200 when the slot warm-api answers, or 500 when it is down.
func requireWarmProxy(t *testing.T, status int, body []byte) {
	t.Helper()
	switch status {
	case http.StatusOK:
		require.Contains(t, string(body), `"ok":true`)
	case http.StatusInternalServerError:
		require.NotEmpty(t, warmpool.ParseError(body))
	default:
		t.Fatalf("warm-api proxy: want 200 or 500, got %d %s", status, body)
	}
}
