package api_test

import (
	"strings"
	"testing"
)

func assertBodyContains(t *testing.T, body []byte, substr string) {
	t.Helper()
	if !strings.Contains(string(body), substr) {
		t.Fatalf("expected body to contain %q, got %s", substr, truncate(string(body), 200))
	}
}

func assertValidMp4(t *testing.T, body []byte, fileName string) {
	t.Helper()
	if len(body) <= 1024 {
		t.Fatalf("expected non-trivial mp4 body for %s, got %d bytes", fileName, len(body))
	}
	if len(body) < 8 || body[4] != 'f' || body[5] != 't' || body[6] != 'y' || body[7] != 'p' {
		t.Fatalf("expected ISO BMFF ftyp box in %s", fileName)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
