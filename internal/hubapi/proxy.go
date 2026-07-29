package hubapi

import (
	"fmt"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// GetProxyExpectStatus GET {prefix}/{sessionID} and require HTTP status.
func GetProxyExpectStatus(cfg *config.Config, prefix, sessionID string, expected int) error {
	path := fmt.Sprintf("%s/%s", prefix, sessionID)
	_, err := hubClient(cfg).GetExpectStatus(path, expected)
	return err
}
