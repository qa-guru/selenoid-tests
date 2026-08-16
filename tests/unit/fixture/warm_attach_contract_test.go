package fixture_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

// Container-reuse contract: orchestrator reserve (loopback) JSON the hub Client decodes.
func TestWarmPoolReserveLoopbackContract(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "warm-pool reserve loopback JSON exposes hub-dialable webdriverUrl",
		Package:   "tests.unit.fixture.ContainerReuseContractTest",
		Layer:     "unit",
		Component: "selenoid",
		Feature:   "Warm-pool container-reuse contract",
		Suite:     "Warm-pool container-reuse contract",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/warm-pool/reserve-loopback.json", func() {
			raw := loadFixture(t, "warm-pool/reserve-loopback.json")
			var out struct {
				OK   bool `json:"ok"`
				Slot struct {
					ID                   string  `json:"id"`
					Protocol             string  `json:"protocol"`
					Browser              string  `json:"browser"`
					WebdriverURL         string  `json:"webdriverUrl"`
					WebdriverURLLoopback string  `json:"webdriverUrlLoopback"`
					ReservedBy           *string `json:"reservedBy"`
				} `json:"slot"`
			}
			require.NoError(t, json.Unmarshal(raw, &out))
			require.True(t, out.OK)
			require.Equal(t, "pool-chrome-1", out.Slot.ID)
			require.Equal(t, "webdriver", out.Slot.Protocol)
			require.Equal(t, "chrome", out.Slot.Browser)
			require.Equal(t, out.Slot.WebdriverURLLoopback, out.Slot.WebdriverURL)
			u, err := url.Parse(out.Slot.WebdriverURL)
			require.NoError(t, err)
			require.Equal(t, "127.0.0.1", u.Hostname())
			require.NotNil(t, out.Slot.ReservedBy)
			require.Equal(t, "hub-1", *out.Slot.ReservedBy)
		})
	})
}
