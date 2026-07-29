package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiBrowsersConfig_ReturnsCatalog(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /browsers-config returns browser catalog JSON",
		Package:   "tests.api.UiBrowsersConfigTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI browsers config",
		Story:     "UI browsers config API",
		Suite:     "UI browsers config API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var catalog uiapi.BrowsersConfig
		a.Step("GET /browsers-config", func() {
			var err error
			catalog, err = uiapi.FetchBrowsersConfig(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify catalog map is present", func() {
			require.NotNil(t, catalog)
		})
	})
}

func TestUiStatusTotal_ExposesTotalCounter(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET UI /status exposes total quota counter",
		Package:   "tests.api.UiStatusTotalTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status proxy",
		Story:     "UI status total counter",
		Suite:     "UI status total counter",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET UI /status", func() {
			var err error
			status, err = uiapi.FetchStatus(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify total counter is non-negative", func() {
			require.GreaterOrEqual(t, status.Total, 0)
		})
	})
}

func TestUiPingVersion_ContainsRevisionMarker(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /ping version contains git revision marker",
		Package:   "tests.api.UiPingVersionTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI ping",
		Story:     "UI ping version",
		Suite:     "UI ping version",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var ping *uiapi.Ping
		a.Step("GET /ping", func() {
			var err error
			ping, err = uiapi.FetchPing(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify version looks like revision", func() {
			require.GreaterOrEqual(t, len(ping.Version), 3)
		})
	})
}

func TestUiSseStream_DeliversPayload(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "SSE /events delivers state or error payload",
		Package:   "tests.api.UiSseStreamTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI SSE",
		Story:     "UI SSE stream",
		Suite:     "UI SSE stream",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var event *uiapi.SseEvent
		a.Step("Read first SSE event", func() {
			var err error
			event, err = uiapi.ReadFirstSSEEvent(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify payload has state or errors", func() {
			require.True(t, event.HasState() || event.HasErrors())
			if event.HasState() {
				require.GreaterOrEqual(t, event.State.Total, 0)
				require.NotNil(t, event.State.Browsers)
			}
		})
	})
}

func TestUiSseMultipleEvents_DeliversTwoEvents(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "SSE /events delivers two consecutive payloads",
		Package:   "tests.api.UiSseMultipleEventsTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI SSE",
		Story:     "UI SSE multiple events",
		Suite:     "UI SSE multiple events",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var events []*uiapi.SseEvent
		a.Step("Read two SSE events", func() {
			var err error
			events, err = uiapi.ReadTwoSSEEvents(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify both events are parseable", func() {
			require.True(t, events[0].HasState() || events[0].HasErrors())
			require.True(t, events[1].HasState() || events[1].HasErrors())
		})
	})
}

func TestUiClipboard_UnknownSessionReturns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /clipboard/{sessionId} via UI proxy returns 404 for unknown session",
		Package:   "tests.api.UiClipboardApiTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI clipboard proxy",
		Story:     "UI /clipboard proxy API",
		Suite:     "UI /clipboard proxy API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET UI /clipboard/unknown-session", func() {
			require.NoError(t, uiapi.GetProxyExpectStatus(cfg, "/clipboard", "unknown-session", http.StatusNotFound))
		})
	})
}

func TestUiLogsWs_RequiresWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /ws/logs/{sessionId} without WebSocket upgrade returns 400 via UI proxy",
		Package:   "tests.api.UiLogsWsApiTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI logs WebSocket proxy",
		Story:     "UI logs WebSocket proxy",
		Suite:     "UI /ws/logs WebSocket API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session", func() {
			var err error
			sessionID, err = hubapi.CreateSession(cfg)
			require.NoError(t, err)
		})
		a.Step("GET /ws/logs/{sessionId} without WebSocket headers", func() {
			require.NoError(t, getExpectStatusBodyContains(uiHTTP(cfg), "/ws/logs/"+sessionID, http.StatusBadRequest, "websocket"))
		})
		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
	})
}

