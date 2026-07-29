package cmapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func apiMeta(name, pkg, feature, story string, tags ...string) allurex.Meta {
	tagList := append([]string{"api", "cm", "positive"}, tags...)
	return allurex.Meta{
		Name:      name,
		Package:   pkg,
		Layer:     "api",
		Component: "cm",
		Epic:      "cm",
		Feature:   feature,
		Story:     story,
		Suite:     story,
		Tags:      tagList,
	}
}

func TestCmHubStatus_ReturnsStatisticsJSON(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, apiMeta(
		"GET /status on CM hub port returns statistics JSON",
		"tests.api.CmHubStatusApiTests",
		"CM-managed hub",
		"CM hub status API",
	), func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET CM hub /status", func() {
			var err error
			status, err = hubapi.FetchFrom(cfg.ResolveCmHubURL())
			require.NoError(t, err)
		})
		a.Step("Verify hub counters are non-negative", func() {
			require.GreaterOrEqual(t, status.Total, 0)
			require.GreaterOrEqual(t, status.Used, 0)
			require.NotNil(t, status.Browsers)
		})
	})
}

func TestCmHubSession_CreateAndDeleteSession(t *testing.T) {
	cfg := config.MustLoad()
	skipUnlessCmSessionReady(t, cfg)
	allurex.Run(t, apiMeta(
		"POST /wd/hub/session on CM hub creates and DELETE removes session",
		"tests.api.CmHubSessionApiTests",
		"CM-managed hub session",
		"CM hub session API",
	), func(a *allurex.A) {
		var sessionID string
		a.Step("Create session on CM hub", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithBrowser(cfg, cfg.Browser, cfg.BrowserVersion)
			require.NoError(t, err)
		})
		a.Step("Verify session id", func() {
			require.NotEmpty(t, sessionID)
		})
		a.Step("Delete session on CM hub", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
	})
}

func TestCmUiStatus_ReturnsProxiedCounters(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, apiMeta(
		"GET /status on CM UI port returns proxied hub counters",
		"tests.api.CmUiStatusApiTests",
		"CM-managed UI",
		"CM UI status API",
	), func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET CM UI /status", func() {
			var err error
			status, err = uiapi.FetchStatusFrom(cfg.ResolveCmUiURL())
			require.NoError(t, err)
		})
		a.Step("Verify UI status counters are present", func() {
			require.GreaterOrEqual(t, status.Total, 0)
			require.GreaterOrEqual(t, status.Used, 0)
			require.NotNil(t, status.Browsers)
		})
	})
}
