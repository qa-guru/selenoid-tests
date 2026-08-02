package allurex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDisplayName_subtestAfterSlash(t *testing.T) {
	got := ResolveDisplayName(
		"TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body",
		"github.com/qa-guru/selenoid.TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body",
	)
	if got != "POST /wd/hub/session rejects malformed JSON body" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveDisplayName_plainName(t *testing.T) {
	got := ResolveDisplayName("TestBrowserNotFound", "github.com/qa-guru/selenoid.TestBrowserNotFound")
	if got != "TestBrowserNotFound" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveDisplayName_fromFullNameOnly(t *testing.T) {
	got := ResolveDisplayName("", "github.com/qa-guru/selenoid.TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body")
	if got != "POST /wd/hub/session rejects malformed JSON body" {
		t.Fatalf("got %q", got)
	}
}

func TestConvertGoTestJSON_writesResults(t *testing.T) {
	dir := t.TempDir()
	input := strings.Join([]string{
		`{"Time":"2026-08-02T12:00:00Z","Action":"run","Package":"github.com/qa-guru/selenoid","Test":"TestFoo"}`,
		`{"Time":"2026-08-02T12:00:00.1Z","Action":"output","Package":"github.com/qa-guru/selenoid","Test":"TestFoo","Output":"hello\n"}`,
		`{"Time":"2026-08-02T12:00:01Z","Action":"pass","Package":"github.com/qa-guru/selenoid","Test":"TestFoo","Elapsed":1.0}`,
		`{"Time":"2026-08-02T12:00:01Z","Action":"run","Package":"github.com/qa-guru/selenoid","Test":"TestBar/case one"}`,
		`{"Time":"2026-08-02T12:00:02Z","Action":"fail","Package":"github.com/qa-guru/selenoid","Test":"TestBar/case one","Elapsed":0.5}`,
	}, "\n") + "\n"

	n, err := ConvertGoTestJSON(strings.NewReader(input), ConvertOptions{
		Epic:      "selenoid",
		Component: "selenoid",
		Layer:     "unit",
		OutputDir: dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("want 2 results, got %d", n)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 files, got %d", len(entries))
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), "-result.json") {
			t.Fatalf("unexpected file %s", e.Name())
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		s := string(raw)
		if !strings.Contains(s, `"name":"language"`) || !strings.Contains(s, `"value":"go"`) {
			t.Fatalf("missing language label in %s: %s", e.Name(), s)
		}
		if !strings.Contains(s, `"name":"epic"`) || !strings.Contains(s, `"value":"selenoid"`) {
			t.Fatalf("missing epic label in %s: %s", e.Name(), s)
		}
	}
}
