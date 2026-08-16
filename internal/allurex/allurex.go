// Package allurex writes Allure 2 *-result.json for Go tests (TestOps 5271 labels).
package allurex

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// Browser flavors for Meta.Browser (Allure tag + parameter). Epic stays on the image family.
const (
	BrowserChrome   = "chrome"
	BrowserFirefox  = "firefox"
	BrowserMsedge   = "msedge"
	BrowserChromium = "chromium"
	BrowserWebkit   = "webkit"
	BrowserAndroid  = "android"
	BrowserIOS      = "ios"
)

// Parameter is an Allure test parameter (shown on the result card).
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Meta maps Java @Layer / @Component / Epic-Feature-Story / @Tag.
type Meta struct {
	Name       string
	FullName   string
	Layer      string
	Component  string
	Epic       string
	Feature    string
	Story      string
	Tags       []string
	Suite      string
	Package    string
	Browser    string
	Parameters []Parameter
}

// A is a per-test Allure context (steps + result file).
type A struct {
	t     *testing.T
	meta  Meta
	start time.Time
	mu    sync.Mutex
	steps []step
}

type step struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	Stage         string `json:"stage"`
	Start         int64  `json:"start"`
	Stop          int64  `json:"stop"`
	Steps         []any  `json:"steps"`
	Attachments   []any  `json:"attachments"`
	Parameters    []any  `json:"parameters"`
	StatusDetails any    `json:"statusDetails,omitempty"`
}

type label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type result struct {
	UUID          string  `json:"uuid"`
	HistoryID     string  `json:"historyId"`
	Name          string  `json:"name"`
	FullName      string  `json:"fullName"`
	Status        string  `json:"status"`
	Stage         string  `json:"stage"`
	Start         int64   `json:"start"`
	Stop          int64   `json:"stop"`
	Steps         []step  `json:"steps"`
	Labels        []label `json:"labels"`
	Attachments   []any   `json:"attachments"`
	Parameters    []any   `json:"parameters"`
	Links         []any   `json:"links"`
	StatusDetails any     `json:"statusDetails,omitempty"`
}

// Run executes fn and writes an Allure result under ALLURE_RESULTS (or build/allure-results/go-hub).
func Run(t *testing.T, meta Meta, fn func(a *A)) {
	t.Helper()
	if meta.FullName == "" {
		meta.FullName = meta.Package + "." + meta.Name
	}
	if meta.Suite == "" {
		meta.Suite = meta.Feature
	}
	a := &A{t: t, meta: meta, start: time.Now()}
	defer a.finish()
	fn(a)
}

// Step records an Allure step (required for api/integration/e2e quality gate).
func (a *A) Step(name string, fn func()) {
	a.t.Helper()
	start := time.Now()
	status := "passed"
	var panicVal any
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicVal = r
				status = "broken"
			}
		}()
		fn()
	}()
	if a.t.Failed() && status == "passed" {
		status = "failed"
	}
	st := step{
		Name:        name,
		Status:      status,
		Stage:       "finished",
		Start:       start.UnixMilli(),
		Stop:        time.Now().UnixMilli(),
		Steps:       []any{},
		Attachments: []any{},
		Parameters:  []any{},
	}
	a.mu.Lock()
	a.steps = append(a.steps, st)
	a.mu.Unlock()
	if panicVal != nil {
		panic(panicVal)
	}
}

func (a *A) finish() {
	status := "passed"
	if a.t.Failed() {
		status = "failed"
	} else if a.t.Skipped() {
		status = "skipped"
	}
	dir := resultsDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		a.t.Logf("allurex: mkdir: %v", err)
		return
	}
	id := newUUID()
	fullName := a.meta.FullName
	payload := result{
		UUID:        id,
		HistoryID:   historyID(fullName, a.meta.Name),
		Name:        a.meta.Name,
		FullName:    fullName,
		Status:      status,
		Stage:       "finished",
		Start:       a.start.UnixMilli(),
		Stop:        time.Now().UnixMilli(),
		Steps:       a.steps,
		Labels:      a.labels(),
		Attachments: []any{},
		Parameters:  a.parameters(),
		Links:       []any{},
	}
	if status == "failed" {
		payload.StatusDetails = map[string]string{
			"message": fmt.Sprintf("%s failed", a.meta.Name),
		}
	}
	path := filepath.Join(dir, id+"-result.json")
	data, err := json.Marshal(payload)
	if err != nil {
		a.t.Logf("allurex: marshal: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		a.t.Logf("allurex: write: %v", err)
	}
}

func (a *A) labels() []label {
	labels := []label{
		{Name: "language", Value: "go"},
		{Name: "framework", Value: "go"},
		{Name: "host", Value: hostname()},
		{Name: "thread", Value: fmt.Sprintf("%d", os.Getpid())},
		{Name: "package", Value: a.meta.Package},
		{Name: "testClass", Value: a.meta.Package},
		{Name: "testMethod", Value: a.meta.Name},
		{Name: "suite", Value: a.meta.Suite},
	}
	if a.meta.Layer != "" {
		labels = append(labels, label{Name: "layer", Value: a.meta.Layer})
	}
	if a.meta.Component != "" {
		labels = append(labels, label{Name: "component", Value: a.meta.Component})
	}
	if a.meta.Epic != "" {
		labels = append(labels, label{Name: "epic", Value: a.meta.Epic})
	}
	if a.meta.Feature != "" {
		labels = append(labels, label{Name: "feature", Value: a.meta.Feature})
	}
	if a.meta.Story != "" {
		labels = append(labels, label{Name: "story", Value: a.meta.Story})
	}
	seen := map[string]struct{}{}
	addTag := func(tag string) {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return
		}
		if _, ok := seen[tag]; ok {
			return
		}
		seen[tag] = struct{}{}
		labels = append(labels, label{Name: "tag", Value: tag})
	}
	addTag(a.meta.Browser)
	for _, tag := range a.meta.Tags {
		addTag(tag)
	}
	return labels
}

func (a *A) parameters() []any {
	out := make([]any, 0, 1+len(a.meta.Parameters))
	if b := strings.TrimSpace(a.meta.Browser); b != "" {
		out = append(out, Parameter{Name: "browser", Value: b})
	}
	for _, p := range a.meta.Parameters {
		out = append(out, p)
	}
	return out
}

func resultsDir() string {
	d := os.Getenv("ALLURE_RESULTS")
	if d == "" {
		d = os.Getenv("ALLURE_RESULTS_DIR")
	}
	if d == "" {
		d = filepath.Join("build", "allure-results", "go-hub")
	}
	if filepath.IsAbs(d) {
		return d
	}
	if root := moduleRoot(); root != "" {
		return filepath.Join(root, d)
	}
	abs, err := filepath.Abs(d)
	if err != nil {
		return d
	}
	return abs
}

// moduleRoot finds the nearest go.mod (repo root for go test package cwd).
func moduleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func historyID(fullName, name string) string {
	sum := md5.Sum([]byte(fullName + "|" + name))
	return hex.EncodeToString(sum[:])
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return runtime.GOOS
	}
	return h
}

func newUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
