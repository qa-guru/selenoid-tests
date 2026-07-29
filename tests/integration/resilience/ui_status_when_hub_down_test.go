package resilience_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/config"
	"github.com/qa-guru/selenoid-tests/internal/stack"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestUiStatusWhenHubDown_UiStatusReportsHubError(t *testing.T) {
	cfg := config.MustLoad()
	allurex.Run(t, allurex.Meta{
		Name:      "UI /status returns errors when hub is unavailable",
		Package:   "tests.integration.UiStatusWhenHubDownTests",
		Layer:     "integration",
		Component: "selenoid-ui",
		Epic:      "selenoid-ui",
		Feature:   "UI status when hub down",
		Story:     "UI status when hub down",
		Suite:     "UI status when hub down",
		Tags:      []string{"integration", "local-only", "positive"},
	}, func(a *allurex.A) {
		a.Step("Ensure controllable ci-bin stack", func() {
			require.NoError(t, stack.EnsureControllable(cfg))
		})

		a.Step("Stop hub", func() {
			require.NoError(t, stack.KillHub(cfg))
			require.NoError(t, stack.WaitForHubDown(cfg, 15*time.Second))
		})

		defer func() {
			a.Step("Restart hub", func() {
				require.NoError(t, stack.StartHubDetached())
				require.NoError(t, stack.WaitForHubReady(cfg, 30*time.Second))
			})
		}()

		a.Step("GET UI /status", func() {
			resp, err := uiapi.FetchStatusWhenHubUnavailable(cfg)
			require.NoError(t, err)
			require.NotNil(t, resp.Errors)
			require.NotEmpty(t, resp.Errors)
			require.NotEmpty(t, resp.Errors[0].Msg)
		})
	})
}
