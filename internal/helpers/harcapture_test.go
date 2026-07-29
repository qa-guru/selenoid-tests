package helpers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/helpers"
)

func TestHarCapture_ToHarBuildsEntriesFromPerformanceLogs(t *testing.T) {
	allurex.Run(t, harMeta("toHar builds entries from performance logs"), func(a *allurex.A) {
		a.Step("toHar", func() {
			har := helpers.ToHar(fixtureEntries(), helpers.HarMeta, nil)
			require.Contains(t, har, "1.2")
			require.Contains(t, har, "example.com")
			require.Contains(t, har, "200")
			require.Contains(t, har, "1280")
			require.True(t, helpers.SupportsBrowser("chrome"))
			require.False(t, helpers.SupportsBrowser("firefox"))
		})
	})
}

func TestHarCapture_ToHarDefaultMetaOmitsContentText(t *testing.T) {
	allurex.Run(t, harMeta("toHar default/meta omits content.text"), func(a *allurex.A) {
		a.Step("compare modes", func() {
			entries := fixtureEntries()
			harDefault := helpers.ToHar(entries, helpers.HarMeta, nil)
			harMeta := helpers.ToHar(entries, helpers.HarMeta, nil)

			statsDefault := helpers.HarStatsFromBytes("unit-meta-default", []byte(harDefault))
			statsMeta := helpers.HarStatsFromBytes("unit-meta", []byte(harMeta))

			require.Equal(t, 0, statsDefault.WithContentText)
			require.Equal(t, 0, statsMeta.WithContentText)
			require.NotContains(t, harDefault, `"text"`)
			require.NotContains(t, harMeta, `"text"`)
		})
	})
}

func TestHarCapture_ToHarBodiesIncludesSyntheticContentText(t *testing.T) {
	allurex.Run(t, harMeta("toHar BODIES includes synthetic content.text"), func(a *allurex.A) {
		a.Step("bodies mode", func() {
			bodies := map[string]helpers.CapturedBody{
				"r1": {Text: "<html>synthetic-body</html>", Base64Encoded: false},
			}
			har := helpers.ToHar(fixtureEntries(), helpers.HarBodies, bodies)
			stats := helpers.HarStatsFromBytes("unit-bodies", []byte(har))
			require.Greater(t, stats.WithContentText, 0)
			require.Contains(t, har, "synthetic-body")
			require.NotContains(t, har, `"encoding"`)
		})
	})
}

func TestHarCapture_ToHarBodiesBase64SetsEncoding(t *testing.T) {
	allurex.Run(t, harMeta("toHar BODIES base64 sets encoding"), func(a *allurex.A) {
		a.Step("base64 body", func() {
			bodies := map[string]helpers.CapturedBody{
				"r1": {Text: "aGVsbG8=", Base64Encoded: true},
			}
			har := helpers.ToHar(fixtureEntries(), helpers.HarBodies, bodies)
			require.Contains(t, har, "aGVsbG8=")
			require.Contains(t, har, "base64")
		})
	})
}

func TestHarCapture_ToHarBodiesWithoutPayloadStaysMetaLike(t *testing.T) {
	allurex.Run(t, harMeta("toHar BODIES without payload stays meta-like"), func(a *allurex.A) {
		a.Step("empty bodies", func() {
			har := helpers.ToHar(fixtureEntries(), helpers.HarBodies, map[string]helpers.CapturedBody{})
			stats := helpers.HarStatsFromBytes("unit-bodies-empty", []byte(har))
			require.Equal(t, 0, stats.WithContentText)
		})
	})
}

func TestHarCapture_FinishedRequestIdsFromLoadingFinished(t *testing.T) {
	allurex.Run(t, harMeta("finishedRequestIds from loadingFinished"), func(a *allurex.A) {
		a.Step("ids", func() {
			require.Equal(t, []string{"r1"}, helpers.FinishedRequestIds(fixtureEntries()))
		})
	})
}

func harMeta(name string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   "helpers.HarCaptureTest",
		Layer:     "unit",
		Component: "selenoid",
		Suite:     "unit",
		Tags:      []string{"unit"},
	}
}

func fixtureEntries() []string {
	return []string{
		`{"message":{"method":"Network.requestWillBeSent","params":{"requestId":"r1","timestamp":1.0,"wallTime":1700000000.0,"request":{"url":"https://example.com/","method":"GET","headers":{"Accept":"*/*"}}}}}`,
		`{"message":{"method":"Network.responseReceived","params":{"requestId":"r1","response":{"status":200,"statusText":"OK","mimeType":"text/html","headers":{"content-type":"text/html"},"protocol":"http/1.1","encodedDataLength":42}}}}`,
		`{"message":{"method":"Network.loadingFinished","params":{"requestId":"r1","timestamp":1.05,"encodedDataLength":1280}}}`,
	}
}

type androidHarBlocker struct {
	Phase           int      `json:"phase"`
	Label           string   `json:"label"`
	Status          string   `json:"status"`
	Verdict         string   `json:"verdict"`
	Scorecard       any      `json:"scorecard"`
	HARPath         *string  `json:"harPath"`
	AllowedClaims   []string `json:"allowedClaims"`
	ForbiddenClaims []string `json:"forbiddenClaims"`
	Blockers        []struct {
		ID      string `json:"id"`
		Layer   string `json:"layer"`
		Summary string `json:"summary"`
	} `json:"blockers"`
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		require.NotEqual(t, parent, dir, "go.mod not found from %s", dir)
		dir = parent
	}
}

func TestHarBenchmarkAndroidNotClaimed(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "Phase 6 Android HAR remains not_claimed until CDP :7070 exists",
		Package:   "helpers.HarBenchmark",
		Layer:     "unit",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "HAR",
		Story:     "Android hub enableHAR is blocked by architecture",
		Suite:     "HAR benchmark SSOT",
		Tags:      []string{"har-benchmark", "android", "not-claimed"},
	}, func(a *allurex.A) {
		a.Step("Load docs/har-benchmark/7-android-blocker.json", func() {
			path := filepath.Join(moduleRoot(t), "docs", "har-benchmark", "7-android-blocker.json")
			raw, err := os.ReadFile(path)
			require.NoError(t, err, "committed SSOT at %s", path)

			var doc androidHarBlocker
			require.NoError(t, json.Unmarshal(raw, &doc))

			require.Equal(t, 6, doc.Phase)
			require.Equal(t, "6-android-hub-enableHAR", doc.Label)
			require.Equal(t, "not_claimed", doc.Status)
			require.Equal(t, "blocker", doc.Verdict)
			require.Nil(t, doc.Scorecard)
			require.Nil(t, doc.HARPath)
			require.Empty(t, doc.AllowedClaims)
			require.NotEmpty(t, doc.ForbiddenClaims)
			require.GreaterOrEqual(t, len(doc.Blockers), 3)

			ids := map[string]bool{}
			for _, b := range doc.Blockers {
				ids[b.ID] = true
			}
			require.True(t, ids["no-7070"])
			require.True(t, ids["hub-cdp-gate"])
			require.True(t, ids["ui-no-cap"])
		})

		a.Step("Load docs/har-benchmark/7-android.NOT_CLAIMED.txt", func() {
			path := filepath.Join(moduleRoot(t), "docs", "har-benchmark", "7-android.NOT_CLAIMED.txt")
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			text := string(raw)
			require.Contains(t, text, "not claimed")
			require.Contains(t, text, "7070")
			require.Contains(t, text, "Appium")
		})
	})
}
