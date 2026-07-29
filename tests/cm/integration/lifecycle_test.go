package cmint_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/cm"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func integrationMeta(name, pkg, feature, story string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   pkg,
		Layer:     "integration",
		Component: "cm",
		Epic:      "cm",
		Feature:   feature,
		Story:     story,
		Suite:     story,
		Tags:      []string{"integration", "cm", "positive"},
	}
}

func newInstaller(t *testing.T) *cm.Installer {
	t.Helper()
	cfg := config.MustLoad()
	installer, err := cm.WithTempConfigDir(cfg)
	require.NoError(t, err)
	installer.StopAll()
	t.Cleanup(func() {
		installer.StopAll()
		installer.DeleteConfigDir()
	})
	return installer
}

func TestCmInstallerLifecycle_ConfigureWritesBrowsersJSON(t *testing.T) {
	installer := newInstaller(t)
	allurex.Run(t, integrationMeta(
		"configure writes browsers.json from dev catalog",
		"tests.integration.CmInstallerLifecycleTests",
		"Installer lifecycle",
		"Installer lifecycle",
	), func(a *allurex.A) {
		var result cm.RunResult
		a.Step("cm selenoid configure -n", func() {
			var err error
			result, err = installer.Configure()
			require.NoError(t, err)
		})
		a.Step("Verify configure exit code", func() {
			result.RequireSuccess("configure")
		})
		a.Step("Verify browsers.json exists and contains chrome", func() {
			path := installer.BrowsersJSONPath()
			info, err := os.Stat(path)
			require.NoError(t, err)
			require.False(t, info.IsDir())

			body, err := os.ReadFile(path)
			require.NoError(t, err)
			var root map[string]any
			require.NoError(t, json.Unmarshal(body, &root))
			chrome := root["chrome"].(map[string]any)
			defaultVersion, _ := chrome["default"].(string)
			require.NotEmpty(t, defaultVersion)
		})
		a.Step("cm selenoid status", func() {
			status := installer.StatusHub()
			status.RequireSuccess("status")
			require.Contains(t, status.Output, "configuration file")
		})
	})
}

func TestCmInstallerLifecycle_StartHubExposesStatusEndpoint(t *testing.T) {
	installer := newInstaller(t)
	cfg := config.MustLoad()
	allurex.Run(t, integrationMeta(
		"start exposes hub /status",
		"tests.integration.CmInstallerLifecycleTests",
		"Installer lifecycle",
		"Installer lifecycle",
	), func(a *allurex.A) {
		a.Step("Configure CM stack", func() {
			result, err := installer.Configure()
			require.NoError(t, err)
			result.RequireSuccess("configure")
		})
		a.Step("Start hub via cm", func() {
			result, err := installer.StartHub()
			require.NoError(t, err)
			result.RequireSuccess("start hub")
		})
		a.Step("Wait for hub readiness", func() {
			require.NoError(t, installer.WaitForHubReady(240*time.Second))
		})
		a.Step("cm selenoid status", func() {
			status := installer.StatusHub()
			status.RequireSuccess("status")
			require.Contains(t, status.Output, "container is running")
		})
		a.Step("GET hub /status", func() {
			_, err := hubapi.FetchFrom(cfg.ResolveCmHubURL())
			require.NoError(t, err)
		})
	})
}

func TestCmInstallerLifecycle_StartFullStackUiProxiesHub(t *testing.T) {
	installer := newInstaller(t)
	cfg := config.MustLoad()
	allurex.Run(t, integrationMeta(
		"full stack start — UI /status mirrors hub counters",
		"tests.integration.CmInstallerLifecycleTests",
		"Installer lifecycle",
		"Installer lifecycle",
	), func(a *allurex.A) {
		a.Step("Configure CM stack", func() {
			result, err := installer.Configure()
			require.NoError(t, err)
			result.RequireSuccess("configure")
		})
		a.Step("Start hub via cm", func() {
			result, err := installer.StartHub()
			require.NoError(t, err)
			result.RequireSuccess("start hub")
		})
		a.Step("Wait for hub readiness", func() {
			require.NoError(t, installer.WaitForHubReady(240*time.Second))
		})
		a.Step("Start UI via cm", func() {
			result, err := installer.StartUi()
			require.NoError(t, err)
			result.RequireSuccess("start UI")
		})
		a.Step("Wait for UI readiness", func() {
			require.NoError(t, installer.WaitForUiReady(120*time.Second))
		})

		var hubStatus, uiStatus *hubapi.Status
		a.Step("GET hub /status", func() {
			var err error
			hubStatus, err = hubapi.FetchFrom(cfg.ResolveCmHubURL())
			require.NoError(t, err)
		})
		a.Step("GET UI /status", func() {
			var err error
			uiStatus, err = uiapi.FetchStatusFrom(cfg.ResolveCmUiURL())
			require.NoError(t, err)
		})
		a.Step("Verify proxied counters match hub", func() {
			require.Equal(t, hubStatus.Total, uiStatus.Total)
			require.Equal(t, hubStatus.Used, uiStatus.Used)
			require.Equal(t, hubStatus.Queued, uiStatus.Queued)
			require.Equal(t, hubStatus.Pending, uiStatus.Pending)
		})
	})
}

func TestCmCliVersion_PrintsRevision(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, integrationMeta(
		"cm version exits zero and prints revision",
		"tests.integration.CmCliVersionTests",
		"CM CLI",
		"CM CLI version",
	), func(a *allurex.A) {
		var result cm.RunResult
		a.Step("Run cm version", func() {
			var err error
			result, err = cm.Version(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify exit code and output", func() {
			require.Equal(t, 0, result.ExitCode, result.Output)
			require.Contains(t, result.Output, "Git Revision:")
		})
	})
}

func TestCmCliHelp_ListsSubcommands(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, integrationMeta(
		"cm --help exits zero and lists subcommands",
		"tests.integration.CmCliVersionTests",
		"CM CLI",
		"CM CLI version",
	), func(a *allurex.A) {
		var result cm.RunResult
		a.Step("Run cm --help", func() {
			var err error
			result, err = cm.Help(cfg)
			require.NoError(t, err)
		})
		a.Step("Verify exit code and output", func() {
			require.Equal(t, 0, result.ExitCode, result.Output)
			require.True(t, strings.Contains(result.Output, "selenoid-ui"))
		})
	})
}
