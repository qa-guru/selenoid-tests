package allurex

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// testEvent is one line from `go test -json`.
type testEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"` // seconds
	Output  string    `json:"Output"`
}

type pendingTest struct {
	pkg     string
	test    string
	start   time.Time
	outputs []string
}

// ConvertOptions controls labels written for foreign product `go test` runs.
type ConvertOptions struct {
	Epic      string
	Component string
	Layer     string
	OutputDir string
}

// ResolveDisplayName turns Go test / t.Run names into readable Allure titles.
func ResolveDisplayName(name, fullName string) string {
	if title := subtestTitle(name); title != "" {
		return title
	}
	if name != "" {
		return name
	}
	if title := subtestTitleFromFullName(fullName); title != "" {
		return title
	}
	if fullName != "" {
		return fullName
	}
	return "unnamed test"
}

// subtestTitle returns the part after the first '/' in a Go test name (t.Run).
func subtestTitle(name string) string {
	slash := strings.IndexByte(name, '/')
	if slash < 0 || slash >= len(name)-1 {
		return ""
	}
	return strings.TrimSpace(name[slash+1:])
}

// subtestTitleFromFullName finds ".Test…/title" — ignores '/' inside module paths.
func subtestTitleFromFullName(fullName string) string {
	idx := strings.Index(fullName, ".Test")
	if idx < 0 {
		return ""
	}
	rest := fullName[idx+1:] // Test…
	slash := strings.IndexByte(rest, '/')
	if slash < 0 || slash >= len(rest)-1 {
		return ""
	}
	return strings.TrimSpace(rest[slash+1:])
}

// ConvertGoTestJSON reads `go test -json` events and writes Allure *-result.json files.
// Returns the number of results written.
func ConvertGoTestJSON(r io.Reader, opt ConvertOptions) (int, error) {
	if opt.OutputDir == "" {
		return 0, fmt.Errorf("output dir required")
	}
	if opt.Epic == "" {
		return 0, fmt.Errorf("epic required")
	}
	if opt.Component == "" {
		opt.Component = opt.Epic
	}
	if opt.Layer == "" {
		opt.Layer = "unit"
	}
	if err := os.MkdirAll(opt.OutputDir, 0o755); err != nil {
		return 0, err
	}

	pending := map[string]*pendingTest{}
	count := 0
	sc := bufio.NewScanner(r)
	// go test output lines can be large (logs); raise limit beyond 64K default.
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for sc.Scan() {
		line := sc.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var ev testEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			// Non-JSON noise on stdout — ignore.
			continue
		}
		if ev.Test == "" {
			continue
		}
		key := ev.Package + "\x00" + ev.Test
		switch ev.Action {
		case "run":
			pending[key] = &pendingTest{
				pkg:   ev.Package,
				test:  ev.Test,
				start: eventTime(ev),
			}
		case "output":
			if p := pending[key]; p != nil && ev.Output != "" {
				p.outputs = append(p.outputs, ev.Output)
			}
		case "pass", "fail", "skip":
			p := pending[key]
			if p == nil {
				p = &pendingTest{
					pkg:   ev.Package,
					test:  ev.Test,
					start: eventTime(ev).Add(-elapsedDuration(ev.Elapsed)),
				}
			}
			if ev.Output != "" {
				p.outputs = append(p.outputs, ev.Output)
			}
			status := mapGoAction(ev.Action)
			stop := eventTime(ev)
			start := p.start
			if ev.Elapsed > 0 {
				start = stop.Add(-elapsedDuration(ev.Elapsed))
			}
			if err := writeConvertedResult(opt, p.pkg, p.test, status, start, stop, strings.Join(p.outputs, "")); err != nil {
				return count, err
			}
			count++
			delete(pending, key)
		}
	}
	if err := sc.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func writeConvertedResult(opt ConvertOptions, pkg, testName, status string, start, stop time.Time, output string) error {
	fullName := pkg + "." + testName
	display := ResolveDisplayName(testName, fullName)
	id := newUUID()
	payload := result{
		UUID:      id,
		HistoryID: historyID(fullName, display),
		Name:      display,
		FullName:  fullName,
		Status:    status,
		Stage:     "finished",
		Start:     start.UnixMilli(),
		Stop:      stop.UnixMilli(),
		Steps:     nil,
		Labels: []label{
			{Name: "language", Value: "go"},
			{Name: "framework", Value: "go"},
			{Name: "host", Value: hostname()},
			{Name: "thread", Value: fmt.Sprintf("%d", os.Getpid())},
			{Name: "package", Value: pkg},
			{Name: "testClass", Value: pkg},
			{Name: "testMethod", Value: testName},
			{Name: "suite", Value: pkg},
			{Name: "layer", Value: opt.Layer},
			{Name: "component", Value: opt.Component},
			{Name: "epic", Value: opt.Epic},
		},
		Attachments: []any{},
		Parameters:  []any{},
		Links:       []any{},
	}
	if status == "failed" || status == "broken" {
		msg := strings.TrimSpace(output)
		if msg == "" {
			msg = display + " " + status
		}
		// Keep statusDetails compact for TestOps.
		if len(msg) > 8000 {
			msg = msg[:8000] + "\n…(truncated)"
		}
		payload.StatusDetails = map[string]string{"message": msg}
	}
	path := filepath.Join(opt.OutputDir, id+"-result.json")
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func mapGoAction(action string) string {
	switch action {
	case "pass":
		return "passed"
	case "fail":
		return "failed"
	case "skip":
		return "skipped"
	default:
		return "unknown"
	}
}

func eventTime(ev testEvent) time.Time {
	if !ev.Time.IsZero() {
		return ev.Time
	}
	return time.Now()
}

func elapsedDuration(sec float64) time.Duration {
	return time.Duration(sec * float64(time.Second))
}
