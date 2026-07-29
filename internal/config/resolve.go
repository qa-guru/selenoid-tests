package config

import (
	"fmt"
	"net/url"
	"strings"
)

// ResolveHubURL trims hubUrl and ensures a trailing slash (ConfigReader.resolveHubUrl).
func (c *Config) ResolveHubURL() (string, error) {
	u := strings.TrimSpace(c.HubURL)
	if u == "" {
		return "", fmt.Errorf("Set hubUrl in config/${env}.properties")
	}
	return withSlash(u), nil
}

// ResolveUiURL trims uiUrl and ensures a trailing slash (ConfigReader.resolveUiUrl).
func (c *Config) ResolveUiURL() (string, error) {
	u := strings.TrimSpace(c.UIURL)
	if u == "" {
		return "", fmt.Errorf("Set uiUrl in config/${env}.properties")
	}
	return withSlash(u), nil
}

// ResolveAPIBaseURL prefers apiBaseUrl, else hubUrl (ConfigReader.resolveApiBaseUrl).
func (c *Config) ResolveAPIBaseURL() (string, error) {
	api := strings.TrimSpace(c.APIBaseURL)
	if api != "" {
		return withSlash(api), nil
	}
	hub := strings.TrimSpace(c.HubURL)
	if hub != "" {
		return withSlash(hub), nil
	}
	return "", fmt.Errorf("Set apiBaseUrl or hubUrl in config/${env}.properties")
}

// ResolveHubStatusPath normalizes hubStatusPath (default /status).
func (c *Config) ResolveHubStatusPath() string {
	return normalizePath(c.HubStatusPath)
}

// ResolveUiBrowserURL maps loopback for remote sessions and strips userinfo
// (ConfigReader.resolveUiBrowserUrl).
func (c *Config) ResolveUiBrowserURL() (string, error) {
	remote := strings.TrimSpace(c.RemoteURL)
	if remote == "" {
		ui, err := c.ResolveUiURL()
		if err != nil {
			return "", err
		}
		return stripUserInfo(stripTrailingSlash(ui)), nil
	}
	ui := strings.TrimSpace(c.UIURL)
	if ui == "" {
		return "", fmt.Errorf("Set uiUrl in config/${env}.properties")
	}
	browserURL := strings.ReplaceAll(ui, "127.0.0.1", "host.docker.internal")
	browserURL = strings.ReplaceAll(browserURL, "localhost", "host.docker.internal")
	return stripUserInfo(stripTrailingSlash(withSlash(browserURL))), nil
}

// ResolveHubBasicAuth reads user:pass from remoteUrl / apiBaseUrl / uiUrl.
func (c *Config) ResolveHubBasicAuth() (user, pass string) {
	for _, candidate := range []string{c.RemoteURL, c.APIBaseURL, c.UIURL} {
		if u, p, ok := parseUserInfo(candidate); ok {
			return u, p
		}
	}
	return "", ""
}

// ResolvePlaywrightWsEndpoint appends session query params when missing
// (ConfigReader.resolvePlaywrightWsEndpoint).
func (c *Config) ResolvePlaywrightWsEndpoint() (string, error) {
	base := strings.TrimSpace(c.PlaywrightWsEndpoint)
	if base == "" {
		return "", fmt.Errorf("Set playwrightWsEndpoint in config/${env}.properties")
	}
	if strings.Contains(base, "?") {
		return base, nil
	}
	q := url.Values{}
	q.Set("name", c.PlaywrightSessionName)
	q.Set("sessionTimeout", c.PlaywrightSessionTimeout)
	q.Set("enableVNC", boolString(c.PlaywrightEnableVNC))
	q.Set("enableVideo", boolString(c.PlaywrightEnableVideo))
	return base + "?" + q.Encode(), nil
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func stripTrailingSlash(value string) string {
	if strings.HasSuffix(value, "/") {
		return value[:len(value)-1]
	}
	return value
}

// stripUserInfo drops user:pass@ so fetch() can create sessions.
func stripUserInfo(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}
	schemeSep := strings.Index(raw, "://")
	if schemeSep < 0 {
		return raw
	}
	afterScheme := schemeSep + 3
	slash := strings.Index(raw[afterScheme:], "/")
	authorityEnd := len(raw)
	if slash >= 0 {
		authorityEnd = afterScheme + slash
	}
	authority := raw[afterScheme:authorityEnd]
	at := strings.LastIndex(authority, "@")
	if at < 0 {
		return raw
	}
	return raw[:afterScheme] + authority[at+1:] + raw[authorityEnd:]
}

func parseUserInfo(raw string) (user, pass string, ok bool) {
	if strings.TrimSpace(raw) == "" {
		return "", "", false
	}
	schemeSep := strings.Index(raw, "://")
	if schemeSep < 0 {
		return "", "", false
	}
	afterScheme := schemeSep + 3
	slash := strings.Index(raw[afterScheme:], "/")
	authorityEnd := len(raw)
	if slash >= 0 {
		authorityEnd = afterScheme + slash
	}
	authority := raw[afterScheme:authorityEnd]
	at := strings.LastIndex(authority, "@")
	if at <= 0 {
		return "", "", false
	}
	userInfo := authority[:at]
	colon := strings.Index(userInfo, ":")
	if colon <= 0 {
		return "", "", false
	}
	return userInfo[:colon], userInfo[colon+1:], true
}
