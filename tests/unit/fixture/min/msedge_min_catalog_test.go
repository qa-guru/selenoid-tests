package min_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestMsedgeMinCatalogJson_ParsesVersionBlockPortAndPath(t *testing.T) {
	minVersion := config.MinVersion("msedge")
	allurex.Run(t, allurex.Meta{
		Name:      "parses msedge-min catalog version block port and path",
		Package:   "tests.unit.fixture.MsedgeMinCatalogJsonTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Suite:     "Edge min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("load catalog version block", func() {
			block := config.VersionBlock("msedge", minVersion)
			require.Equal(t, "4444", fmt.Sprint(block["port"]))
			require.Equal(t, "/", fmt.Sprint(block["path"]))
		})
	})
}

func TestMsedgeMinCatalogJson_ParsesMinImageTag(t *testing.T) {
	minVersion := config.MinVersion("msedge")
	allurex.Run(t, allurex.Meta{
		Name:      "parses min image tag in msedge catalog entry",
		Package:   "tests.unit.fixture.MsedgeMinCatalogJsonTest",
		Layer:     "unit",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Suite:     "Edge min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("verify min image tag", func() {
			block := config.VersionBlock("msedge", minVersion)
			require.True(t, strings.Contains(fmt.Sprint(block["image"]), config.MinImageMajor("msedge")+"-min"))
		})
	})
}
