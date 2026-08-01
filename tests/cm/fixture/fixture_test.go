package cm_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func loadFixture(t *testing.T, rel string) string {
	t.Helper()
	path := filepath.Join(moduleRoot(t), "src", "test", "resources", filepath.FromSlash(rel))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture %s: %v", rel, err)
	}
	return string(body)
}

func loadProjectFixture(t *testing.T, rel string) string {
	t.Helper()
	path := filepath.Join(moduleRoot(t), filepath.FromSlash(rel))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture %s: %v", rel, err)
	}
	return string(body)
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found from test cwd")
		}
		dir = parent
	}
}

func unitMeta(name, pkg, suite string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   pkg,
		Layer:     "unit",
		Component: "cm",
		Epic:      "cm",
		Suite:     suite,
		Tags:      []string{"unit", "cm"},
	}
}

func TestCmBrowsersConfigJson_ParsesDefaultChromeVersion(t *testing.T) {
	allurex.Run(t, unitMeta(
		"parses cm browsers.json default chrome version",
		"tests.unit.fixture.CmBrowsersConfigJsonTest",
		"CM browsers.json fixture",
	), func(a *allurex.A) {
		a.Step("parse ci-browsers.json", func() {
			var root map[string]any
			require.NoError(t, json.Unmarshal([]byte(loadProjectFixture(t, "fixtures/ci-browsers.json")), &root))

			chrome, ok := root["chrome"].(map[string]any)
			require.True(t, ok)
			defaultVersion, _ := chrome["default"].(string)
			require.Regexp(t, regexp.MustCompile(`\d+\.\d+`), defaultVersion)

			require.Contains(t, root, "firefox")
			require.Contains(t, root, "msedge")

			versions, ok := chrome["versions"].(map[string]any)
			require.True(t, ok)
			require.Contains(t, versions, defaultVersion)
			require.Contains(t, versions, defaultVersion+"-min")
		})
	})
}

func TestCmBrowsersConfigJson_ParsesChromeMinImage(t *testing.T) {
	allurex.Run(t, unitMeta(
		"parses cm browsers.json chrome-min image tag",
		"tests.unit.fixture.CmBrowsersConfigJsonTest",
		"CM browsers.json fixture",
	), func(a *allurex.A) {
		a.Step("verify chrome-min image tag", func() {
			var root map[string]any
			require.NoError(t, json.Unmarshal([]byte(loadProjectFixture(t, "fixtures/ci-browsers.json")), &root))
			chrome := root["chrome"].(map[string]any)
			defaultVersion := chrome["default"].(string)
			minVersion := defaultVersion + "-min"
			major := defaultVersion[:strings.Index(defaultVersion, ".")]
			versions := chrome["versions"].(map[string]any)
			block := versions[minVersion].(map[string]any)
			image := block["image"].(string)
			require.Contains(t, image, "webdriver-chrome")
			require.Contains(t, image, major+"-min")
		})
	})
}

func TestCmHelpOutput_ListsCoreSubcommands(t *testing.T) {
	allurex.Run(t, unitMeta(
		"lists selenoid and selenoid-ui subcommands in help",
		"tests.unit.fixture.CmHelpOutputTest",
		"CM help output fixture",
	), func(a *allurex.A) {
		a.Step("parse help fixture", func() {
			output := loadFixture(t, "fixtures/cm/help-output.txt")
			require.Contains(t, output, "selenoid")
			require.Contains(t, output, "selenoid-ui")
			require.Contains(t, output, "version")
		})
	})
}

func TestCmStatusOutput_DetectsRunningContainer(t *testing.T) {
	allurex.Run(t, unitMeta(
		"detects running container in cm status output",
		"tests.unit.fixture.CmStatusOutputTest",
		"CM status output fixture",
	), func(a *allurex.A) {
		a.Step("parse status fixture", func() {
			output := loadFixture(t, "fixtures/cm/status-running.txt")
			require.Contains(t, output, "Configuration file:")
			require.Contains(t, output, "container is running on port 4445")
		})
	})
}

func TestCmVersionOutput_ParsesGitRevision(t *testing.T) {
	allurex.Run(t, unitMeta(
		"parses git revision from cm version output",
		"tests.unit.fixture.CmVersionOutputTest",
		"CM version output fixture",
	), func(a *allurex.A) {
		a.Step("parse version fixture", func() {
			output := loadFixture(t, "fixtures/cm/version-output.txt")
			require.Contains(t, output, "Git Revision:")
			require.Contains(t, output, "UTC Build Time:")
		})
	})
}
