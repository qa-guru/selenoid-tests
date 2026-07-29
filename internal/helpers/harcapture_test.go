package helpers_test

import (
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
