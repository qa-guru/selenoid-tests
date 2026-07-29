package component_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestHubStatusJson_ParsesIdleCounters(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses idle hub status counters",
		Package:   "tests.component.HubStatusJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Hub status fixture",
		Suite:     "Hub status fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/status-idle.json", func() {
			st, err := hubapi.Parse(loadFixture(t, "hub/status-idle.json"))
			require.NoError(t, err)
			require.Equal(t, 5, st.Total)
			require.Equal(t, 0, st.Used)
			require.NotNil(t, st.Browsers)
		})
	})
}

func TestHubStatusParser_FlatAndUIWrapped(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses flat hub /status JSON",
		Package:   "tests.component.HubStatusParserTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Hub status parser",
		Suite:     "Hub status parser",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse flat hub status", func() {
			st, err := hubapi.Parse(loadFixture(t, "hub/status-idle.json"))
			require.NoError(t, err)
			require.Equal(t, 0, st.Used)
			require.Equal(t, 5, st.Total)
		})
		a.Step("parse UI-shaped /status with .state wrapper", func() {
			st, err := hubapi.Parse(loadFixture(t, "ui/status.json"))
			require.NoError(t, err)
			require.Equal(t, 1, st.Used)
			require.Equal(t, 5, st.Total)
		})
	})
}

func TestHubStatusBrowsersJson_ParsesBrowsersMap(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses browsers map with version entries",
		Package:   "tests.component.HubStatusBrowsersJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Hub status browsers fixture",
		Suite:     "Hub status browsers fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/status-with-browsers.json", func() {
			st, err := hubapi.Parse(loadFixture(t, "hub/status-with-browsers.json"))
			require.NoError(t, err)
			_, ok := st.Browsers["chrome"]
			require.True(t, ok)
			require.NotEmpty(t, st.Browsers)
		})
	})
}

func TestHubStatusForwardCompatJson_IgnoresUnknownFields(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "ignores unknown JSON fields",
		Package:   "tests.component.HubStatusForwardCompatJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Hub status forward compatibility",
		Suite:     "Hub status forward compatibility",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/status-unknown-field.json", func() {
			st, err := hubapi.Parse(loadFixture(t, "hub/status-unknown-field.json"))
			require.NoError(t, err)
			require.Equal(t, 3, st.Total)
		})
	})
}

func TestHubLogsListJson_ParsesSessionLogNames(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses hub logs JSON array with session file names",
		Package:   "tests.component.HubLogsListJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Epic:      "selenoid",
		Feature:   "Hub logs list fixture",
		Suite:     "Hub logs list fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/logs-list.json", func() {
			files, err := hubapi.ParseLogsList(loadFixture(t, "hub/logs-list.json"))
			require.NoError(t, err)
			require.NotEmpty(t, files)
			found := false
			for _, name := range files {
				if strings.HasSuffix(name, ".log") {
					found = true
					break
				}
			}
			require.True(t, found)
		})
	})
}

func TestHubWebDriverStatusJson_ParsesReadyPayload(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses ready WebDriver status payload",
		Package:   "tests.component.HubWebDriverStatusJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Hub WebDriver status fixture",
		Suite:     "Hub WebDriver status fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/wd-status-ready.json", func() {
			st, err := hubapi.ParseWebDriverStatus(loadFixture(t, "hub/wd-status-ready.json"))
			require.NoError(t, err)
			require.True(t, st.Value.Ready)
			require.Equal(t, "Selenoid ready", st.Value.Message)
		})
	})
}

func TestSessionCreateJson_ExtractsSessionID(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "extracts sessionId from create session response",
		Package:   "tests.component.SessionCreateJsonTest",
		Layer:     "component",
		Component: "selenoid",
		Feature:   "Session create fixture",
		Suite:     "Session create fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/hub/session-create.json", func() {
			id, err := hubapi.ParseSessionID(loadFixture(t, "hub/session-create.json"))
			require.NoError(t, err)
			require.Equal(t, "abc123-session", id)
		})
	})
}
