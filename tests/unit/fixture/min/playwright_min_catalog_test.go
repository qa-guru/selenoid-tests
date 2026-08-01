package min_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestPlaywrightMinCatalogJson_ParsesPlaywrightVersion(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses playwrightVersion from min catalog version block",
		Package:   "tests.unit.fixture.PlaywrightMinCatalogJsonTest",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Suite:     "Playwright min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("parse min catalog version block", func() {
			cat, err := playwrightapi.ParseCatalog(loadFixture(t, "playwright/browser-catalog.json"))
			require.NoError(t, err)
			family := cat["playwright-chromium"]
			minVer := family.Default + "-min"
			block := family.Versions[minVer]
			require.Equal(t, family.Default, block.PlaywrightVersion)
		})
	})
}

func TestPlaywrightMinCatalogJson_ParsesMinImageTag(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses min image tag in catalog entry",
		Package:   "tests.unit.fixture.PlaywrightMinCatalogJsonTest",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Suite:     "Playwright min browser catalog fixture",
		Tags:      []string{"unit", "min"},
	}, func(a *allurex.A) {
		a.Step("verify min protocol and image tag", func() {
			cat, err := playwrightapi.ParseCatalog(loadFixture(t, "playwright/browser-catalog.json"))
			require.NoError(t, err)
			family := cat["playwright-chromium"]
			minVer := family.Default + "-min"
			block := family.Versions[minVer]
			require.Equal(t, "playwright", block.Protocol)
			require.True(t, strings.Contains(block.Image, minVer))
		})
	})
}
