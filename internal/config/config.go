// Package config loads Owner-compatible Java .properties profiles.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Config is the runtime test profile (keys match TestConfig / *.properties).
type Config struct {
	Env                      string
	HubURL                   string
	UIURL                    string
	APIBaseURL               string
	HubStatusPath            string
	RemoteURL                string
	Browser                  string
	BrowserVersion           string
	ChromeVersion            string
	ChromeMinVersion         string
	FirefoxVersion           string
	FirefoxMinVersion        string
	MsedgeVersion            string
	MsedgeMinVersion         string
	PlaywrightWsEndpoint     string
	PlaywrightSessionName    string
	PlaywrightSessionTimeout string
	PlaywrightEnableVNC      bool
	PlaywrightEnableVideo    bool
	SmokeURL                 string
	LogToConsole             bool
	SkipHealthCheck          bool
	CmHubPort                int
	CmUiPort                 int
	CmBinaryPath             string
	CmBrowsersJSON           string
	CmSelenoidBinary         string
	CmSelenoidUiBinary       string
	CmUseLocalBinaries       bool
}

var (
	loadOnce sync.Once
	loaded   *Config
	loadErr  error
)

// Load returns the process-wide config (defaults < env file < process env).
func Load() (*Config, error) {
	loadOnce.Do(func() {
		loaded, loadErr = load()
	})
	return loaded, loadErr
}

// MustLoad panics if config cannot be loaded.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// ResetForTest clears the cached config (unit tests only).
func ResetForTest() {
	loadOnce = sync.Once{}
	loaded = nil
	loadErr = nil
}

// FromMap builds a Config from Owner-style key overrides (unit tests).
// Unset keys use TestConfig @DefaultValue equivalents for resolve coverage.
func FromMap(overrides map[string]string) *Config {
	props := map[string]string{
		"hubUrl":                   "http://127.0.0.1:4444/",
		"uiUrl":                    "http://127.0.0.1:8080/",
		"apiBaseUrl":               "",
		"hubStatusPath":            "/status",
		"remoteUrl":                "",
		"playwrightWsEndpoint":     "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
		"playwrightSessionName":    "java-playwright-tests",
		"playwrightSessionTimeout": "5m",
		"playwrightEnableVnc":      "false",
		"playwrightEnableVideo":    "false",
	}
	for k, v := range overrides {
		props[k] = v
	}
	return configFromProps("from-map", props)
}

// ChromeVersionForSession returns chromeVersion or browserVersion (TestConfig.chromeVersion parity).
func (c *Config) ChromeVersionForSession() string {
	if strings.TrimSpace(c.ChromeVersion) != "" {
		return strings.TrimSpace(c.ChromeVersion)
	}
	return strings.TrimSpace(c.BrowserVersion)
}

// FirefoxVersionForSession returns firefoxVersion (TestConfig.firefoxVersion).
func (c *Config) FirefoxVersionForSession() string {
	return strings.TrimSpace(c.FirefoxVersion)
}

// MsedgeVersionForSession returns msedgeVersion (TestConfig.msedgeVersion).
func (c *Config) MsedgeVersionForSession() string {
	return strings.TrimSpace(c.MsedgeVersion)
}

// ChromeMinVersionForSession returns chromeMinVersion (TestConfig.chromeMinVersion).
func (c *Config) ChromeMinVersionForSession() string {
	return strings.TrimSpace(c.ChromeMinVersion)
}

// FirefoxMinVersionForSession returns firefoxMinVersion (TestConfig.firefoxMinVersion).
func (c *Config) FirefoxMinVersionForSession() string {
	return strings.TrimSpace(c.FirefoxMinVersion)
}

// MsedgeMinVersionForSession returns msedgeMinVersion (TestConfig.msedgeMinVersion).
func (c *Config) MsedgeMinVersionForSession() string {
	return strings.TrimSpace(c.MsedgeMinVersion)
}

func configFromProps(envName string, props map[string]string) *Config {
	return &Config{
		Env:                      envName,
		HubURL:                   strings.TrimSpace(props["hubUrl"]),
		UIURL:                    strings.TrimSpace(props["uiUrl"]),
		APIBaseURL:               strings.TrimSpace(props["apiBaseUrl"]),
		HubStatusPath:            normalizePath(firstNonEmpty(props["hubStatusPath"], "/status")),
		RemoteURL:                strings.TrimSpace(props["remoteUrl"]),
		Browser:                  firstNonEmpty(props["browser"], "chrome"),
		BrowserVersion:           firstNonEmpty(props["browserVersion"], "149.0"),
		ChromeVersion:            firstNonEmpty(props["chromeVersion"], props["browserVersion"]),
		ChromeMinVersion:         firstNonEmpty(props["chromeMinVersion"], MinVersion("chrome")),
		FirefoxVersion:           firstNonEmpty(props["firefoxVersion"], "151.0"),
		FirefoxMinVersion:        firstNonEmpty(props["firefoxMinVersion"], MinVersion("firefox")),
		MsedgeVersion:            firstNonEmpty(props["msedgeVersion"], "145.0"),
		MsedgeMinVersion:         firstNonEmpty(props["msedgeMinVersion"], MinVersion("msedge")),
		PlaywrightWsEndpoint:     strings.TrimSpace(props["playwrightWsEndpoint"]),
		PlaywrightSessionName:    firstNonEmpty(props["playwrightSessionName"], "java-playwright-tests"),
		PlaywrightSessionTimeout: firstNonEmpty(props["playwrightSessionTimeout"], "5m"),
		PlaywrightEnableVNC:      parseBool(props["playwrightEnableVnc"], false),
		PlaywrightEnableVideo:    parseBool(props["playwrightEnableVideo"], false),
		SmokeURL:                 firstNonEmpty(props["smokeUrl"], "https://example.com/"),
		LogToConsole:             parseBool(props["logToConsole"], true),
		SkipHealthCheck:          parseBool(firstNonEmpty(os.Getenv("SELENOID_TEST_SKIP_HEALTH_CHECK"), props["skipHealthCheck"]), false),
		CmHubPort:                parseInt(firstNonEmpty(props["cmHubPort"], "4445"), 4445),
		CmUiPort:                 parseInt(firstNonEmpty(props["cmUiPort"], "8081"), 8081),
		CmBinaryPath:             firstNonEmpty(props["cmBinaryPath"], "../cm/cm"),
		CmBrowsersJSON:           firstNonEmpty(props["cmBrowsersJson"], "../dev/browsers.json"),
		CmSelenoidBinary:         firstNonEmpty(props["cmSelenoidBinary"], "../dev/bin/selenoid"),
		CmSelenoidUiBinary:       firstNonEmpty(props["cmSelenoidUiBinary"], "../dev/bin/selenoid-ui"),
		CmUseLocalBinaries:       parseBool(props["cmUseLocalBinaries"], false),
	}
}

