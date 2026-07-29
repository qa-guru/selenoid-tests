package playwrightapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// AssertUpgradeRequired GET playwright path without WebSocket upgrade (expects 400).
func AssertUpgradeRequired(cfg *config.Config) error {
	path, err := ResolveHTTPPath(cfg)
	if err != nil {
		return err
	}
	return getExpectStatus(cfg, path, http.StatusBadRequest)
}

// AssertUnknownPathRejected GET unknown playwright path (expects 400).
func AssertUnknownPathRejected(cfg *config.Config) error {
	rawQuery := playwrightRawQuery(cfg)
	path := withRawQuery("/playwright/unknown-browser/0.0.0", rawQuery)
	return getExpectStatus(cfg, path, http.StatusBadRequest)
}

// ResolveHTTPPath maps playwrightWsEndpoint ws(s) URL to HTTP GET path (+ query).
func ResolveHTTPPath(cfg *config.Config) (string, error) {
	endpoint := strings.TrimSpace(cfg.PlaywrightWsEndpoint)
	if endpoint == "" {
		return "", fmt.Errorf("playwrightWsEndpoint is empty")
	}
	httpURI := strings.Replace(strings.Replace(endpoint, "ws://", "http://", 1), "wss://", "https://", 1)
	u, err := url.Parse(httpURI)
	if err != nil {
		return "", err
	}
	if u.Path == "" || u.Path == "/" {
		return "", fmt.Errorf("playwright path not found in: %s", endpoint)
	}
	return withRawQuery(u.Path, u.RawQuery), nil
}

func playwrightRawQuery(cfg *config.Config) string {
	endpoint := strings.TrimSpace(cfg.PlaywrightWsEndpoint)
	httpURI := strings.Replace(strings.Replace(endpoint, "ws://", "http://", 1), "wss://", "https://", 1)
	u, err := url.Parse(httpURI)
	if err != nil {
		return ""
	}
	return u.RawQuery
}

func withRawQuery(path, rawQuery string) string {
	if rawQuery == "" {
		return path
	}
	return path + "?" + rawQuery
}

func getExpectStatus(cfg *config.Config, path string, expected int) error {
	base := strings.TrimRight(cfg.APIBase(), "/")
	_, err := httpx.New(base).GetExpectStatus(path, expected)
	return err
}
