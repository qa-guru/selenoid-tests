package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReader_ResolveCmHubURLUsesCmHubPort(t *testing.T) {
	allurex.Run(t, cmUnitMeta("resolveCmHubUrl uses cmHubPort with trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"cmHubPort": "4445"})
			require.Equal(t, "http://127.0.0.1:4445/", cfg.ResolveCmHubURL())
		})
	})
}

func TestConfigReader_ResolveCmHubURLRespectsCustomPort(t *testing.T) {
	allurex.Run(t, cmUnitMeta("resolveCmHubUrl respects custom cmHubPort"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"cmHubPort": "9444"})
			require.Equal(t, "http://127.0.0.1:9444/", cfg.ResolveCmHubURL())
		})
	})
}

func TestConfigReader_ResolveCmUiURLUsesCmUiPort(t *testing.T) {
	allurex.Run(t, cmUnitMeta("resolveCmUiUrl uses cmUiPort with trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"cmUiPort": "8081"})
			require.Equal(t, "http://127.0.0.1:8081/", cfg.ResolveCmUiURL())
		})
	})
}

func TestConfigReader_ResolveCmUiURLRespectsCustomPort(t *testing.T) {
	allurex.Run(t, cmUnitMeta("resolveCmUiUrl respects custom cmUiPort"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"cmUiPort": "9080"})
			require.Equal(t, "http://127.0.0.1:9080/", cfg.ResolveCmUiURL())
		})
	})
}

func TestConfigReader_ResolveCmRemoteURLCombinesHubBaseAndWdHubPath(t *testing.T) {
	allurex.Run(t, cmUnitMeta("resolveCmRemoteUrl points WebDriver hub at CM-managed port"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"cmHubPort": "4445"})
			require.Equal(t, "http://127.0.0.1:4445/wd/hub", cfg.ResolveCmRemoteURL())
		})
	})
}

func cmUnitMeta(name string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   "config.ConfigReaderCmTest",
		Layer:     "unit",
		Component: "cm",
		Epic:      "cm",
		Suite:     "ConfigReader CM URLs",
		Tags:      []string{"unit", "cm"},
	}
}
