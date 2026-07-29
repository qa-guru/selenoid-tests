package uiapi

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

const sseRequestTimeout = 30 * time.Second

// ReadFirstSSEEvent reads one data: payload from UI /events.
func ReadFirstSSEEvent(cfg *config.Config) (*SseEvent, error) {
	return readSSEEvents(cfg, 1)
}

// ReadTwoSSEEvents reads two consecutive SSE payloads from UI /events.
func ReadTwoSSEEvents(cfg *config.Config) ([]*SseEvent, error) {
	return readSSEEventsList(cfg, 2)
}

func readSSEEvents(cfg *config.Config, count int) (*SseEvent, error) {
	events, err := readSSEEventsList(cfg, count)
	if err != nil {
		return nil, err
	}
	return events[0], nil
}

func readSSEEventsList(cfg *config.Config, count int) ([]*SseEvent, error) {
	base := strings.TrimRight(cfg.UIURL, "/")
	url := base + "/events"
	client := &http.Client{Timeout: sseRequestTimeout}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSE GET %s returned HTTP %d", url, resp.StatusCode)
	}
	return parseSSEPayloads(resp.Body, count)
}

func parseSSEPayloads(r io.Reader, count int) ([]*SseEvent, error) {
	scanner := bufio.NewScanner(r)
	var events []*SseEvent
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		ev, err := ParseSseEvent([]byte(payload))
		if err != nil {
			return nil, err
		}
		events = append(events, ev)
		if len(events) >= count {
			return events, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if count == 1 {
		return nil, fmt.Errorf("SSE stream ended without data event")
	}
	return nil, fmt.Errorf("expected %d SSE events, got %d", count, len(events))
}
