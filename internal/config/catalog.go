package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CatalogResource is the classpath-relative fixture used by Java WebDriverCatalog.
const CatalogResource = "fixtures/webdriver/browser-catalog.json"

type catalogBrowser struct {
	Default  string                    `json:"default"`
	Versions map[string]map[string]any `json:"versions"`
}

var (
	catalogOnce sync.Once
	catalogData map[string]catalogBrowser
	catalogErr  error
)

func loadCatalog() (map[string]catalogBrowser, error) {
	catalogOnce.Do(func() {
		root, err := findModuleRoot()
		if err != nil {
			catalogErr = err
			return
		}
		path := filepath.Join(root, "src", "test", "resources", CatalogResource)
		raw, err := os.ReadFile(path)
		if err != nil {
			catalogErr = fmt.Errorf("fixture not found: %s: %w", CatalogResource, err)
			return
		}
		var data map[string]catalogBrowser
		if err := json.Unmarshal(raw, &data); err != nil {
			catalogErr = err
			return
		}
		catalogData = data
	})
	return catalogData, catalogErr
}

func applyStandWebDriverCatalog(root string, props map[string]string) error {
	path := filepath.Join(root, "fixtures", "ci-browsers.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load WebDriver catalog: %w", err)
	}
	var catalog map[string]catalogBrowser
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return fmt.Errorf("decode WebDriver catalog: %w", err)
	}
	for _, browser := range []string{"chrome", "firefox", "msedge"} {
		entry := catalog[browser]
		if entry.Default == "" || entry.Versions[entry.Default] == nil || entry.Versions[entry.Default+"-min"] == nil {
			return fmt.Errorf("WebDriver catalog lacks default/min versions for %s", browser)
		}
		props[browser+"Version"] = entry.Default
		props[browser+"MinVersion"] = entry.Default + "-min"
	}
	browser := firstNonEmpty(props["browser"], "chrome")
	if catalog[browser].Default == "" {
		return fmt.Errorf("WebDriver catalog lacks default for %s", browser)
	}
	version := catalog[browser].Default
	if strings.HasSuffix(props["browserVersion"], "-min") {
		version += "-min"
	}
	props["browserVersion"] = version
	return nil
}

// DefaultVersion returns catalog default version for browser (WebDriverCatalog.defaultVersion).
func DefaultVersion(browser string) string {
	data, err := loadCatalog()
	if err != nil {
		panic(err)
	}
	b, ok := data[browser]
	if !ok || b.Default == "" {
		panic(fmt.Sprintf("browser not in catalog: %s", browser))
	}
	return b.Default
}

// MinVersion returns default+"-min" (WebDriverCatalog.minVersion).
func MinVersion(browser string) string {
	return DefaultVersion(browser) + "-min"
}

// VersionBlock returns the versions[version] object from the catalog.
func VersionBlock(browser, version string) map[string]any {
	data, err := loadCatalog()
	if err != nil {
		panic(err)
	}
	b, ok := data[browser]
	if !ok {
		panic(fmt.Sprintf("browser not in catalog: %s", browser))
	}
	block, ok := b.Versions[version]
	if !ok {
		panic(fmt.Sprintf("version not in catalog: %s/%s", browser, version))
	}
	return block
}

// MinImageMajor is the major prefix of the default version (e.g. "149" from "149.0").
func MinImageMajor(browser string) string {
	def := DefaultVersion(browser)
	dot := strings.Index(def, ".")
	if dot < 0 {
		return def
	}
	return def[:dot]
}
