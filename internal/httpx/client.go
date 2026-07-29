package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response is a raw HTTP result (status + body).
type Response struct {
	StatusCode int
	Body       []byte
}

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
	resp, err := c.Do(http.MethodGet, pathOrURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, statusErr(http.MethodGet, pathOrURL, resp)
	}
	return resp.Body, nil
}

// GetText GETs and returns the body as a string (expects HTTP 200).
func (c *Client) GetText(pathOrURL string) (string, error) {
	body, err := c.GetBytes(pathOrURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetExpectStatus GETs and requires a specific status code.
func (c *Client) GetExpectStatus(pathOrURL string, expected int) (*Response, error) {
	resp, err := c.Do(http.MethodGet, pathOrURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expected {
		return resp, statusErr(http.MethodGet, pathOrURL, resp)
	}
	return resp, nil
}

// GetQuery GETs with query parameters (expects HTTP 200).
func (c *Client) GetQuery(pathOrURL string, query url.Values) ([]byte, error) {
	rawURL, err := c.resolve(pathOrURL)
	if err != nil {
		return nil, err
	}
	if len(query) > 0 {
		sep := "?"
		if strings.Contains(rawURL, "?") {
			sep = "&"
		}
		rawURL += sep + query.Encode()
	}
	resp, err := c.Do(http.MethodGet, rawURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, statusErr(http.MethodGet, pathOrURL, resp)
	}
	return resp.Body, nil
}

// PostJSON POSTs JSON and returns the response (checks expectedStatus when > 0).
func (c *Client) PostJSON(pathOrURL string, payload any, expectedStatus int) (*Response, error) {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := c.Do(http.MethodPost, pathOrURL, body, headers)
	if err != nil {
		return nil, err
	}
	if expectedStatus > 0 && resp.StatusCode != expectedStatus {
		return resp, statusErr(http.MethodPost, pathOrURL, resp)
	}
	return resp, nil
}

// Delete sends DELETE and checks status when expectedStatus > 0.
func (c *Client) Delete(pathOrURL string, expectedStatus int) (*Response, error) {
	resp, err := c.Do(http.MethodDelete, pathOrURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if expectedStatus > 0 && resp.StatusCode != expectedStatus {
		return resp, statusErr(http.MethodDelete, pathOrURL, resp)
	}
	return resp, nil
}

// Do executes an HTTP request against base URL or absolute path.
func (c *Client) Do(method, pathOrURL string, body []byte, headers map[string]string) (*Response, error) {
	rawURL, err := c.resolve(pathOrURL)
	if err != nil {
		return nil, err
	}
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, stripUserInfo(rawURL), reader)
	if err != nil {
		return nil, err
	}
	if user, pass, ok := basicAuth(rawURL); ok {
		req.SetBasicAuth(user, pass)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{StatusCode: resp.StatusCode, Body: respBody}, nil
}

func statusErr(method, path string, resp *Response) error {
	return fmt.Errorf("%s %s: HTTP %d: %s", method, path, resp.StatusCode, truncate(string(resp.Body), 300))
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
