package hubapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestCreateSessionRequestJSON_IncludesDockerSafeChromeArgs(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "session body includes docker-safe chrome args",
		Package:   "config.CreateSessionRequestJsonTest",
		Layer:     "unit",
		Component: "selenoid",
		Suite:     "unit",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("marshal createSessionBody", func() {
			minVersion := config.MinVersion("chrome")
			raw, err := json.Marshal(hubapi.CreateSessionBody("chrome", minVersion))
			require.NoError(t, err)
			jsonStr := string(raw)
			require.Contains(t, jsonStr, "browserName")
			require.Contains(t, jsonStr, "chrome")
			require.Contains(t, jsonStr, "browserVersion")
			require.Contains(t, jsonStr, minVersion)
			require.Contains(t, jsonStr, "goog:chromeOptions")
			require.Contains(t, jsonStr, "no-sandbox")
			require.Contains(t, jsonStr, "disable-dev-shm-usage")
		})
	})
}

func TestWebDriverCreateSessionBody_FirefoxWarm(t *testing.T) {
	assertWarmSessionBody(t, "firefox", "moz:firefoxOptions", "-headless")
}

func TestWebDriverCreateSessionBody_MsEdgeWarm(t *testing.T) {
	assertWarmSessionBody(t, "msedge", "ms:edgeOptions", "no-sandbox")
}

func TestWebDriverCreateSessionBody_ChromeWarm(t *testing.T) {
	assertWarmSessionBody(t, "chrome", "goog:chromeOptions", "no-sandbox")
}

func TestWebDriverCreateSessionBody_ChromeMin(t *testing.T) {
	assertMinSessionBody(t, "chrome", "goog:chromeOptions", "no-sandbox")
}

func TestWebDriverCreateSessionBody_FirefoxMin(t *testing.T) {
	assertMinSessionBody(t, "firefox", "moz:firefoxOptions", "-headless")
}

func TestWebDriverCreateSessionBody_MsEdgeMin(t *testing.T) {
	assertMinSessionBody(t, "msedge", "ms:edgeOptions", "no-sandbox")
}

func assertWarmSessionBody(t *testing.T, browser, optionsKey, argFragment string) {
	t.Helper()
	allurex.Run(t, allurex.Meta{
		Name:      browser + " warm session body matches catalog default version",
		Package:   "config.WebDriverCreateSessionBodyTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Suite:     "unit",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("assert body", func() {
			warm := config.DefaultVersion(browser)
			raw, err := json.Marshal(hubapi.CreateSessionBody(browser, warm))
			require.NoError(t, err)
			assertSessionBody(t, string(raw), browser, warm, optionsKey, argFragment)
		})
	})
}

func assertMinSessionBody(t *testing.T, browser, optionsKey, argFragment string) {
	t.Helper()
	allurex.Run(t, allurex.Meta{
		Name:      browser + " min session body includes docker-safe args",
		Package:   "config.WebDriverCreateSessionBodyTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Suite:     "unit",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("assert body", func() {
			minVersion := config.MinVersion(browser)
			raw, err := json.Marshal(hubapi.CreateSessionBody(browser, minVersion))
			require.NoError(t, err)
			assertSessionBody(t, string(raw), browser, minVersion, optionsKey, argFragment)
			image, _ := config.VersionBlock(browser, minVersion)["image"].(string)
			require.Contains(t, image, config.MinImageMajor(browser)+"-min")
		})
	})
}

func assertSessionBody(t *testing.T, jsonStr, browserName, browserVersion, optionsKey, argFragment string) {
	t.Helper()
	require.Contains(t, jsonStr, "capabilities")
	require.Contains(t, jsonStr, "alwaysMatch")
	require.Contains(t, jsonStr, "browserName")
	require.Contains(t, jsonStr, browserName)
	require.Contains(t, jsonStr, "browserVersion")
	require.Contains(t, jsonStr, browserVersion)
	require.Contains(t, jsonStr, optionsKey)
	require.Contains(t, jsonStr, argFragment)

	var body map[string]any
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &body))
	caps := body["capabilities"].(map[string]any)
	always := caps["alwaysMatch"].(map[string]any)
	require.Equal(t, browserName, always["browserName"])
	require.Equal(t, browserVersion, always["browserVersion"])
}
