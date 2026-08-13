package warmpool_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestWarmPoolHealthAndRoot(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "GET / and /health return ok plus slot count",
		Package:   "tests.api.warmpool.WarmPoolHealthApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Health",
		Suite:     "Warm-pool health",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		var n int
		a.Step("GET /health", func() {
			h, err := cli.Health()
			require.NoError(t, err)
			require.True(t, h.OK)
			require.GreaterOrEqual(t, h.Slots, 1)
			n = h.Slots
		})
		a.Step("GET / matches health", func() {
			h, err := cli.Root()
			require.NoError(t, err)
			require.True(t, h.OK)
			require.Equal(t, n, h.Slots)
		})
	})
}

func TestWarmPoolSlots_ChromeLoopbackCatalog(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "GET /pool/slots lists chrome webdriver loopback URLs",
		Package:   "tests.api.warmpool.WarmPoolSlotsApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Slots",
		Suite:     "Warm-pool slots",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("GET /pool/slots matches /health count", func() {
			h, err := cli.Health()
			require.NoError(t, err)
			slots, err := cli.Slots()
			require.NoError(t, err)
			require.Len(t, slots, h.Slots)
			var chrome int
			for _, s := range slots {
				require.NotEmpty(t, s.ID)
				require.NotEmpty(t, s.Protocol)
				require.NotEmpty(t, s.Browser)
				if s.Protocol == "webdriver" && s.Browser == "chrome" {
					chrome++
					require.True(t, s.IsLoopback(), "slot %s DialURL=%s", s.ID, s.DialURL())
				}
			}
			require.GreaterOrEqual(t, chrome, 1)
		})
		a.Step("response is a JSON array", func() {
			status, body, err := cli.GetStatus("/pool/slots")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status)
			var raw []json.RawMessage
			require.NoError(t, json.Unmarshal(body, &raw))
			require.NotEmpty(t, raw)
		})
	})
}