func TestUiVncWs_RequiresWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /ws/vnc/{sessionId} without WebSocket upgrade returns 400 via UI proxy",
		Package:   "tests.api.UiVncWsApiTests",
		Layer:     "api",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI VNC WebSocket proxy",
		Story:     "UI /ws/vnc WebSocket API",
		Suite:     "UI /ws/vnc WebSocket API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		a.Step("GET /ws/vnc/{sessionId} without WebSocket headers", func() {
			require.NoError(t, getExpectStatusBodyContains(uiHTTP(cfg), "/ws/vnc/abc", http.StatusBadRequest, "websocket"))
		})
	})
}

func TestUiVideo_ListJsonReturnsPaginatedPage(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /video/?json via UI proxy returns paginated page (default limit 10)",
		Package:   "tests.api.UiVideoApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "UI video proxy",
		Story:     "UI /video proxy API",
		Suite:     "UI /video proxy API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var listed *hubapi.VideoListResponse
		a.Step("GET UI /video/?json", func() {
			var err error
			listed, err = uiapi.ListVideoJSON(cfg, 10, 0, "")
			require.NoError(t, err)
		})
		a.Step("Verify paginated payload", func() {
			require.NotNil(t, listed)
			require.NotNil(t, listed.Videos)
			require.Equal(t, 10, listed.Limit)
			require.Equal(t, 0, listed.Offset)
			require.LessOrEqual(t, len(listed.Videos), listed.Limit)
			require.GreaterOrEqual(t, listed.Total, len(listed.Videos))
		})
	})
}

func TestUiVideo_MissingFileReturns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /video/{missing}.mp4 via UI proxy returns 404",
		Package:   "tests.api.UiVideoApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "UI video proxy",
		Story:     "UI /video proxy API",
		Suite:     "UI /video proxy API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET missing video via UI", func() {
			require.NoError(t, uiapi.GetVideoExpectStatus(cfg, "missing-session-id.mp4", http.StatusNotFound))
		})
	})
}

func TestUiVideoSession_VideoListedAfterClose(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Session video is downloadable via UI /video proxy after delete",
		Package:   "tests.api.UiVideoSessionApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "UI session video proxy",
		Story:     "UI session video proxy",
		Suite:     "UI session video via /video proxy",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var sessionID string
		a.Step("Create hub session with video", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithSelenoidOptions(cfg, "chrome", cfg.ChromeVersionForSession(), map[string]any{"enableVideo": true})
			require.NoError(t, err)
		})
		a.Step("Keep session briefly for recorder", func() {
			time.Sleep(3 * time.Second)
		})
		a.Step("Delete hub session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})
		var videoFile string
		a.Step("Wait for video artifact in UI GET /video/?json&q=sessionId", func() {
			deadline := time.Now().Add(30 * time.Second)
			for time.Now().Before(deadline) {
				match, err := uiapi.FindVideoBySessionID(cfg, sessionID)
				require.NoError(t, err)
				if match != "" {
					videoFile = match
					return
				}
				time.Sleep(time.Second)
			}
			match, err := uiapi.FindVideoBySessionID(cfg, sessionID)
			require.NoError(t, err)
			videoFile = match
		})
		a.Step("Verify session video is listed via UI proxy", func() {
			require.NotEmpty(t, videoFile)
		})
		var body []byte
		a.Step("Download session video via UI proxy", func() {
			var err error
			body, err = uiapi.DownloadVideo(cfg, videoFile)
			require.NoError(t, err)
		})
		a.Step("Verify MP4 payload via UI proxy", func() {
			assertValidMp4(t, body, videoFile)
		})
	})
}

func TestPlaywrightEndpoint_RequiresWebSocketUpgrade(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET playwright path without upgrade returns 400",
		Package:   "tests.api.PlaywrightEndpointTests",
		Layer:     "api",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright endpoint",
		Story:     "Playwright HTTP endpoint",
		Suite:     "Playwright HTTP endpoint",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		a.Step("GET playwright path without WebSocket headers", func() {
			require.NoError(t, playwrightapi.AssertUpgradeRequired(cfg))
		})
	})
}

func TestPlaywrightUnknownPath_Rejected(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET unknown playwright path returns 400",
		Package:   "tests.api.PlaywrightUnknownPathTests",
		Layer:     "api",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright endpoint",
		Story:     "Playwright unknown path",
		Suite:     "Playwright unknown path",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET unknown playwright path", func() {
			require.NoError(t, playwrightapi.AssertUnknownPathRejected(cfg))
		})
	})
}
