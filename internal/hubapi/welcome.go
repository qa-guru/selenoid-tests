package hubapi

import (
	"github.com/qa-guru/selenoid-tests/internal/config"
)

// FetchWelcomeText GET hub / welcome page.
func FetchWelcomeText(cfg *config.Config) (string, error) {
	return hubClient(cfg).GetText("/")
}
