package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

func TestLoad_LocalUnitDefaults(t *testing.T) {
	config.ResetForTest()
	t.Setenv("SELENOID_TEST_ENV", "local_unit")
	_ = os.Unsetenv("hubUrl")
	_ = os.Unsetenv("SELENOID_TEST_HUB_URL")

	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "local_unit", cfg.Env)
	require.Equal(t, "http://127.0.0.1:4444/", cfg.HubURL)
	require.Equal(t, "http://127.0.0.1:8080/", cfg.UIURL)
	require.Equal(t, "/status", cfg.HubStatusPath)
}

func TestLoad_ProdApiHubStatusPath(t *testing.T) {
	config.ResetForTest()
	t.Setenv("SELENOID_TEST_ENV", "selenoid_qa_guru_api")

	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, "/hub/status", cfg.HubStatusPath)
	require.Contains(t, cfg.APIBase(), "selenoid.qa.guru")
}
