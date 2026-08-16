package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatus_ListsAndroidWhenConfigured(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /status lists android when the Appium image is configured",
		Package:   "tests.api.HubAndroidStatusTests",
		Layer:     "api",
		Component: "android",
		Epic:      "android",
		Feature:   "Appium session",
		Story:     "Android hub status",
		Suite:     "Android hub status",
		Browser:   allurex.BrowserAndroid,
		Tags:      []string{"api", "positive"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET /status", func() {
			var err error
			status, err = hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		a.Step("Skip unless this hub advertises android", func() {
			if status.Browsers == nil {
				t.Skip("hub /status has no browsers map")
			}
			if _, ok := status.Browsers["android"]; !ok {
				t.Skip("android not configured on this hub — live Appium session smoke stays behind this skip")
			}
		})
		a.Step("Verify android family is present", func() {
			require.Contains(t, status.Browsers, "android")
		})
	})
}

func TestHubStatus_IosFamilyNotConfigured(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "GET /status has no ios family (roadmap stub)",
		Package:   "tests.api.HubIosStatusTests",
		Layer:     "api",
		Component: "ios",
		Epic:      "ios",
		Feature:   "Appium session",
		Story:     "iOS hub status stub",
		Suite:     "iOS hub status stub",
		Browser:   allurex.BrowserIOS,
		Tags:      []string{"api", "future"},
	}, func(a *allurex.A) {
		var status *hubapi.Status
		a.Step("GET /status", func() {
			var err error
			status, err = hubapi.Fetch(cfg)
			require.NoError(t, err)
		})
		a.Step("ios is not advertised until qaguru/ios ships", func() {
			_, ok := status.Browsers["ios"]
			require.False(t, ok, "replace this stub with Appium session smoke when the hub lists ios")
		})
	})
}
