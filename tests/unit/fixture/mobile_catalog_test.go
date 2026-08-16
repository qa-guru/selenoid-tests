package fixture_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func TestAndroidCatalog_PresentInCiBrowsers(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "ci-browsers.json lists android Appium family",
		Package:   "tests.unit.fixture.AndroidCatalogTests",
		Layer:     "unit",
		Component: "android",
		Epic:      "android",
		Feature:   "Appium catalog",
		Story:     "Android catalog",
		Suite:     "Android catalog",
		Browser:   allurex.BrowserAndroid,
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("verify android family in fixtures/ci-browsers.json", func() {
			var root map[string]any
			require.NoError(t, json.Unmarshal(loadProjectFixture(t, "fixtures/ci-browsers.json"), &root))
			android, ok := root["android"].(map[string]any)
			require.True(t, ok, "android key")
			require.Equal(t, "16.0", android["default"])
			versions, ok := android["versions"].(map[string]any)
			require.True(t, ok)
			require.Contains(t, versions, "16.0")
		})
	})
}

func TestIosCatalog_NotClaimedInCiBrowsers(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "ci-browsers.json has no ios family yet (roadmap stub)",
		Package:   "tests.unit.fixture.IosCatalogTests",
		Layer:     "unit",
		Component: "ios",
		Epic:      "ios",
		Feature:   "Appium catalog",
		Story:     "iOS catalog stub",
		Suite:     "iOS catalog stub",
		Browser:   allurex.BrowserIOS,
		Tags:      []string{"unit", "future"},
	}, func(a *allurex.A) {
		a.Step("ios key is absent until qaguru/ios ships", func() {
			var root map[string]any
			require.NoError(t, json.Unmarshal(loadProjectFixture(t, "fixtures/ci-browsers.json"), &root))
			_, ok := root["ios"]
			require.False(t, ok, "replace this stub with Appium session smoke when ios is added to the catalog")
		})
	})
}
