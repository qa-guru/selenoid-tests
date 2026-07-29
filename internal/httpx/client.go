package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a small JSON HTTP helper with optional basic auth from URLs.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

// New builds a client. baseURL may embed user:pass@ (prod profiles).
func New(baseURL string) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
	}
}

// GetJSON GETs path (or absolute URL) and decodes JSON into dest. Expects HTTP 200.
func (c *Client) GetJSON(pathOrURL string, dest any) error {
	rawURL, err := c.resolve(pathOrURL)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, stripUserInfo(rawURL), nil)
	if err != nil {
		return err
	}
	if user, pass, ok := basicAuth(rawURL); ok {
		req.SetBasicAuth(user, pass)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: HTTP %d: %s", stripUserInfo(rawURL), resp.StatusCode, truncate(string(body), 300))
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode %s: %w", stripUserInfo(rawURL), err)
	}
	return nil
}

// GetBytes GETs and returns body for custom parsing.
func (c *Client) GetBytes(pathOrURL string) ([]byte, error) {
	rawURL, err := c.resolve(pathOrURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, stripUserInfo(rawURL), nil)
	if err != nil {
		return nil, err
	}
	if user, pass, ok := basicAuth(rawURL); ok {
		req.SetBasicAuth(user, pass)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", stripUserInfo(rawURL), resp.StatusCode, truncate(string(body), 300))
	}
	return body, nil
}

func (c *Client) resolve(pathOrURL string) (string, error) {
	pathOrURL = strings.TrimSpace(pathOrURL)
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		return pathOrURL, nil
	}
	if c.BaseURL == "" {
		return "", fmt.Errorf("empty base URL")
	}
	if !strings.HasPrefix(pathOrURL, "/") {
		pathOrURL = "/" + pathOrURL
	}
	return c.BaseURL + pathOrURL, nil
}

func basicAuth(raw string) (user, pass string, ok bool) {
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		return "", "", false
	}
	pass, _ = u.User.Password()
	return u.User.Username(), pass, u.User.Username() != ""
}

func stripUserInfo(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.User = nil
	return u.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
