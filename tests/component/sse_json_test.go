package component_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/uiapi"
)

func TestSseStateJson_ParsesHubState(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses SSE payload with hub state",
		Package:   "tests.component.SseStateJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "SSE state fixture",
		Suite:     "SSE state fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/sse/state.json", func() {
			ev, err := uiapi.ParseSseEvent(loadFixture(t, "sse/state.json"))
			require.NoError(t, err)
			require.True(t, ev.HasState())
		})
	})
}

func TestSseErrorsJson_ParsesErrorsList(t *testing.T) {
	allurex.Run(t, allurex.Meta{
		Name:      "parses SSE payload with errors list",
		Package:   "tests.component.SseErrorsJsonTest",
		Layer:     "component",
		Component: "selenoid-ui",
		Feature:   "SSE errors fixture",
		Suite:     "SSE errors fixture",
		Tags:      []string{"component"},
	}, func(a *allurex.A) {
		a.Step("parse fixtures/sse/errors.json", func() {
			ev, err := uiapi.ParseSseEvent(loadFixture(t, "sse/errors.json"))
			require.NoError(t, err)
			require.True(t, ev.HasErrors())
		})
	})
}
