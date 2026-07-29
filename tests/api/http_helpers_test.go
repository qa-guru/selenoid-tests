package api_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

func hubHTTP(cfg *config.Config) *httpx.Client {
	return httpx.New(strings.TrimRight(cfg.APIBase(), "/"))
}

func uiHTTP(cfg *config.Config) *httpx.Client {
	return httpx.New(strings.TrimRight(cfg.UIURL, "/"))
}

func getExpectStatusBodyContains(client *httpx.Client, path string, expected int, substr string) error {
	resp, err := client.Do(http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != expected {
		return fmt.Errorf("GET %s: HTTP %d", path, resp.StatusCode)
	}
	if !strings.Contains(string(resp.Body), substr) {
		return fmt.Errorf("GET %s: body missing %q", path, substr)
	}
	return nil
}
