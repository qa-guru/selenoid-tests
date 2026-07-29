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
