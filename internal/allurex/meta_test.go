package allurex

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_writesBrowserTagAndParameter(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ALLURE_RESULTS", dir)

	Run(t, Meta{
		Name:      "browser meta",
		Package:   "allurex.MetaTests",
		Layer:     "unit",
		Component: "webdriver-image",
		Epic:      "webdriver-image",
		Browser:   BrowserFirefox,
		Tags:      []string{"positive", "firefox"},
	}, func(a *A) {
		a.Step("record result", func() {})
	})

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 result, got %d", len(entries))
	}
	raw, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	var got result
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if !hasLabel(got.Labels, "tag", BrowserFirefox) {
		t.Fatalf("missing firefox tag: %s", raw)
	}
	tagCount := 0
	for _, l := range got.Labels {
		if l.Name == "tag" && l.Value == BrowserFirefox {
			tagCount++
		}
	}
	if tagCount != 1 {
		t.Fatalf("duplicate firefox tags: %d in %s", tagCount, raw)
	}
	if !hasLabel(got.Labels, "epic", "webdriver-image") {
		t.Fatalf("missing epic: %s", raw)
	}
	foundBrowser := false
	for _, p := range got.Parameters {
		m, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if m["name"] == "browser" && m["value"] == BrowserFirefox {
			foundBrowser = true
		}
	}
	if !foundBrowser {
		t.Fatalf("missing browser parameter: %s", raw)
	}
}

func hasLabel(labels []label, name, value string) bool {
	for _, l := range labels {
		if l.Name == name && l.Value == value {
			return true
		}
	}
	return false
}
