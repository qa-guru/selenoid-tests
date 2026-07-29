package hubapi

import (
	"encoding/json"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// Ping is hub GET /ping payload.
type Ping struct {
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

// FetchPing GET hub /ping.
func FetchPing(cfg *config.Config) (*Ping, error) {
	body, err := hubClient(cfg).GetBytes("/ping")
	if err != nil {
		return nil, err
	}
	var ping Ping
	if err := json.Unmarshal(body, &ping); err != nil {
		return nil, err
	}
	return &ping, nil
}
