package uiapi

import (
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// FetchStatus GET UI /status and returns proxied hub counters (.state or flat).
func FetchStatus(cfg *config.Config) (*hubapi.Status, error) {
	client := httpx.New(cfg.UIURL)
	body, err := client.GetBytes("/status")
	if err != nil {
		return nil, err
	}
	return hubapi.Parse(body)
}
