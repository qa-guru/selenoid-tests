package hubapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// ListVideoJSON GET /video/?json with pagination params.
func ListVideoJSON(cfg *config.Config, limit, offset int, q string) (*VideoListResponse, error) {
	query := url.Values{}
	query.Set("json", "")
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))
	if strings.TrimSpace(q) != "" {
		query.Set("q", q)
	}
	body, err := hubClient(cfg).GetQuery("/video/", query)
	if err != nil {
		return nil, err
	}
	return ParseVideoList(body)
}

// FindVideoBySessionID searches paginated video list for a session id substring.
func FindVideoBySessionID(cfg *config.Config, sessionID string) (string, error) {
	listed, err := ListVideoJSON(cfg, 10, 0, sessionID)
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

// DownloadVideo GET /video/{fileName} bytes.
func DownloadVideo(cfg *config.Config, fileName string) ([]byte, error) {
	return hubClient(cfg).GetBytes("/video/" + fileName)
}

// GetVideoExpectStatus GET /video/{fileName} and require HTTP status.
func GetVideoExpectStatus(cfg *config.Config, fileName string, expected int) error {
	_, err := hubClient(cfg).GetExpectStatus("/video/"+fileName, expected)
	return err
}
