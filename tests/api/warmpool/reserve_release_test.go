package warmpool_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/warmpool"
)

func TestWarmPoolReserveRelease_LoopbackChrome(t *testing.T) {
	cli := liveClient(t)
	t.Cleanup(func() { releaseOwned(t, cli, testOwner) })

	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve loopback chrome then release",
		Package:   "tests.api.warmpool.WarmPoolReserveApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
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
			require.Equal(t, "webdriver", slot.Protocol)
			require.Equal(t, "chrome", slot.Browser)
			require.True(t, slot.IsLoopback(), "webdriverUrl=%s", slot.DialURL())
			require.NotNil(t, slot.ReservedBy)
			require.Equal(t, testOwner, *slot.ReservedBy)
			slotID = slot.ID
		})
		a.Step("GET /pool/slots shows reservedBy", func() {
			slots, err := cli.Slots()
			require.NoError(t, err)
			found := false
			for _, s := range slots {
				if s.ID == slotID {
					found = true
					require.NotNil(t, s.ReservedBy)
					require.Equal(t, testOwner, *s.ReservedBy)
				}
			}
			require.True(t, found)
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
		a.Step("second release of a free slot is still 200", func() {
			require.NoError(t, cli.Release(slotID))
		})
	})
}

func TestWarmPoolReserve_ExhaustedReturns409(t *testing.T) {
	cli := liveClient(t)
	t.Cleanup(func() { releaseOwned(t, cli, testOwner) })

	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve returns 409 when no loopback chrome slots left",
		Package:   "tests.api.warmpool.WarmPoolReserveExhaustedApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Reserve exhausted",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		var held []string
		a.Step("reserve all chrome loopback slots", func() {
			for i := 0; i < 16; i++ {
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
		a.Step("next reserve is 409 no available slots", func() {
			_, status, err := cli.Reserve("webdriver", "chrome", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, status)
			_, body, err := cli.Post("/pool/reserve", map[string]any{
				"protocol": "webdriver",
				"browser":  "chrome",
				"owner":    testOwner,
				"loopback": true,
			})
			require.NoError(t, err)
			require.Equal(t, "no available slots", warmpool.ParseError(body))
		})
		a.Step("release held slots", func() {
			for _, id := range held {
				require.NoError(t, cli.Release(id))
			}
		})
	})
}

func TestWarmPoolReserve_UnknownBrowser409(t *testing.T) {
	cli := liveClient(t)
	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve firefox or playwright returns 409 on chrome-only pool",
		Package:   "tests.api.warmpool.WarmPoolReserveFilterApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Reserve filter",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "negative", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("browser=firefox → 409", func() {
			_, status, err := cli.Reserve("webdriver", "firefox", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, status)
		})
		a.Step("protocol=playwright → 409", func() {
			_, status, err := cli.Reserve("playwright", "chromium", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, status)
		})
	})
}

func TestWarmPoolReserve_EmptyFilterTakesAnyLoopback(t *testing.T) {
	cli := liveClient(t)
	t.Cleanup(func() { releaseOwned(t, cli, testOwner) })

	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve with empty protocol/browser still returns a loopback slot",
		Package:   "tests.api.warmpool.WarmPoolReserveEmptyFilterApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Reserve filter",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("reserve without protocol/browser", func() {
			slot, status, err := cli.Reserve("", "", testOwner, true)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status)
			require.NotEmpty(t, slot.ID)
			require.True(t, slot.IsLoopback())
			require.NoError(t, cli.Release(slot.ID))
		})
	})
}

func TestWarmPoolReserve_DefaultOwnerAnonymous(t *testing.T) {
	cli := liveClient(t)
	var slotID string
	t.Cleanup(func() {
		if slotID != "" {
			_ = cli.Release(slotID)
		}
	})
	allurex.Run(t, allurex.Meta{
		Name:      "POST /pool/reserve without owner sets reservedBy anonymous",
		Package:   "tests.api.warmpool.WarmPoolReserveAnonymousApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool API",
		Story:     "Reserve owner",
		Suite:     "Warm-pool reserve",
		Tags:      []string{"api", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		a.Step("reserve with empty owner", func() {
			status, body, err := cli.Post("/pool/reserve", map[string]any{
				"protocol": "webdriver",
				"browser":  "chrome",
				"loopback": true,
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, status)
			var out struct {
				OK   bool          `json:"ok"`
				Slot warmpool.Slot `json:"slot"`
			}
			require.NoError(t, json.Unmarshal(body, &out))
			require.True(t, out.OK)
			require.NotNil(t, out.Slot.ReservedBy)
			require.Equal(t, "anonymous", *out.Slot.ReservedBy)
			slotID = out.Slot.ID
		})
	})
}
