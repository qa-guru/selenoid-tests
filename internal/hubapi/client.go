package hubapi

import (
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

func hubClient(cfg *config.Config) *httpx.Client {
	base := strings.TrimRight(cfg.APIBase(), "/")
	return httpx.New(base)
}
