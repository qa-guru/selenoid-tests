package uiapi

import (
	"encoding/json"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// Ping is UI GET /ping payload.
type Ping struct {
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

// ParsePing unmarshals a UI /ping JSON body.
func ParsePing(body []byte) (*Ping, error) {
	var ping Ping
	if err := json.Unmarshal(body, &ping); err != nil {
		return nil, err
	}
	return &ping, nil
}

// FetchPing GET UI /ping.
func FetchPing(cfg *config.Config) (*Ping, error) {
	client := httpx.New(cfg.UIURL)
	body, err := client.GetBytes("/ping")
	if err != nil {
		return nil, err
	}
	return ParsePing(body)
}
