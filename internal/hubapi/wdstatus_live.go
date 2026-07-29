package hubapi

import (
	"github.com/qa-guru/selenoid-tests/internal/config"
)

// FetchWebDriverStatus GET /wd/hub/status.
func FetchWebDriverStatus(cfg *config.Config) (*WebDriverStatus, error) {
	body, err := hubClient(cfg).GetBytes("/wd/hub/status")
	if err != nil {
		return nil, err
	}
	return ParseWebDriverStatus(body)
}
