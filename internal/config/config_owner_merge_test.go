package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigOwnerMerge_SystemPropertyOverridesFileDefault(t *testing.T) {
	allurex.Run(t, unitMeta("selenoid", "env/system hubUrl overrides file default"), func(a *allurex.A) {
		a.Step("load with hubUrl override", func() {
			config.ResetForTest()
			t.Setenv("SELENOID_TEST_ENV", "local_unit")
			_ = os.Unsetenv("SELENOID_TEST_HUB_URL")
			t.Setenv("hubUrl", "https://override.autotests.ai:4444")

			cfg, err := config.Load()
			require.NoError(t, err)
			require.Equal(t, "https://override.autotests.ai:4444/", cfg.HubURL)
		})
	})
}
