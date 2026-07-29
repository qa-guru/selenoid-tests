package hubapi

import "encoding/json"

// ParseLogsList unmarshals hub GET /logs JSON array of session log file names.
func ParseLogsList(body []byte) ([]string, error) {
	var files []string
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}
	return files, nil
}
