package fixture_test

import (
	"os"
	"path/filepath"
	"testing"
)

// loadFixture reads src/test/resources/fixtures/<rel> (Java classpath parity).
func loadProjectFixture(t *testing.T, rel string) []byte {
	t.Helper()
	path := filepath.Join(moduleRoot(t), filepath.FromSlash(rel))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture %s: %v", rel, err)
	}
	return body
}

func loadFixture(t *testing.T, rel string) []byte {
	t.Helper()
	path := filepath.Join(moduleRoot(t), "src", "test", "resources", "fixtures", filepath.FromSlash(rel))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture %s: %v", rel, err)
	}
	return body
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found from test cwd")
		}
		dir = parent
	}
}
