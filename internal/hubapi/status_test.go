package hubapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/hubapi"
)

func TestParse_FlatHubStatus(t *testing.T) {
	st, err := hubapi.Parse([]byte(`{"total":10,"used":1,"queued":0,"pending":0,"browsers":{"chrome":{}}}`))
	require.NoError(t, err)
	require.Equal(t, 10, st.Total)
	require.Equal(t, 1, st.Used)
	require.NotNil(t, st.Browsers)
}

func TestParse_UIStateEnvelope(t *testing.T) {
	st, err := hubapi.Parse([]byte(`{"state":{"total":5,"used":0,"queued":0,"pending":0,"browsers":{}}}`))
	require.NoError(t, err)
	require.Equal(t, 5, st.Total)
	require.NotNil(t, st.Browsers)
}
