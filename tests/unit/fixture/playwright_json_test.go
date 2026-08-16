package fixture_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/playwrightapi"
)

func TestPlaywrightWsPathJson_DefaultAndProtocol(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses playwright browser catalog default version",
		Package:   "tests.unit.fixture.PlaywrightWsPathJsonTest",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS path fixture",
		Suite:     "Playwright WS path fixture",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		cat, err := playwrightapi.ParseCatalog(loadFixture(t, "playwright/browser-catalog.json"))
		require.NoError(t, err)
		family := cat["playwright-chromium"]
		a.Step("default version", func() {
			require.Equal(t, "1.61.1", family.Default)
		})
		a.Step("protocol flag in catalog entry", func() {
			block := family.Versions["1.61.1"]
			require.Equal(t, "playwright", block.Protocol)
			require.True(t, strings.Contains(block.Image, "playwright-chromium"))
		})
	})
}

func TestPlaywrightBrowserCapsJson_VersionAndFamily(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses playwrightVersion from catalog version block",
		Package:   "tests.unit.fixture.PlaywrightBrowserCapsJsonTest",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright browser caps fixture",
		Suite:     "Playwright browser caps fixture",
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		cat, err := playwrightapi.ParseCatalog(loadFixture(t, "playwright/browser-catalog.json"))
		require.NoError(t, err)
		a.Step("playwrightVersion from version block", func() {
			require.Equal(t, "1.61.1", cat["playwright-chromium"].Versions["1.61.1"].PlaywrightVersion)
		})
		a.Step("lists playwright-chromium family", func() {
			_, ok := cat["playwright-chromium"]
			require.True(t, ok)
		})
	})
}

func TestPlaywrightChromeCatalog_PresentInCiBrowsers(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "ci-browsers.json lists playwright-chrome family",
		Package:   "tests.unit.fixture.PlaywrightChromeCatalogTests",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright Chrome catalog",
		Story:     "Playwright Chrome catalog",
		Suite:     "Playwright Chrome catalog",
		Browser:   allurex.BrowserChrome,
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("verify playwright-chrome family in fixtures/ci-browsers.json", func() {
			cat, err := playwrightapi.ParseCatalog(loadProjectFixture(t, "fixtures/ci-browsers.json"))
			require.NoError(t, err)
			family, ok := cat["playwright-chrome"]
			require.True(t, ok, "playwright-chrome key")
			block := family.Versions[family.Default]
			require.Equal(t, "playwright", block.Protocol)
			require.Contains(t, block.Image, "playwright-chrome")
		})
	})
}

func TestPlaywrightMsedgeCatalog_PresentInCiBrowsers(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "ci-browsers.json lists playwright-msedge family",
		Package:   "tests.unit.fixture.PlaywrightMsedgeCatalogTests",
		Layer:     "unit",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright Microsoft Edge catalog",
		Story:     "Playwright Microsoft Edge catalog",
		Suite:     "Playwright Microsoft Edge catalog",
		Browser:   allurex.BrowserMsedge,
		Tags:      []string{"unit"},
	}, func(a *allurex.A) {
		a.Step("verify playwright-msedge family in fixtures/ci-browsers.json", func() {
			cat, err := playwrightapi.ParseCatalog(loadProjectFixture(t, "fixtures/ci-browsers.json"))
			require.NoError(t, err)
			family, ok := cat["playwright-msedge"]
			require.True(t, ok, "playwright-msedge key")
			block := family.Versions[family.Default]
			require.Equal(t, "playwright", block.Protocol)
			require.Contains(t, block.Image, "playwright-msedge")
		})
	})
}
