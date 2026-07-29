package cm_test

import "github.com/qa-guru/selenoid-tests/internal/allurex"

func unitMeta(name, pkg, suite string) allurex.Meta {
	return allurex.Meta{
		Name:      name,
		Package:   pkg,
		Layer:     "unit",
		Component: "cm",
		Epic:      "cm",
		Suite:     suite,
		Tags:      []string{"unit", "cm"},
	}
}
