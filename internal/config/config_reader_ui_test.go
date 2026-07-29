package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReaderUI_ResolveUiURLFailsWhenEmpty(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid-ui", "ConfigReader.resolveUiUrl fails when empty"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"uiUrl": ""})
			_, err := cfg.ResolveUiURL()
			require.Error(t, err)
			require.Contains(t, err.Error(), "uiUrl")
		})
	})
}

func TestConfigReaderUI_ResolveUiURLKeepsTrailingSlash(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid-ui", "ConfigReader.resolveUiUrl keeps trailing slash"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"uiUrl": "http://127.0.0.1:8080/"})
			got, err := cfg.ResolveUiURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:8080/", got)
		})
	})
}
