// Package cm wraps the qa-guru/cm CLI for installer lifecycle tests.
package cm

import (
	"fmt"
)

// RunResult is cm command stdout + exit code (Java CmInstallerHelper.CmRunResult).
type RunResult struct {
	ExitCode int
	Output   string
}

// RequireSuccess fails the test when exit code is non-zero.
func (r RunResult) RequireSuccess(action string) {
	if r.ExitCode != 0 {
		panic(fmt.Sprintf("%s failed (exit %d):\n%s", action, r.ExitCode, r.Output))
	}
}
