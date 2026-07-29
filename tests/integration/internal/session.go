// Package integration provides shared helpers for integration tests (WD + UI slices).
package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// AssertBrowserVersionListed checks hub/UI browsers map contains family+version.
func AssertBrowserVersionListed(t *testing.T, browsers map[string]any, family, version string) {
	t.Helper()
	raw, ok := browsers[family]
	require.True(t, ok, "expected %q in browsers map", family)
	versions, ok := raw.(map[string]any)
	require.True(t, ok, "expected %q browsers entry to be a version map", family)
	require.NotNil(t, versions[version], "expected version %q listed for %q", version, family)
}
