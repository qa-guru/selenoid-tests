package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReader_ResolveHubURLAddsTrailingSlash(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubUrl adds trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"hubUrl": "http://127.0.0.1:4444"})
			got, err := cfg.ResolveHubURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:4444/", got)
		})
	})
}

func TestConfigReader_ResolveHubURLKeepsTrailingSlash(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubUrl keeps trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"hubUrl": "http://127.0.0.1:4444/"})
			got, err := cfg.ResolveHubURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:4444/", got)
		})
	})
}

func TestConfigReader_ResolveHubURLFailsWhenEmpty(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubUrl fails when empty"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"hubUrl": ""})
			_, err := cfg.ResolveHubURL()
			require.Error(t, err)
			require.Contains(t, err.Error(), "hubUrl")
		})
	})
}

func TestConfigReader_ResolveUiURLAddsTrailingSlash(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveUiUrl adds trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"uiUrl": "http://127.0.0.1:8080"})
			got, err := cfg.ResolveUiURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:8080/", got)
		})
	})
}

func TestConfigReader_ResolveAPIBaseURLPrefersExplicitKey(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveApiBaseUrl prefers apiBaseUrl"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"apiBaseUrl": "https://autotests.ai/stack/backend-java-spring",
				"hubUrl":     "https://selenoid.qa.guru/",
			})
			got, err := cfg.ResolveAPIBaseURL()
			require.NoError(t, err)
			require.Equal(t, "https://autotests.ai/stack/backend-java-spring/", got)
		})
	})
}

func TestConfigReader_ResolveAPIBaseURLFallsBackToHubURL(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveApiBaseUrl falls back to hubUrl"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"apiBaseUrl": "", "hubUrl": "http://127.0.0.1:4444"})
			got, err := cfg.ResolveAPIBaseURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:4444/", got)
		})
	})
}

func TestConfigReader_ResolveAPIBaseURLFailsWhenBothEmpty(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveApiBaseUrl fails when both empty"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"apiBaseUrl": "", "hubUrl": ""})
			_, err := cfg.ResolveAPIBaseURL()
			require.Error(t, err)
			require.Contains(t, err.Error(), "apiBaseUrl or hubUrl")
		})
	})
}

func TestConfigReader_ResolveHubStatusPathDefaultsToStatus(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubStatusPath defaults to /status"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{})
			require.Equal(t, "/status", cfg.ResolveHubStatusPath())
		})
	})
}

func TestConfigReader_ResolveHubStatusPathNormalizesPath(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubStatusPath normalizes path"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"hubStatusPath": "hub/status"})
			require.Equal(t, "/hub/status", cfg.ResolveHubStatusPath())
		})
	})
}

func unitMeta(component, name string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   "config",
		Layer:     "unit",
		Component: component,
		Suite:     "unit",
		Tags:      []string{"unit"},
	}
}
