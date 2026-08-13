package warmpool_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/warmpool"
)

const testOwner = "selenoid-tests-api"

func liveClient(t *testing.T) *warmpool.Client {
	t.Helper()
	cli := warmpool.Default()
	if err := cli.Ping(); err != nil {
		t.Skipf("warm-pool stand down (%s): %v", warmpool.BaseURL(), err)
	}
	return cli
}

func releaseOwned(t *testing.T, cli *warmpool.Client, owner string) {
	t.Helper()
	slots, err := cli.Slots()
	if err != nil {
		return
	}
	for _, s := range slots {
		if s.ReservedBy != nil && *s.ReservedBy == owner {
			_ = cli.Release(s.ID)
		}
	}
}

func requireAPIError(t *testing.T, status, want int, body []byte, msg string) {
	t.Helper()
	require.Equal(t, want, status)
	require.Equal(t, msg, warmpool.ParseError(body))
}

func anySlotID(t *testing.T, cli *warmpool.Client) string {
	t.Helper()
	slots, err := cli.Slots()
	require.NoError(t, err)
	require.NotEmpty(t, slots)
	return slots[0].ID
}
