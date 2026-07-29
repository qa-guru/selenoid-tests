package uiapi

import (
	"net/http"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

// FetchStatusWhenHubUnavailable GET UI /status and expects HTTP 500 with error payload
// (Java UiStatusApi.fetchWhenHubUnavailable parity).
func FetchStatusWhenHubUnavailable(cfg *config.Config) (*ErrorResponse, error) {
	base := strings.TrimRight(strings.TrimSpace(cfg.UIURL), "/")
	client := httpx.New(base)
	resp, err := client.GetExpectStatus("/status", http.StatusInternalServerError)
	if err != nil {
		return nil, err
	}
	return ParseError(resp.Body)
}
