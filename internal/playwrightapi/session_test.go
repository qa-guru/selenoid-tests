package playwrightapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendQuery(t *testing.T) {
	require.Equal(t, "ws://h/playwright-chromium/1.62.1-min?enableHAR=true",
		AppendQuery("ws://h/playwright-chromium/1.62.1-min", "enableHAR=true"))
	require.Equal(t, "wss://h/p?accessKey=u:p&enableHAR=true&enableVideo=false",
		AppendQuery("wss://h/p?accessKey=u:p", "enableHAR=true&enableVideo=false"))
	require.Equal(t, "ws://h/p", AppendQuery("ws://h/p", "  "))
}
