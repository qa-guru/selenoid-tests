// Package warmpool is the HTTP client for selenoid-warm-pool orchestrator (port 9090).
package warmpool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/qa-guru/selenoid-tests/internal/httpx"
)

const defaultBaseURL = "http://127.0.0.1:9090"

// Slot is one orchestrator pool entry (GET /pool/slots, POST /pool/reserve).
type Slot struct {
	ID                   string  `json:"id"`
	Protocol             string  `json:"protocol"`
	Browser              string  `json:"browser"`
	SessionID            string  `json:"sessionId"`
	WarmURL              string  `json:"warmUrl"`
	WebdriverURL         string  `json:"webdriverUrl"`
	WebdriverURLLoopback string  `json:"webdriverUrlLoopback"`
	PlaywrightWsURL      string  `json:"playwrightWsUrl"`
	ReservedBy           *string `json:"reservedBy"`
}

// Health is GET /health (and GET /).
type Health struct {
	OK    bool `json:"ok"`
	Slots int  `json:"slots"`
}

// Client talks to the orchestrator.
type Client struct {
	http *httpx.Client
}

// BaseURL from WARM_POOL_URL / SELENOID_WARM_POOL_URL, else local stand :9090.
func BaseURL() string {
	for _, key := range []string{"WARM_POOL_URL", "SELENOID_WARM_POOL_URL"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return strings.TrimRight(v, "/")
		}
	}
	return defaultBaseURL
}

// New wraps an orchestrator base URL.
func New(baseURL string) *Client {
	return &Client{http: httpx.New(strings.TrimRight(strings.TrimSpace(baseURL), "/"))}
}

// Default uses BaseURL().
func Default() *Client {
	return New(BaseURL())
}

// Ping GET /health.
func (c *Client) Ping() error {
	var h Health
	return c.http.GetJSON("/health", &h)
}

// Health GET /health.
func (c *Client) Health() (*Health, error) {
	var h Health
	if err := c.http.GetJSON("/health", &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// Root GET / (stand URL gate; same payload as health).
func (c *Client) Root() (*Health, error) {
	var h Health
	if err := c.http.GetJSON("/", &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// Slots GET /pool/slots.
func (c *Client) Slots() ([]Slot, error) {
	var slots []Slot
	if err := c.http.GetJSON("/pool/slots", &slots); err != nil {
		return nil, err
	}
	return slots, nil
}

// Reserve POST /pool/reserve. HTTP 409 is not an error — Status=409, Slot empty.
func (c *Client) Reserve(protocol, browser, owner string, loopback bool) (slot Slot, status int, err error) {
	resp, err := c.http.PostJSON("/pool/reserve", map[string]any{
		"protocol": protocol,
		"browser":  browser,
		"owner":    owner,
		"loopback": loopback,
	}, 0)
	if err != nil {
		return Slot{}, 0, err
	}
	if resp.StatusCode == http.StatusConflict {
		return Slot{}, resp.StatusCode, nil
	}
	if resp.StatusCode != http.StatusOK {
		return Slot{}, resp.StatusCode, fmt.Errorf("reserve HTTP %d: %s", resp.StatusCode, truncate(string(resp.Body), 200))
	}
	var out struct {
		OK   bool `json:"ok"`
		Slot Slot `json:"slot"`
	}
	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return Slot{}, resp.StatusCode, err
	}
	return out.Slot, resp.StatusCode, nil
}

// Release POST /pool/release.
func (c *Client) Release(slotID string) error {
	_, err := c.http.PostJSON("/pool/release", map[string]string{"slotId": slotID}, http.StatusOK)
	return err
}

// DialURL is the host-reachable WebDriver URL (loopback field wins).
func (s Slot) DialURL() string {
	if strings.TrimSpace(s.WebdriverURLLoopback) != "" {
		return s.WebdriverURLLoopback
	}
	return s.WebdriverURL
}

// IsLoopback reports whether DialURL is 127.0.0.1 / localhost / ::1.
func (s Slot) IsLoopback() bool {
	u, err := url.Parse(strings.TrimSpace(s.DialURL()))
	if err != nil || u.Host == "" {
		return false
	}
	switch strings.ToLower(u.Hostname()) {
	case "127.0.0.1", "localhost", "::1":
		return true
	default:
		return false
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
