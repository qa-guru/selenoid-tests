package hubapi

import "encoding/json"

// VideoListResponse is hub GET /video/?json paginated list (HubVideoApi).
type VideoListResponse struct {
	Videos []string `json:"videos"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

// ParseVideoList unmarshals a paginated /video/?json body.
func ParseVideoList(body []byte) (*VideoListResponse, error) {
	var listed VideoListResponse
	if err := json.Unmarshal(body, &listed); err != nil {
		return nil, err
	}
	return &listed, nil
}
