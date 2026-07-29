package playwrightapi

import "encoding/json"

// Catalog is a browsers.json-shaped map of family name → entry.
type Catalog map[string]Family

// Family is one browser family block (default + versions).
type Family struct {
	Default  string                  `json:"default"`
	Versions map[string]VersionBlock `json:"versions"`
}

// VersionBlock is one version entry under a family.
type VersionBlock struct {
	Image             string `json:"image"`
	Protocol          string `json:"protocol"`
	PlaywrightVersion string `json:"playwrightVersion"`
}

// ParseCatalog unmarshals a Playwright (or generic) browser catalog JSON.
func ParseCatalog(body []byte) (Catalog, error) {
	var cat Catalog
	if err := json.Unmarshal(body, &cat); err != nil {
		return nil, err
	}
	return cat, nil
}
