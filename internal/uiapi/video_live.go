package uiapi

import (
	"fmt"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

// GetProxyExpectStatus GET UI {prefix}/{sessionID} and require HTTP status.
func GetProxyExpectStatus(cfg *config.Config, prefix, sessionID string, expected int) error {
	path := fmt.Sprintf("%s/%s", prefix, sessionID)
	_, err := uiClient(cfg).GetExpectStatus(path, expected)
	return err
}

// ListVideoJSON GET UI /video/?json with pagination.
func ListVideoJSON(cfg *config.Config, limit, offset int, q string) (*hubapi.VideoListResponse, error) {
	query := buildVideoQuery(limit, offset, q)
	body, err := uiClient(cfg).GetQuery("/video/", query)
	if err != nil {
		return nil, err
	}
	return hubapi.ParseVideoList(body)
}

// FindVideoBySessionID searches UI video list for session id.
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

// DownloadVideo GET UI /video/{fileName}.
func DownloadVideo(cfg *config.Config, fileName string) ([]byte, error) {
	return uiClient(cfg).GetBytes("/video/" + fileName)
}

// GetVideoExpectStatus GET UI /video/{fileName} and require status.
func GetVideoExpectStatus(cfg *config.Config, fileName string, expected int) error {
	_, err := uiClient(cfg).GetExpectStatus("/video/"+fileName, expected)
	return err
}
