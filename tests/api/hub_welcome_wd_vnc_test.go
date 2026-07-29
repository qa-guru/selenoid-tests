package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubWelcome_ReturnsSelenoidBanner(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET / returns Selenoid welcome text",
		Package:   "tests.api.HubWelcomeApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub welcome",
		Story:     "Hub welcome page API",
		Suite:     "Hub welcome page API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var body string
		a.Step("GET /", func() {
			var err error
			body, err = hubapi.FetchWelcomeText(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify welcome text", func() {
			require.True(t, strings.Contains(body, "Selenoid"))
		})
	})
}

func TestHubWebDriverStatus_ReportsReady(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /wd/hub/status reports ready hub",
		Package:   "tests.api.HubWebDriverStatusApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub WebDriver status",
		Story:     "Hub WebDriver status API",
		Suite:     "Hub WebDriver status API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.WebDriverStatus
		a.Step("GET /wd/hub/status", func() {
			var err error
			status, err = hubapi.FetchWebDriverStatus(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify ready flag and message", func() {
			require.True(t, status.Value.Ready)
			require.NotEmpty(t, status.Value.Message)
		})
	})
}

func TestWebDriverStatus_ReportsReadyForWebDriverImage(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /wd/hub/status reports ready hub for webdriver-image sessions",
		Package:   "tests.api.WebDriverStatusApiTests",
		Layer:     "api",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver hub status",
		Story:     "WebDriver image hub status API",
		Suite:     "WebDriver image hub status API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.WebDriverStatus
		a.Step("GET /wd/hub/status", func() {
			var err error
			status, err = hubapi.FetchWebDriverStatus(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify ready flag and message", func() {
			require.True(t, status.Value.Ready)
			require.NotEmpty(t, status.Value.Message)
		})
	})
}

func TestWebDriverSession_CreateAndDeleteChromeSession(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "POST /wd/hub/session creates chrome session via webdriver-image node",
		Package:   "tests.api.WebDriverSessionApiTests",
		Layer:     "api",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Feature:   "WebDriver session API",
		Story:     "WebDriver session API",
		Suite:     "WebDriver image session API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var created *hubapi.SessionCreateResult
		a.Step("Create remote session", func() {
			var err error
			created, err = hubapi.CreateSessionWithCapabilities(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify session id and browserName", func() {
			require.NotEmpty(t, created.SessionID)
			require.Equal(t, cfg.Browser, created.BrowserName)
		})
		a.Step("Delete session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, created.SessionID))
		})
	})
}

func TestHubVncSession_PathRequiresWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /vnc/{sessionId} without WebSocket upgrade returns 400 when enableVNC=true",
		Package:   "tests.api.HubVncSessionApiTests",
		Layer:     "api",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub VNC WebSocket",
		Story:     "Hub VNC WebSocket",
		Suite:     "Hub session VNC WebSocket API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session with VNC", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithSelenoidOptions(cfg, cfg.Browser, cfg.BrowserVersion, map[string]any{"enableVNC": true})
			require.NoError(t, err)
		})
		a.Step("GET /vnc/{sessionId} without WebSocket headers", func() {
			require.NoError(t, getExpectStatusBodyContains(hubHTTP(cfg), "/vnc/"+sessionID, http.StatusBadRequest, "websocket"))
		})
		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
	})
}
