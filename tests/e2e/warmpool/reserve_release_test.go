package warmpool_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestWarmPoolReserveRelease_LoopbackChrome(t *testing.T) {
	cli := liveClient(t)
	t.Cleanup(func() { releaseOwned(t, cli, testOwner) })

	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve loopback chrome then release",
		Package:   "tests.e2e.warmpool.WarmPoolReserveTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool orchestrator",
		Story:     "Reserve and release",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		var slotID string
		a.Step("POST /pool/reserve loopback chrome", func() {
			slot, status, err := cli.Reserve("webdriver", "chrome", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status)
			require.NotEmpty(t, slot.ID)
			require.True(t, slot.IsLoopback(), "webdriverUrl=%s", slot.DialURL())
			require.NotNil(t, slot.ReservedBy)
			require.Equal(t, testOwner, *slot.ReservedBy)
			slotID = slot.ID
		})
		a.Step("POST /pool/release", func() {
			require.NoError(t, cli.Release(slotID))
		})
		a.Step("slot is free again", func() {
			slots, err := cli.Slots()
			require.NoError(t, err)
			for _, s := range slots {
				if s.ID == slotID {
					require.Nil(t, s.ReservedBy)
				}
			}
		})
	})
}

func TestWarmPoolReserve_ExhaustedReturns409(t *testing.T) {
	cli := liveClient(t)
	t.Cleanup(func() { releaseOwned(t, cli, testOwner) })

	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve returns 409 when no loopback chrome slots left",
		Package:   "tests.e2e.warmpool.WarmPoolReserveExhaustedTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool orchestrator",
		Story:     "Reserve exhausted",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		var held []string
		a.Step("reserve all chrome loopback slots", func() {
			for i := 0; i < 8; i++ {
				slot, status, err := cli.Reserve("webdriver", "chrome", testOwner, true)
				require.NoError(t, err)
				if status == http.StatusConflict {
					break
				}
				require.Equal(t, http.StatusOK, status)
				held = append(held, slot.ID)
			}
			require.NotEmpty(t, held)
		})
		a.Step("next reserve is 409", func() {
			_, status, err := cli.Reserve("webdriver", "chrome", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, status)
		})
		a.Step("release held slots", func() {
			for _, id := range held {
				require.NoError(t, cli.Release(id))
			}
		})
	})
}
