package cm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qa-guru/selenoid-tests/internal/allurex"
	"github.com/qa-guru/selenoid-tests/internal/cm"
)

func TestCmRunResult_RequireSuccessPassesOnZeroExit(t *testing.T) {
	allurex.Run(t, unitMeta("requireSuccess passes on exit code zero", "helpers.CmRunResultTest", "CmRunResult"), func(a *allurex.A) {
		a.Step("requireSuccess", func() {
			result := cm.RunResult{ExitCode: 0, Output: "ok"}
			require.NotPanics(t, func() { result.RequireSuccess("configure") })
		})
	})
}

func TestCmRunResult_RequireSuccessFailsOnNonZeroExit(t *testing.T) {
	allurex.Run(t, unitMeta("requireSuccess fails on non-zero exit code", "helpers.CmRunResultTest", "CmRunResult"), func(a *allurex.A) {
		a.Step("requireSuccess", func() {
			result := cm.RunResult{ExitCode: 1, Output: "boom"}
			require.Panics(t, func() { result.RequireSuccess("start") })
		})
	})
}
