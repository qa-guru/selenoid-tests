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
	"os"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
)

func main() {
	input := flag.String("input", "", "path to go test -json lines (default: stdin)")
	output := flag.String("output", "", "Allure results directory")
	epic := flag.String("epic", "", "Allure epic label")
	component := flag.String("component", "", "Allure component label (default: epic)")
	layer := flag.String("layer", "unit", "Allure layer label")
	flag.Parse()

	if *output == "" || *epic == "" {
		fmt.Fprintln(os.Stderr, "usage: gotest2allure --output DIR --epic NAME [--input FILE] [--component NAME] [--layer unit]")
		os.Exit(2)
	}

	in := os.Stdin
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gotest2allure: open input: %v\n", err)
			os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "gotest2allure: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %d Allure results to %s (epic=%s, layer=%s)\n", n, *output, *epic, *layer)
}
