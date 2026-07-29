package uiapi

import (
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

func uiClient(cfg *config.Config) *httpx.Client {
	base := strings.TrimRight(cfg.UIURL, "/")
	return httpx.New(base)
}
