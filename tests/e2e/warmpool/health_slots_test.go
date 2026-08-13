package warmpool_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestWarmPoolHealthAndRoot_LiveStand(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "GET / and /health return ok plus slot count",
		Package:   "tests.e2e.warmpool.WarmPoolHealthTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool orchestrator",
		Story:     "Health",
		Suite:     "Warm-pool health",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("GET /health", func() {
			h, err := cli.Health()
			require.NoError(t, err)
			require.True(t, h.OK)
			require.GreaterOrEqual(t, h.Slots, 1)
		})
		a.Step("GET / (stand URL gate)", func() {
			h, err := cli.Root()
			require.NoError(t, err)
			require.True(t, h.OK)
		})
	})
}

func TestWarmPoolSlots_ListsChromeLoopback(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "GET /pool/slots lists chrome webdriver loopback URLs",
		Package:   "tests.e2e.warmpool.WarmPoolSlotsTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool orchestrator",
		Story:     "Slots",
		Suite:     "Warm-pool slots",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("GET /pool/slots", func() {
			slots, err := cli.Slots()
			require.NoError(t, err)
			require.NotEmpty(t, slots)
			var chrome int
			for _, s := range slots {
				if s.Protocol == "webdriver" && s.Browser == "chrome" {
					chrome++
					require.NotEmpty(t, s.ID)
					require.True(t, s.IsLoopback(), "slot %s DialURL=%s must be loopback for host hub", s.ID, s.DialURL())
				}
			}
			require.GreaterOrEqual(t, chrome, 1)
		})
	})
}
