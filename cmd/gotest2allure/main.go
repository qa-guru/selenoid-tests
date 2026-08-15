// Command gotest2allure converts `go test -json` output into Allure *-result.json.
//
// Usage:
//
//	go test -json ./... > events.jsonl
//	gotest2allure --input events.jsonl --output build/allure-results/go-selenoid --epic selenoid
//
//	go test -json ./... | gotest2allure --output DIR --epic NAME
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("gotest2allure", flag.ContinueOnError)
	fs.SetOutput(stderr)
	input := fs.String("input", "", "path to go test -json lines (default: stdin)")
	output := fs.String("output", "", "Allure results directory")
	epic := fs.String("epic", "", "Allure epic label")
	component := fs.String("component", "", "Allure component label (default: epic)")
	layer := fs.String("layer", "unit", "Allure layer label")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *output == "" || *epic == "" {
		fmt.Fprintln(stderr, "usage: gotest2allure --output DIR --epic NAME [--input FILE] [--component NAME] [--layer unit]")
		return 2
	}

	in := stdin
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			fmt.Fprintf(stderr, "gotest2allure: open input: %v\n", err)
			return 1
		}
		defer f.Close()
		in = f
	}

	n, err := allurex.ConvertGoTestJSON(in, allurex.ConvertOptions{
		Epic:      *epic,
		Component: *component,
		Layer:     *layer,
		OutputDir: *output,
	})
	if err != nil {
		fmt.Fprintf(stderr, "gotest2allure: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "Wrote %d Allure results to %s (epic=%s, layer=%s)\n", n, *output, *epic, *layer)
	return 0
}
