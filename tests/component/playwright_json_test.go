package component_test

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
		Package:   "tests.component.PlaywrightWsPathJsonTest",
		Layer:     "component",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright WS path fixture",
		Suite:     "Playwright WS path fixture",
		Tags:      []string{"component"},
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
		Package:   "tests.component.PlaywrightBrowserCapsJsonTest",
		Layer:     "component",
		Component: "playwright-image",
		Epic:      "playwright-image",
		Feature:   "Playwright browser caps fixture",
		Suite:     "Playwright browser caps fixture",
		Tags:      []string{"component"},
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
