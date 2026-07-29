package helpers

import (
	"encoding/json"
	"strings"
)

// HarStats is a lightweight HAR 1.2 metrics summary (helpers.HarStats).
type HarStats struct {
	Label              string
	Entries            int
	HTTPEntries        int
	WithStatus         int
	WithResponseHeaders int
	WithRequestHeaders int
	WithContentSize    int
	WithContentText    int
	WithPostData       int
	WithPositiveTime   int
	Creator            string
}

// HarStatsFromBytes parses HAR JSON bytes into metrics.
func HarStatsFromBytes(label string, harBytes []byte) HarStats {
	var root map[string]any
	_ = json.Unmarshal(harBytes, &root)
	log, _ := root["log"].(map[string]any)
	if log == nil {
		log = map[string]any{}
	}
	entries, _ := log["entries"].([]any)
	creatorMap, _ := log["creator"].(map[string]any)
	creator := stringVal(creatorMap["name"])

	stats := HarStats{Label: label, Entries: len(entries), Creator: creator}
	for _, e := range entries {
		entry, _ := e.(map[string]any)
		req, _ := entry["request"].(map[string]any)
		resp, _ := entry["response"].(map[string]any)
		if req == nil {
			req = map[string]any{}
		}
		if resp == nil {
			resp = map[string]any{}
		}
		content, _ := resp["content"].(map[string]any)
		if content == nil {
			content = map[string]any{}
		}
		url := stringVal(req["url"])
		if isHTTPURL(url) {
			stats.HTTPEntries++
		}
		if intVal(resp["status"]) > 0 {
			stats.WithStatus++
		}
		if listSize(resp["headers"]) > 0 {
			stats.WithResponseHeaders++
		}
		if listSize(req["headers"]) > 0 {
			stats.WithRequestHeaders++
		}
		if longVal(content["size"]) > 0 {
			stats.WithContentSize++
		}
		if text, ok := content["text"].(string); ok && strings.TrimSpace(text) != "" {
			stats.WithContentText++
		}
		if pd, ok := req["postData"].(map[string]any); ok && len(pd) > 0 {
			stats.WithPostData++
		}
		if doubleVal(entry["time"]) > 0 {
			stats.WithPositiveTime++
		}
	}
	return stats
}

func isHTTPURL(url string) bool {
	u := strings.ToLower(url)
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

func listSize(o any) int {
	switch v := o.(type) {
	case []any:
		return len(v)
	case []map[string]string:
		return len(v)
	default:
		return 0
	}
}
