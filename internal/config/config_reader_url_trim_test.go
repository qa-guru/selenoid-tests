package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReaderURLTrim_ResolveHubURLTrimsWhitespace(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveHubUrl trims whitespace"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"hubUrl": "  http://127.0.0.1:4444  "})
			got, err := cfg.ResolveHubURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:4444/", got)
		})
	})
}

func TestConfigReaderURLTrim_ResolveUiURLTrimsWhitespace(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveUiUrl trims whitespace"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"uiUrl": "  http://127.0.0.1:8080  "})
			got, err := cfg.ResolveUiURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:8080/", got)
		})
	})
}

func TestConfigReaderURLTrim_ResolveAPIBaseURLTrimsWhitespace(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "ConfigReader.resolveApiBaseUrl trims whitespace"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"apiBaseUrl": "  http://api.example.com  ",
				"hubUrl":     "http://hub/",
			})
			got, err := cfg.ResolveAPIBaseURL()
			require.NoError(t, err)
			require.Equal(t, "http://api.example.com/", got)
		})
	})
}
