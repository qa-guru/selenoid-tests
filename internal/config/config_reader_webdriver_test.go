package config_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReaderWebdriver_UiBrowserURLUsesUiWhenRemoteBlank(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveUiBrowserUrl uses uiUrl when remote blank"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":     "http://127.0.0.1:8080/",
				"remoteUrl": "",
			})
			got, err := cfg.ResolveUiBrowserURL()
			require.NoError(t, err)
			require.Equal(t, "http://127.0.0.1:8080", got)
		})
	})
}

func TestConfigReaderWebdriver_UiBrowserURLMapsLoopbackForRemote(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveUiBrowserUrl maps loopback for remote"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":     "http://127.0.0.1:8080",
				"remoteUrl": "http://127.0.0.1:4444/wd/hub",
			})
			got, err := cfg.ResolveUiBrowserURL()
			require.NoError(t, err)
			require.Contains(t, got, "host.docker.internal")
			require.True(t, strings.HasSuffix(got, ":8080"))
		})
	})
}

func TestConfigReaderWebdriver_UiBrowserURLMapsLocalhostForRemote(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveUiBrowserUrl maps localhost for remote"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":     "http://localhost:8080/",
				"remoteUrl": "http://127.0.0.1:4444/wd/hub",
			})
			got, err := cfg.ResolveUiBrowserURL()
			require.NoError(t, err)
			require.Equal(t, "http://host.docker.internal:8080", got)
		})
	})
}

func TestConfigReaderWebdriver_UiBrowserURLFailsWhenUiEmptyWithRemote(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveUiBrowserUrl fails when uiUrl empty"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":     "",
				"remoteUrl": "http://127.0.0.1:4444/wd/hub",
			})
			_, err := cfg.ResolveUiBrowserURL()
			require.Error(t, err)
			require.Contains(t, err.Error(), "uiUrl")
		})
	})
}

func TestConfigReaderWebdriver_UiBrowserURLStripsBasicAuth(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveUiBrowserUrl strips embedded basic auth"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":     "https://user1:1234@selenoid.qa.guru/",
				"remoteUrl": "https://user1:1234@selenoid.qa.guru/wd/hub",
			})
			got, err := cfg.ResolveUiBrowserURL()
			require.NoError(t, err)
			require.Equal(t, "https://selenoid.qa.guru", got)
		})
	})
}

func TestConfigReaderWebdriver_HubBasicAuthFromRemoteURL(t *testing.T) {
	allurex.Run(t, unitMeta("webdriver-image", "resolveHubBasicAuth from remoteUrl"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"uiUrl":      "https://selenoid.qa.guru/",
				"remoteUrl":  "https://user1:secret@selenoid.qa.guru/wd/hub",
				"apiBaseUrl": "",
			})
			user, pass := cfg.ResolveHubBasicAuth()
			require.Equal(t, "user1", user)
			require.Equal(t, "secret", pass)
		})
	})
}
