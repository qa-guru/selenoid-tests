package har_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/helpers"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestMain(m *testing.M) {
	if err := playwright.Install(&playwright.RunOptions{SkipInstallBrowsers: true}); err != nil {
		fmt.Fprintf(os.Stderr, "playwright driver install: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

const (
	harProdOutDir     = "build/har-prod-e2e"
	harCompareOutDir  = "build/har-compare"
	harProdStep5Dir   = "build/har-compare/prod-step5"
	bodiesMinContent  = 1
)

func hubHarMeta(name string, tags ...string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   "tests.HubHarProdE2eTests",
		Layer:     "e2e",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "HAR",
		Story:     "Hub enableHAR produces downloadable /har artifact",
		Suite:     "Hub HAR prod e2e",
		Tags:      append([]string{"smoke", "hub-har", "positive"}, tags...),
	}
}

func harCaptureMeta(name string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   "tests.HarCaptureProdE2eTests",
		Layer:     "e2e",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "HAR",
		Story:     "HarCapture BODIES on prod Selenoid",
		Suite:     "HarCapture prod e2e",
		Tags:      []string{"smoke", "har-capture", "positive"},
	}
}

func stripTrailingSlash(raw string) string {
	if raw == "" {
		return raw
	}
	if strings.HasSuffix(raw, "/") {
		return raw[:len(raw)-1]
	}
	return raw
}

func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Host
}

func assertValidHubHar(t *testing.T, body []byte, label, targetHost string) {
	t.Helper()
	var root map[string]any
	require.NoError(t, json.Unmarshal(body, &root))

	log, ok := root["log"].(map[string]any)
	require.True(t, ok, "HAR must contain log")
	require.Equal(t, "1.2", fmt.Sprint(log["version"]), "HAR log.version")

	entries, ok := log["entries"].([]any)
	require.True(t, ok, "HAR must contain log.entries")
	require.NotEmpty(t, entries, "HAR entries must not be empty")

	hasTarget := false
	hasMethod := false
	for _, raw := range entries {
		entry, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		req, _ := entry["request"].(map[string]any)
		if req == nil {
			req = map[string]any{}
		}
		entryURL := fmt.Sprint(req["url"])
		if strings.Contains(entryURL, targetHost) {
			hasTarget = true
		}
		if method, ok := req["method"].(string); ok && strings.TrimSpace(method) != "" {
			hasMethod = true
		}
		if timeVal, ok := entry["time"].(float64); ok {
			require.GreaterOrEqual(t, timeVal, 0.0, "entry.time must be >= 0 when present")
		}
	}
	require.True(t, hasTarget, "expected ≥1 request.url containing host %s", targetHost)
	require.True(t, hasMethod, "expected ≥1 request.method")

	stats := helpers.HarStatsFromBytes(label, body)
	require.Greater(t, stats.Entries, 0, "HarStats.entries must be > 0")
	require.Greater(t, stats.HTTPEntries, 0, "HarStats.httpEntries must be > 0")
}

func playwrightWsWithEnableHAR(cfg *config.Config) (string, error) {
	base, err := cfg.ResolvePlaywrightWsEndpoint()
	if err != nil {
		return "", err
	}
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	return base + sep + "enableHAR=true&enableVideo=false&name=hub-har-pw-prod-e2e", nil
}

func writeHarArtifact(t *testing.T, dir, name string, body []byte) string {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, body, 0o644))
	return path
}

func TestHubHarProd_WebDriverEnableHarProducesValidArtifact(t *testing.T) {
	cfg := config.MustLoad()
	targetURL := stripTrailingSlash(cfg.SmokeURL)
	targetHost := hostOf(targetURL)

	allurex.Run(t, hubHarMeta("WebDriver enableHAR → /har is valid HAR 1.2 with target URL"), func(a *allurex.A) {
		var sessionID string
		a.Step("Create WebDriver session with enableHAR", func() {
			var err error
			sessionID, err = hubapi.CreateSessionWithSelenoidOptions(cfg, "chrome", cfg.ChromeVersionForSession(), map[string]any{
				"enableVNC":      false,
				"enableVideo":    false,
				"enableHAR":      true,
				"enableLog":      false,
				"sessionTimeout": "2m",
				"name":           "hub-har-wd-prod-e2e",
			})
			require.NoError(t, err)
		})

		a.Step("Navigate to "+targetURL, func() {
			require.NoError(t, hubapi.NavigateSession(cfg, sessionID, targetURL))
			time.Sleep(2 * time.Second)
		})

		a.Step("Delete WebDriver session", func() {
			require.NoError(t, hubapi.DeleteSession(cfg, sessionID))
		})

		var body []byte
		a.Step("Wait and download /har for "+sessionID, func() {
			file, err := hubapi.WaitForSessionHar(cfg, sessionID, 30*time.Second)
			require.NoError(t, err)
			require.NotEmpty(t, file, "expected hub HAR for WD session %s", sessionID)
			body, err = hubapi.DownloadHar(cfg, file)
			require.NoError(t, err)
		})

		a.Step("Assert HAR 1.2 + entries + target host", func() {
			writeHarArtifact(t, harProdOutDir, "wd-"+sessionID+".har", body)
			assertValidHubHar(t, body, "wd-hub-enableHAR", targetHost)
		})
	})
}

