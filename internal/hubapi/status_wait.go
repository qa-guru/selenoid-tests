package hubapi

import (
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// WaitUntilUsed polls hub status until used == expected or timeout elapses.
func WaitUntilUsed(cfg *config.Config, expected int, timeout time.Duration) (*Status, error) {
	deadline := time.Now().Add(timeout)
	var last *Status
	for time.Now().Before(deadline) {
		st, err := Fetch(cfg)
		if err != nil {
			return nil, err
		}
		last = st
		if st.Used == expected {
			return st, nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	if last != nil {
		return last, nil
	}
	return Fetch(cfg)
}
