#!/usr/bin/env bash
# Offline coverage for Sonar — own orchestrator code only (no repos/**, no prod sessions).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GO111MODULE=on
export SELENOID_TEST_ENV=local_unit

go test -coverprofile=coverage.txt -covermode=atomic \
  ./internal/... \
  ./tests/component/...

test -s coverage.txt
