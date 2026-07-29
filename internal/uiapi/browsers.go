package uiapi

import "encoding/json"

// BrowsersConfig is UI /browsers-config catalog: family → version → block.
type BrowsersConfig map[string]map[string]any

// ParseBrowsersConfig unmarshals a browsers-config JSON object.
func ParseBrowsersConfig(body []byte) (BrowsersConfig, error) {
	var cfg BrowsersConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
