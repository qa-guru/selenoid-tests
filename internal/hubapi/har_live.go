package hubapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// HarListResponse is hub GET /har/?json paginated list (HubHarApi; videos JSON field).
type HarListResponse struct {
	Videos []string `json:"videos"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

// ParseHarList unmarshals a paginated /har/?json body.
func ParseHarList(body []byte) (*HarListResponse, error) {
	var listed HarListResponse
	if err := json.Unmarshal(body, &listed); err != nil {
		return nil, err
	}
	return &listed, nil
}

// ListHarJSON GET /har/?json with pagination params.
func ListHarJSON(cfg *config.Config, limit, offset int, q string) (*HarListResponse, error) {
	query := url.Values{}
	query.Set("json", "")
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))
	if strings.TrimSpace(q) != "" {
		query.Set("q", q)
	}
	body, err := hubClient(cfg).GetQuery("/har/", query)
	if err != nil {
		return nil, err
	}
	return ParseHarList(body)
}

// FindHarBySessionID searches paginated HAR list for a session id substring.
func FindHarBySessionID(cfg *config.Config, sessionID string) (string, error) {
	listed, err := ListHarJSON(cfg, 10, 0, sessionID)
	if err != nil {
		return "", err
	}
	for _, name := range listed.Videos {
		if strings.Contains(name, sessionID) {
			return name, nil
		}
	}
	return "", nil
}

// DownloadHar GET /har/{fileName} bytes.
func DownloadHar(cfg *config.Config, fileName string) ([]byte, error) {
	return hubClient(cfg).GetBytes("/har/" + fileName)
}

// WaitForSessionHar polls GET /har/?json&q=sessionId until a match appears or timeout elapses.
func WaitForSessionHar(cfg *config.Config, sessionID string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		match, err := FindHarBySessionID(cfg, sessionID)
		if err != nil {
			return "", err
		}
		if match != "" {
			return match, nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return FindHarBySessionID(cfg, sessionID)
}
