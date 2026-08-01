package min_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestChromeMinCatalogJson_ParsesVersionBlockPortAndPath(t *testing.T) {
	minVersion := config.MinVersion("chrome")
	allurex.Run(t, allurex.Meta{
		Name:      "parses chrome-min catalog version block port and path",
		Package:   "tests.unit.fixture.ChromeMinCatalogJsonTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Suite:     "Chrome min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("load catalog version block", func() {
			block := config.VersionBlock("chrome", minVersion)
			require.Equal(t, "4444", fmt.Sprint(block["port"]))
			require.Equal(t, "/", fmt.Sprint(block["path"]))
		})
	})
}

func TestChromeMinCatalogJson_ParsesMinImageTag(t *testing.T) {
	minVersion := config.MinVersion("chrome")
	allurex.Run(t, allurex.Meta{
		Name:      "parses min image tag in chrome catalog entry",
		Package:   "tests.unit.fixture.ChromeMinCatalogJsonTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Suite:     "Chrome min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("verify min image tag", func() {
			block := config.VersionBlock("chrome", minVersion)
			require.True(t, strings.Contains(fmt.Sprint(block["image"]), config.MinImageMajor("chrome")+"-min"))
		})
	})
}