func load() (*Config, error) {
	envName := firstNonEmpty(os.Getenv("SELENOID_TEST_ENV"), os.Getenv("env"), "local")
	root, err := findModuleRoot()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(root, "src", "test", "resources", "config")

	props := map[string]string{}
	if err := mergePropertiesFile(props, filepath.Join(configDir, "default.properties")); err != nil {
		return nil, err
	}
	envFile := filepath.Join(configDir, envName+".properties")
	if _, err := os.Stat(envFile); err == nil {
		if err := mergePropertiesFile(props, envFile); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	applyEnvOverrides(props)

	cfg := configFromProps(envName, props)
	// Keep Load() URLs slash-normalized for API clients (P0 contract).
	cfg.HubURL = withSlash(cfg.HubURL)
	cfg.UIURL = withSlash(cfg.UIURL)
	cfg.APIBaseURL = withSlash(cfg.APIBaseURL)
	if cfg.HubURL == "" {
		return nil, fmt.Errorf("set hubUrl in config/%s.properties", envName)
	}
	if cfg.UIURL == "" {
		return nil, fmt.Errorf("set uiUrl in config/%s.properties", envName)
	}
	return cfg, nil
}

// APIBase returns apiBaseUrl or hubUrl (Owner resolveApiBaseUrl).
func (c *Config) APIBase() string {
	if strings.TrimSpace(c.APIBaseURL) != "" {
		return withSlash(c.APIBaseURL)
	}
	return withSlash(c.HubURL)
}

// HubStatusURL is GET target for hub capacity (prod: /hub/status).
func (c *Config) HubStatusURL() string {
	return strings.TrimRight(c.APIBase(), "/") + c.ResolveHubStatusPath()
}

func findModuleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}

func mergePropertiesFile(dst map[string]string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		dst[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	return sc.Err()
}

func applyEnvOverrides(props map[string]string) {
	set := func(propKey, envKey string) {
		if v := os.Getenv(envKey); v != "" {
			props[propKey] = v
		}
	}
	// Java -Dkey / system property style (plain key in process env).
	set("hubUrl", "hubUrl")
	set("uiUrl", "uiUrl")
	set("apiBaseUrl", "apiBaseUrl")
	set("hubStatusPath", "hubStatusPath")
	set("remoteUrl", "remoteUrl")
	set("browser", "browser")
	set("browserVersion", "browserVersion")
	set("chromeVersion", "chromeVersion")
	set("chromeMinVersion", "chromeMinVersion")
	set("firefoxVersion", "firefoxVersion")
	set("firefoxMinVersion", "firefoxMinVersion")
	set("msedgeVersion", "msedgeVersion")
	set("msedgeMinVersion", "msedgeMinVersion")
	set("logToConsole", "logToConsole")
	set("playwrightWsEndpoint", "playwrightWsEndpoint")
	set("playwrightSessionName", "playwrightSessionName")
	set("playwrightSessionTimeout", "playwrightSessionTimeout")
	set("playwrightEnableVnc", "playwrightEnableVnc")
	set("playwrightEnableVideo", "playwrightEnableVideo")
	set("smokeUrl", "smokeUrl")
	set("cmBinaryPath", "cmBinaryPath")
	set("cmBrowsersJson", "cmBrowsersJson")
	set("cmSelenoidBinary", "cmSelenoidBinary")
	set("cmSelenoidUiBinary", "cmSelenoidUiBinary")
	set("cmUseLocalBinaries", "cmUseLocalBinaries")
	set("cmHubPort", "cmHubPort")
	set("cmUiPort", "cmUiPort")
	// Explicit SELENOID_TEST_* overrides.
	set("hubUrl", "SELENOID_TEST_HUB_URL")
	set("uiUrl", "SELENOID_TEST_UI_URL")
	set("apiBaseUrl", "SELENOID_TEST_API_BASE_URL")
	set("hubStatusPath", "SELENOID_TEST_HUB_STATUS_PATH")
	set("remoteUrl", "SELENOID_TEST_REMOTE_URL")
	set("logToConsole", "SELENOID_TEST_LOG_TO_CONSOLE")
}

func withSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if !strings.HasSuffix(s, "/") {
		return s + "/"
	}
	return s
}

func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/status"
	}
	if !strings.HasPrefix(p, "/") {
		return "/" + p
	}
	return p
}

func parseBool(s string, def bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return def
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func parseInt(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil || n <= 0 {
		return def
	}
	return n
}
