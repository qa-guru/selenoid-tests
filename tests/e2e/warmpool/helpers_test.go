package warmpool_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/warmpool"
)

const testOwner = "selenoid-tests-e2e"

func liveClient(t *testing.T) *warmpool.Client {
	t.Helper()
	cli := warmpool.Default()
	if err := cli.Ping(); err != nil {
		t.Skipf("warm-pool stand down (%s): %v", warmpool.BaseURL(), err)
	}
	return cli
}

func releaseOwned(t *testing.T, cli *warmpool.Client, prefix string) {
	t.Helper()
	slots, err := cli.Slots()
	if err != nil {
		return
	}
	for _, s := range slots {
		if s.ReservedBy != nil && *s.ReservedBy == prefix {
			_ = cli.Release(s.ID)
		}
	}
}

func hasDialableLoopback(t *testing.T, cli *warmpool.Client) bool {
	t.Helper()
	slots, err := cli.Slots()
	if err != nil {
		return false
	}
	httpCli := &http.Client{Timeout: 800 * time.Millisecond}
	for _, s := range slots {
		if s.Protocol != "webdriver" || s.Browser != "chrome" || !s.IsLoopback() || s.ReservedBy != nil {
			continue
		}
		u := s.DialURL()
		if u == "" {
			continue
		}
		req, err := http.NewRequest(http.MethodHead, u, nil)
		if err != nil {
			continue
		}
		resp, err := httpCli.Do(req)
		if err != nil {
			continue
		}
		_ = resp.Body.Close()
		return true
	}
	return false
}

func hubWarmTotal(t *testing.T) int {
	t.Helper()
	cfg := config.MustLoad()
	st, err := hubapi.Fetch(cfg)
	if err != nil {
		return 0
	}
	return st.WarmTotal
}
