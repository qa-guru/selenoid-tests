package hubapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// SessionArchiveEntry is one row from hub GET /sessions/?json.
type SessionArchiveEntry struct {
	ID       string `json:"id"`
	Video    string `json:"video,omitempty"`
	Log      string `json:"log,omitempty"`
	HAR      string `json:"har,omitempty"`
	Name     string `json:"name,omitempty"`
	Quota    string `json:"quota,omitempty"`
}

// SessionsArchiveResponse is the paginated finished-session listing.
type SessionsArchiveResponse struct {
	Sessions []SessionArchiveEntry `json:"sessions"`
	Total    int                   `json:"total"`
	Limit    int                   `json:"limit"`
	Offset   int                   `json:"offset"`
}

// ListSessionsJSON GET /sessions/?json with pagination params.
func ListSessionsJSON(cfg *config.Config, limit, offset int, q string) (*SessionsArchiveResponse, error) {
	query := url.Values{}
	query.Set("json", "")
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))
	if strings.TrimSpace(q) != "" {
		query.Set("q", q)
	}
	body, err := hubClient(cfg).GetQuery("/sessions/", query)
	if err != nil {
		return nil, err
	}
	var listed SessionsArchiveResponse
	if err := json.Unmarshal(body, &listed); err != nil {
		return nil, err
	}
	return &listed, nil
}

// FindArchivedSession returns the exact session row from /sessions/?json?q=id.
func FindArchivedSession(cfg *config.Config, sessionID string) (*SessionArchiveEntry, error) {
	listed, err := ListSessionsJSON(cfg, 50, 0, sessionID)
	if err != nil {
		return nil, err
	}
	for _, row := range listed.Sessions {
		if row.ID == sessionID {
			copy := row
			return &copy, nil
		}
	}
	return nil, nil
}

// WaitArchivedSessionHar polls /sessions/?json until the session row includes a har file.
func WaitArchivedSessionHar(cfg *config.Config, sessionID string, timeout time.Duration) (*SessionArchiveEntry, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		row, err := FindArchivedSession(cfg, sessionID)
		if err != nil {
			return nil, err
		}
		if row != nil && strings.TrimSpace(row.HAR) != "" {
			return row, nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return FindArchivedSession(cfg, sessionID)
}
