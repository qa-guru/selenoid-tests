package ui_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
	inthelpers "github.com/qa-guru/selenoid-tests/tests/integration/internal"
)

func TestUiBrowsersConfigIntegration_BrowsersConfigIncludesHubBrowserVersions(t *testing.T) {
	cfg := config.MustLoad()
	version := cfg.BrowserVersion
	allurex.Run(t, allurex.Meta{
		Name:      "UI /browsers-config includes chrome versions from hub /status",
		Package:   "tests.integration.UiBrowsersConfigIntegrationTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI browsers config proxy",
		Story:     "UI browsers config vs hub status",
		Suite:     "UI browsers config vs hub status",
		Tags:      []string{"integration", "positive"},
	}, func(a *allurex.A) {
		var hubBrowsers map[string]any
		var uiCatalog uiapi.BrowsersConfig
		a.Step("GET hub /status browsers", func() {
			st, err := hubapi.Fetch(cfg)
			require.NoError(t, err)
			require.NotNil(t, st.Browsers)
			hubBrowsers = st.Browsers
		})
		a.Step("GET UI /browsers-config", func() {
			var err error
			uiCatalog, err = uiapi.FetchBrowsersConfig(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify configured chrome version appears in hub status and UI catalog", func() {
			require.NotNil(t, hubBrowsers["chrome"], "hub should expose chrome family")
			require.Contains(t, uiCatalog, "chrome", "UI catalog should expose chrome family")
			inthelpers.AssertBrowserVersionListed(t, hubBrowsers, "chrome", version)
			inthelpers.AssertBrowserVersionListed(t, browsersConfigAsAny(uiCatalog), "chrome", version)
		})
	})
}

func browsersConfigAsAny(cfg uiapi.BrowsersConfig) map[string]any {
	out := make(map[string]any, len(cfg))
	for family, versions := range cfg {
		inner := make(map[string]any, len(versions))
		for ver, block := range versions {
			inner[ver] = block
		}
		out[family] = inner
	}
	return out
}
