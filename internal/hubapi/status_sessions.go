package hubapi

import (
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// CollectSessionIDs returns all session ids from hub /status browsers tree.
func CollectSessionIDs(cfg *config.Config) (map[string]struct{}, error) {
	st, err := Fetch(cfg)
	if err != nil {
		return nil, err
	}
	ids := map[string]struct{}{}
	collectSessionIDs(st.Browsers, ids)
	return ids, nil
}

// WaitNewSessionID polls /status until a session id not in before appears.
func WaitNewSessionID(cfg *config.Config, before map[string]struct{}, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		now, err := CollectSessionIDs(cfg)
		if err != nil {
			return "", err
		}
		for id := range now {
			if _, seen := before[id]; !seen {
				return id, nil
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	return "", nil
}

func collectSessionIDs(node any, ids map[string]struct{}) {
	switch v := node.(type) {
	case map[string]any:
		if sessions, ok := v["sessions"].([]any); ok {
			for _, item := range sessions {
				if sess, ok := item.(map[string]any); ok {
					if id, ok := sess["id"].(string); ok && id != "" {
						ids[id] = struct{}{}
					}
				}
			}
		}
		for _, child := range v {
			collectSessionIDs(child, ids)
		}
	case []any:
		for _, child := range v {
			collectSessionIDs(child, ids)
		}
	}
}
