package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestConfigReaderPlaywright_ResolveAddsQueryParams(t *testing.T) {
	allurex.Run(t, unitMeta("playwright-image", "resolvePlaywrightWsEndpoint appends query params"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"playwrightWsEndpoint":     "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
				"playwrightSessionName":    "smoke",
				"playwrightSessionTimeout": "5m",
			})
			endpoint, err := cfg.ResolvePlaywrightWsEndpoint()
			require.NoError(t, err)
			require.Contains(t, endpoint, "name=smoke")
			require.Contains(t, endpoint, "sessionTimeout=5m")
			require.Contains(t, endpoint, "enableVNC=false")
			require.Contains(t, endpoint, "enableVideo=false")
		})
	})
}

func TestConfigReaderPlaywright_ResolveFailsWhenEmpty(t *testing.T) {
	allurex.Run(t, unitMeta("playwright-image", "resolvePlaywrightWsEndpoint fails when empty"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{"playwrightWsEndpoint": ""})
			_, err := cfg.ResolvePlaywrightWsEndpoint()
			require.Error(t, err)
			require.Contains(t, err.Error(), "playwrightWsEndpoint")
		})
	})
}

func TestConfigReaderPlaywright_KeepsExistingQueryString(t *testing.T) {
	allurex.Run(t, unitMeta("playwright-image", "resolvePlaywrightWsEndpoint keeps existing query"), func(a *allurex.A) {
		a.Step("resolve", func() {
			preset := "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1?name=preset&enableVNC=true"
			cfg := config.FromMap(map[string]string{
				"playwrightWsEndpoint":  preset,
				"playwrightSessionName": "ignored",
				"playwrightEnableVnc":   "false",
			})
			got, err := cfg.ResolvePlaywrightWsEndpoint()
			require.NoError(t, err)
			require.Equal(t, preset, got)
		})
	})
}

func TestConfigReaderPlaywright_AppendsVncAndVideoFlags(t *testing.T) {
	allurex.Run(t, unitMeta("playwright-image", "resolvePlaywrightWsEndpoint appends VNC/video"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"playwrightWsEndpoint":     "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
				"playwrightSessionName":    "rec",
				"playwrightSessionTimeout": "10m",
				"playwrightEnableVnc":      "true",
				"playwrightEnableVideo":    "true",
			})
			endpoint, err := cfg.ResolvePlaywrightWsEndpoint()
			require.NoError(t, err)
			require.Contains(t, endpoint, "enableVNC=true")
			require.Contains(t, endpoint, "enableVideo=true")
			require.Contains(t, endpoint, "sessionTimeout=10m")
		})
	})
}

func TestConfigReaderPlaywright_URLEncodesSessionName(t *testing.T) {
	allurex.Run(t, unitMeta("playwright-image", "resolvePlaywrightWsEndpoint URL-encodes name"), func(a *allurex.A) {
		a.Step("resolve", func() {
			cfg := config.FromMap(map[string]string{
				"playwrightWsEndpoint":  "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
				"playwrightSessionName": "my session",
			})
			endpoint, err := cfg.ResolvePlaywrightWsEndpoint()
			require.NoError(t, err)
			require.Contains(t, endpoint, "name=my+session")
		})
	})
}
