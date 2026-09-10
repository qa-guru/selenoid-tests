package playwrightapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHTTPURL(t *testing.T) {
	require.Equal(t, "http://h/playwright-chromium/1.62.1-min?enableHAR=true",
		HTTPURL("ws://h/playwright-chromium/1.62.1-min?enableHAR=true"))
	require.Equal(t, "https://h/p?accessKey=u:p", HTTPURL("wss://h/p?accessKey=u:p"))
	require.Equal(t, "http://h/p", HTTPURL("http://h/p"))
}

func TestHandshake_SendsUpgradeAndEnableHAR(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "websocket", strings.ToLower(r.Header.Get("Upgrade")))
		require.Equal(t, "true", r.URL.Query().Get("enableHAR"))
		require.True(t, strings.HasSuffix(r.URL.Path, "/playwright-chromium/1.62.1-min"))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("1.62.1-min is a headless CI image and does not support enableHAR"))
	}))
	defer srv.Close()

	ws := "ws://" + strings.TrimPrefix(srv.URL, "http://") + "/playwright/playwright-chromium/1.62.1-min"
	resp, err := Handshake(nil, AppendQuery(ws, "enableHAR=true"))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Contains(t, string(resp.Body), "enableHAR")
}