func TestHubHarProd_PlaywrightEnableHarProducesValidArtifact(t *testing.T) {
	cfg := config.MustLoad()
	targetURL := stripTrailingSlash(cfg.SmokeURL)
	targetHost := hostOf(targetURL)

	allurex.Run(t, hubHarMeta("Playwright enableHAR → /har is valid HAR 1.2 with target URL"), func(a *allurex.A) {
		var sessionID string
		a.Step("Connect Playwright with enableHAR and navigate", func() {
			before, err := hubapi.CollectSessionIDs(cfg)
			require.NoError(t, err)

			ws, err := playwrightWsWithEnableHAR(cfg)
			require.NoError(t, err)

			pw, err := playwright.Run()
			require.NoError(t, err)
			defer func() { require.NoError(t, pw.Stop()) }()

			browser, err := playwrightapi.Connect(pw, cfg, ws)
			require.NoError(t, err)

			id, err := hubapi.WaitNewSessionID(cfg, before, 30*time.Second)
			require.NoError(t, err)
			require.NotEmpty(t, id)
			sessionID = id

			ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{
				ServiceWorkers: playwright.ServiceWorkerPolicyBlock,
			})
			require.NoError(t, err)
			page, err := ctx.NewPage()
			require.NoError(t, err)

			time.Sleep(750 * time.Millisecond)
			_, err = page.Goto(targetURL)
			require.NoError(t, err)
			require.NoError(t, page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
				State: playwright.LoadStateLoad,
			}))
			time.Sleep(1500 * time.Millisecond)

			_ = ctx.Close()
			_ = playwrightapi.Close(browser)
		})

		require.NotEmpty(t, sessionID, "Playwright hub session id not observed in /status")

		var body []byte
		a.Step("Wait and download /har for "+sessionID, func() {
			file, err := hubapi.WaitForSessionHar(cfg, sessionID, 30*time.Second)
			require.NoError(t, err)
			require.NotEmpty(t, file, "expected hub HAR for PW session %s — Chromium-family image with DevTools :7070", sessionID)
			body, err = hubapi.DownloadHar(cfg, file)
			require.NoError(t, err)
		})

		a.Step("Assert HAR 1.2 + entries + target host", func() {
			writeHarArtifact(t, harProdOutDir, "pw-"+sessionID+".har", body)
			assertValidHubHar(t, body, "pw-hub-enableHAR", targetHost)
		})
	})
}

func TestHarCaptureProd_BodiesProducesTextOnProd(t *testing.T) {
	cfg := config.MustLoad()
	targetURL := stripTrailingSlash(cfg.SmokeURL)
	label := "3b-selenide-HarCapture-bodies-prod"

	allurex.Run(t, harCaptureMeta("Selenide HarCapture BODIES → valid HAR with content.text on prod"), func(a *allurex.A) {
		var stats helpers.HarStats
		a.Step("Capture HarCapture BODIES on "+targetURL, func() {
			var sessionID string
			var err error
			sessionID, err = hubapi.CreateHarCaptureSession(cfg, label)
			require.NoError(t, err)
			defer func() { _ = hubapi.DeleteSession(cfg, sessionID) }()

			require.NoError(t, hubapi.NavigateSession(cfg, sessionID, targetURL))
			time.Sleep(2 * time.Second)

			harBytes, err := hubapi.CollectHarFromSession(cfg, sessionID, helpers.HarBodies)
			require.NoError(t, err)
			require.NotEmpty(t, harBytes)

			require.NoError(t, os.MkdirAll(harCompareOutDir, 0o755))
			require.NoError(t, os.MkdirAll(harProdStep5Dir, 0o755))

			harPath := filepath.Join(harCompareOutDir, "3b-selenide-HarCapture-bodies.har")
			prodHarPath := filepath.Join(harCompareOutDir, "3b-selenide-HarCapture-bodies-prod.har")
			step5HarPath := filepath.Join(harProdStep5Dir, "3b-selenide-HarCapture-bodies-prod.har")

			require.NoError(t, os.WriteFile(harPath, harBytes, 0o644))
			require.NoError(t, os.WriteFile(prodHarPath, harBytes, 0o644))
			require.NoError(t, os.WriteFile(step5HarPath, harBytes, 0o644))

			stats = helpers.HarStatsFromBytes(label, harBytes)
			if stats.WithContentSize < 1 {
				t.Logf("NOTE: HarCapture prod bodies withContentSize=%d (best-effort gate; text gate is authoritative)", stats.WithContentSize)
			}
		})

		a.Step("Assert HarCapture prod bodies gates", func() {
			require.GreaterOrEqual(t, stats.HTTPEntries, 1, "expected ≥1 http entry on prod smoke")
			require.GreaterOrEqual(t, stats.WithContentText, bodiesMinContent,
				"HarCapture BODIES prod expected withContentText>=%d, got %d",
				bodiesMinContent, stats.WithContentText)
		})
	})
}
