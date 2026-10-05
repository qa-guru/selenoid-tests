package config_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestLoad_LocalUnitDefaults(t *testing.T) {
	config.ResetForTest()
	t.Setenv("SELENOID_TEST_ENV", "local_unit")
	_ = os.Unsetenv("hubUrl")
	_ = os.Unsetenv("SELENOID_TEST_HUB_URL")

	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "local_unit", cfg.Env)
	require.Equal(t, "http://127.0.0.1:4444/", cfg.HubURL)
	require.Equal(t, "http://127.0.0.1:8080/", cfg.UIURL)
	require.Equal(t, "/status", cfg.HubStatusPath)
}

func TestLoad_ProdApiHubStatusPath(t *testing.T) {
	config.ResetForTest()
	t.Setenv("SELENOID_TEST_ENV", "selenoid_qa_guru_api")

	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "/hub/status", cfg.HubStatusPath)
	require.Contains(t, cfg.APIBase(), "selenoid.qa.guru")
}

func TestLoad_StandWebDriverVersionsFollowCICatalog(t *testing.T) {
	raw, err := os.ReadFile("../../fixtures/ci-browsers.json")
	require.NoError(t, err)
	var catalog map[string]struct {
		Default string `json:"default"`
	}
	require.NoError(t, json.Unmarshal(raw, &catalog))
	for _, key := range []string{"browserVersion", "chromeVersion", "chromeMinVersion", "firefoxVersion", "firefoxMinVersion", "msedgeVersion", "msedgeMinVersion"} {
		t.Setenv(key, "")
	}
	for _, tc := range []struct {
		profile string
		min     bool
	}{
		{"selenoid_qa_guru_api", false},
		{"selenoid_qa_guru_e2e", false},
		{"selenoid_qa_guru_integration", false},
		{"selenoid_github_api", true},
		{"selenoid_github_e2e", true},
		{"selenoid_github_integration", true},
		{"selenoid_github_min_integration", false},
		{"selenoid_github_cm_integration", true},
	} {
		t.Run(tc.profile, func(t *testing.T) {
			config.ResetForTest()
			t.Cleanup(config.ResetForTest)
			t.Setenv("SELENOID_TEST_ENV", tc.profile)
			cfg, err := config.Load()
			require.NoError(t, err)
			browserVersion := catalog["chrome"].Default
			if tc.min {
				browserVersion += "-min"
			}
			require.Equal(t, browserVersion, cfg.BrowserVersion)
			require.Equal(t, catalog["chrome"].Default, cfg.ChromeVersionForSession())
			require.Equal(t, catalog["chrome"].Default+"-min", cfg.ChromeMinVersionForSession())
			require.Equal(t, catalog["firefox"].Default, cfg.FirefoxVersionForSession())
			require.Equal(t, catalog["firefox"].Default+"-min", cfg.FirefoxMinVersionForSession())
			require.Equal(t, catalog["msedge"].Default, cfg.MsedgeVersionForSession())
			require.Equal(t, catalog["msedge"].Default+"-min", cfg.MsedgeMinVersionForSession())
		})
	}
}

func TestLoad_ProdWebDriverVersionOverridesCatalog(t *testing.T) {
	config.ResetForTest()
	t.Cleanup(config.ResetForTest)
	t.Setenv("SELENOID_TEST_ENV", "selenoid_qa_guru_api")
	for key, value := range map[string]string{
		"browserVersion": "custom-browser", "chromeVersion": "custom-chrome",
		"chromeMinVersion": "custom-chrome-min", "firefoxVersion": "custom-firefox",
		"firefoxMinVersion": "custom-firefox-min", "msedgeVersion": "custom-msedge",
		"msedgeMinVersion": "custom-msedge-min",
	} {
		t.Setenv(key, value)
	}
	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "custom-browser", cfg.BrowserVersion)
	require.Equal(t, "custom-chrome", cfg.ChromeVersionForSession())
	require.Equal(t, "custom-chrome-min", cfg.ChromeMinVersionForSession())
	require.Equal(t, "custom-firefox", cfg.FirefoxVersionForSession())
	require.Equal(t, "custom-firefox-min", cfg.FirefoxMinVersionForSession())
	require.Equal(t, "custom-msedge", cfg.MsedgeVersionForSession())
	require.Equal(t, "custom-msedge-min", cfg.MsedgeMinVersionForSession())
}

func TestAdvertisedCatalogMustStart(t *testing.T) {
	require.True(t, config.FromMap(map[string]string{}).AdvertisedCatalogMustStart() == false)
	qa := config.FromMap(map[string]string{})
	qa.Env = "selenoid_qa_guru_e2e"
	require.True(t, qa.AdvertisedCatalogMustStart())
	gh := config.FromMap(map[string]string{})
	gh.Env = "selenoid_github_min_integration"
	require.True(t, gh.AdvertisedCatalogMustStart())
	local := config.FromMap(map[string]string{})
	local.Env = "local_integration"
	require.False(t, local.AdvertisedCatalogMustStart())
}
