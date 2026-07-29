package cm_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/cm"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestCmInstallerHelper_WithTempConfigDirCreatesConfigDirectory(t *testing.T) {
	cfg := config.FromMap(map[string]string{
		"cmBinaryPath": "../cm/cm",
	})
	allurex.Run(t, unitMeta("withTempConfigDir creates isolated config directory", "helpers.CmInstallerHelperTest", "CmInstallerHelper paths"), func(a *allurex.A) {
		var installer *cm.Installer
		a.Step("create temp config dir", func() {
			var err error
			installer, err = cm.WithTempConfigDir(cfg)
			require.NoError(t, err)
		})
		defer installer.DeleteConfigDir()

		a.Step("verify directory layout", func() {
			info, err := os.Stat(installer.ConfigDir())
			require.NoError(t, err)
			require.True(t, info.IsDir())
			_, err = os.Stat(installer.BrowsersJSONPath())
			require.Error(t, err)
		})
	})
}
