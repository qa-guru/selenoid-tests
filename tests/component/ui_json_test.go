package component_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiStatusJson_ParsesProxiedWrapper(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses proxied UI status wrapper",
		Package:   "tests.component.UiStatusJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "UI status fixture",
		Suite:     "UI status fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/ui/status.json", func() {
			resp, err := uiapi.ParseStatusResponse(loadFixture(t, "ui/status.json"))
			require.NoError(t, err)
			require.NotNil(t, resp.State)
			require.Equal(t, 1, resp.State.Used)
		})
	})
}

func TestUiPingJson_ParsesUptimeAndVersion(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses uptime and version fields",
		Package:   "tests.component.UiPingJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "UI ping fixture",
		Suite:     "UI ping fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/ui/ping.json", func() {
			ping, err := uiapi.ParsePing(loadFixture(t, "ui/ping.json"))
			require.NoError(t, err)
			require.Equal(t, "1h2m3s", ping.Uptime)
			require.Equal(t, "1.0.0-test", ping.Version)
		})
	})
}

func TestUiErrorJson_ParsesErrorList(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses error list payload",
		Package:   "tests.component.UiErrorJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "UI error fixture",
		Suite:     "UI error fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/ui/error.json", func() {
			errResp, err := uiapi.ParseError(loadFixture(t, "ui/error.json"))
			require.NoError(t, err)
			require.NotNil(t, errResp.Errors)
			require.NotEmpty(t, errResp.Errors)
		})
	})
}

func TestBrowsersConfigJson_ParsesCatalogMap(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses browsers catalog map",
		Package:   "tests.component.BrowsersConfigJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "Browsers config fixture",
		Suite:     "Browsers config fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/ui/browsers-config.json", func() {
			cfg, err := uiapi.ParseBrowsersConfig(loadFixture(t, "ui/browsers-config.json"))
			require.NoError(t, err)
			chrome := cfg["chrome"]
			require.Contains(t, chrome, "148.0")
			require.Contains(t, chrome, "149.0")
		})
	})
}
