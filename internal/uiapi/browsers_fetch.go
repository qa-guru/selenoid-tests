package uiapi

import (
	"github.com/qa-guru/selenoid-tests/internal/config"
)

// FetchBrowsersConfig GET UI /browsers-config.
func FetchBrowsersConfig(cfg *config.Config) (BrowsersConfig, error) {
	body, err := uiClient(cfg).GetBytes("/browsers-config")
	if err != nil {
		return nil, err
	}
	return ParseBrowsersConfig(body)
}
