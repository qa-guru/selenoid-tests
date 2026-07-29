package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubVideo_ListJsonReturnsPaginatedPage(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /video/?json returns paginated page (default limit 10)",
		Package:   "tests.api.HubVideoApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "Hub video",
		Story:     "Hub video API",
		Suite:     "Hub video API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var listed *hubapi.VideoListResponse
		a.Step("GET /video/?json", func() {
			var err error
			listed, err = hubapi.ListVideoJSON(cfg, 10, 0, "")
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

func TestHubVideo_ListJsonSecondPage(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /video/?json&offset=10 returns second page metadata",
		Package:   "tests.api.HubVideoApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "Hub video",
		Story:     "Hub video API",
		Suite:     "Hub video API",
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var first, second *hubapi.VideoListResponse
		a.Step("GET first page", func() {
			var err error
			first, err = hubapi.ListVideoJSON(cfg, 10, 0, "")
			require.NoError(t, err)
		})
		a.Step("GET second page", func() {
			var err error
			second, err = hubapi.ListVideoJSON(cfg, 10, 10, "")
			require.NoError(t, err)
		})
		a.Step("Verify second page contract", func() {
			require.Equal(t, 10, second.Limit)
			require.Equal(t, 10, second.Offset)
			require.Equal(t, first.Total, second.Total)
			require.LessOrEqual(t, len(second.Videos), 10)
			if first.Total <= 10 {
				require.Empty(t, second.Videos)
			}
		})
	})
}

func TestHubVideo_MissingFileReturns404(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /video/{missing}.mp4 returns 404",
		Package:   "tests.api.HubVideoApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "Hub video",
		Story:     "Hub video API",
		Suite:     "Hub video API",
		Tags:      []string{"api", "negative"},
	}, func(a *allurex.A) {
		a.Step("GET missing video file", func() {
			require.NoError(t, hubapi.GetVideoExpectStatus(cfg, "missing-session-id.mp4", http.StatusNotFound))
		})
	})
}

func TestHubVideoSession_VideoListedAfterClose(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "Session with enableVideo=true produces downloadable MP4 after delete",
		Package:   "tests.api.HubVideoSessionApiTests",
		Layer:     "api",
		Component: "video-recorder",
		Epic:      "video-recorder",
		Feature:   "Hub session video",
		Story:     "Hub session video",
		Suite:     "Hub session video API",
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
		a.Step("Wait for video artifact in GET /video/?json&q=sessionId", func() {
			deadline := time.Now().Add(30 * time.Second)
			for time.Now().Before(deadline) {
				match, err := hubapi.FindVideoBySessionID(cfg, sessionID)
				require.NoError(t, err)
				if match != "" {
					videoFile = match
					return
				}
				time.Sleep(time.Second)
			}
			match, err := hubapi.FindVideoBySessionID(cfg, sessionID)
			require.NoError(t, err)
			videoFile = match
		})
		a.Step("Verify session video is listed", func() {
			require.NotEmpty(t, videoFile, "expected video for session %s in hub video list", sessionID)
		})
		var body []byte
		a.Step("Download session video", func() {
			var err error
			body, err = hubapi.DownloadVideo(cfg, videoFile)
			require.NoError(t, err)
		})
		a.Step("Verify MP4 payload", func() {
			assertValidMp4(t, body, videoFile)
		})
	})
}
