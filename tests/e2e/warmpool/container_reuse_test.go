package warmpool_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestContainerReuse_ReservesAndReleasesChromeSlot(t *testing.T) {
	cli := liveClient(t)
	if hubWarmTotal(t) == 0 {
		t.Skip("hub /status warmTotal=0 — start hub with SELENOID_WARM_POOL_URL=http://127.0.0.1:9090 (dev/scripts/start-selenoid.sh)")
	}
	if !hasDialableLoopback(t, cli) {
		t.Skip("no dialable loopback ChromeDriver — docker compose -f docker-compose.local.yml up -d in selenoid-pool")
	}

	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Hub Chrome session reuses a warm-pool container and releases it",
		Package:   "tests.e2e.warmpool.ContainerReuseTests",
		Layer:     "e2e",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Warm-pool container-reuse",
		Story:     "Reserve and release",
		Suite:     "Warm-pool container-reuse",
		Tags:      []string{"e2e", "positive", "warm-pool"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create Chrome WD session (no video/VNC/HAR)", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithBrowser(cfg, "chrome", cfg.ChromeVersionForSession())
			require.NoError(t, err)
			require.NotEmpty(t, sessionID)
		})
		defer func() {
			if sessionID != "" {
				_ = hubapi.DeleteSession(cfg, sessionID)
			}
		}()

		a.Step("A chrome slot is reserved by hub-*", func() {
			deadline := time.Now().Add(8 * time.Second)
			var reserved string
			for time.Now().Before(deadline) {
				slots, err := cli.Slots()
				require.NoError(t, err)
				for _, s := range slots {
					if s.ReservedBy != nil && *s.ReservedBy != "" && s.Browser == "chrome" {
						reserved = *s.ReservedBy
						break
					}
				}
				if reserved != "" {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
			require.NotEmpty(t, reserved, "expected hub to reserve a warm chrome slot")
			require.Contains(t, reserved, "hub-")
		})

		a.Step("Delete session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
			sessionID = ""
		})

		a.Step("Slot is released", func() {
			deadline := time.Now().Add(8 * time.Second)
			for time.Now().Before(deadline) {
				slots, err := cli.Slots()
				require.NoError(t, err)
				busy := false
				for _, s := range slots {
					if s.ReservedBy != nil && *s.ReservedBy != "" {
						busy = true
						break
					}
				}
				if !busy {
					return
				}
				time.Sleep(200 * time.Millisecond)
			}
			t.Fatal("expected all slots free after session delete")
		})
	})
}
