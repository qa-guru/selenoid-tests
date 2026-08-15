package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_usageWhenOutputOrEpicMissing(t *testing.T) {
	var stderr bytes.Buffer
	code := run(nil, strings.NewReader(""), ioDiscard(), &stderr)
	if code != 2 {
		t.Fatalf("code=%d", code)
	}
	if !strings.Contains(stderr.String(), "usage: gotest2allure") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRun_flagParseError(t *testing.T) {
	var stderr bytes.Buffer
	code := run([]string{"-bogus"}, strings.NewReader(""), ioDiscard(), &stderr)
	if code != 2 {
		t.Fatalf("code=%d", code)
	}
}

func TestRun_openInputMissing(t *testing.T) {
	var stderr bytes.Buffer
	code := run([]string{
		"--input", filepath.Join(t.TempDir(), "missing.jsonl"),
		"--output", t.TempDir(),
		"--epic", "selenoid",
	}, strings.NewReader(""), ioDiscard(), &stderr)
	if code != 1 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "open input") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestRun_convertErrorWhenOutputIsFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(out, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	code := run([]string{"--output", out, "--epic", "selenoid"}, strings.NewReader("{}\n"), ioDiscard(), &stderr)
	if code != 1 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
}

func TestRun_writesFromFileAndStdin(t *testing.T) {
	input := strings.Join([]string{
		`{"Time":"2026-08-02T12:00:00Z","Action":"run","Package":"pkg","Test":"TestFoo"}`,
		`{"Time":"2026-08-02T12:00:01Z","Action":"pass","Package":"pkg","Test":"TestFoo","Elapsed":1.0}`,
	}, "\n") + "\n"

	dir := t.TempDir()
	inFile := filepath.Join(dir, "events.jsonl")
	if err := os.WriteFile(inFile, []byte(input), 0o644); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "allure")
	var stdout bytes.Buffer
	code := run([]string{
		"--input", inFile,
		"--output", outDir,
		"--epic", "selenoid",
		"--component", "hub",
		"--layer", "unit",
	}, strings.NewReader(""), &stdout, ioDiscard())
	if code != 0 {
		t.Fatalf("file input code=%d", code)
	}
	if !strings.Contains(stdout.String(), "Wrote 1 Allure results") {
		t.Fatalf("stdout=%q", stdout.String())
	}

	stdout.Reset()
	code = run([]string{"--output", filepath.Join(dir, "allure-stdin"), "--epic", "selenoid"}, strings.NewReader(input), &stdout, ioDiscard())
	if code != 0 {
		t.Fatalf("stdin code=%d", code)
	}
}

func ioDiscard() *bytes.Buffer {
	return &bytes.Buffer{}
}
