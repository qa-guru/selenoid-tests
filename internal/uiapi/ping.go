package uiapi

import (
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// Ping is UI GET /ping payload.
type Ping struct {
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

// FetchPing GET UI /ping.
func FetchPing(cfg *config.Config) (*Ping, error) {
	client := httpx.New(cfg.UIURL)
	var ping Ping
	if err := client.GetJSON("/ping", &ping); err != nil {
		return nil, err
	}
	return &ping, nil
}
