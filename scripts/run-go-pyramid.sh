#!/usr/bin/env bash
# Run Go pyramid slices (ADR-go-pyramid). Mirrors Gradle testApi / testHubProd entrypoints.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

SLICE="${1:-api}"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export ALLURE_RESULTS="${ALLURE_RESULTS:-${ROOT}/build/allure-results/go-hub}"
mkdir -p "${ALLURE_RESULTS}"

# Profile: SELENOID_TEST_ENV or env=… (Owner-compatible). Default per slice.
if [[ -z "${SELENOID_TEST_ENV:-}" && -z "${env:-}" ]]; then
  case "${SLICE}" in
    unit|component) export SELENOID_TEST_ENV=local_unit ;;
    api|hub-all)    export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_github}_api" ;;
    integration)    export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_github}_integration" ;;
    e2e|webdriver|ui|playwright|hub-prod)
                    export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_github}_e2e" ;;
    min)            export SELENOID_TEST_ENV=selenoid_github_min_integration ;;
    *)              export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_github}_api" ;;
  esac
fi

export SELENOID_TEST_SKIP_HEALTH_CHECK="${SELENOID_TEST_SKIP_HEALTH_CHECK:-true}"

run_pkgs() {
  local pkgs=("$@")
  echo "go pyramid slice=${SLICE} env=${SELENOID_TEST_ENV:-} allure=${ALLURE_RESULTS}"
  go test "${pkgs[@]}" -count=1 -timeout=15m
}

case "${SLICE}" in
  unit)
    # Offline unit (−cm): config resolve/merge, WD session body, HarCapture pure.
    run_pkgs ./internal/config/... ./internal/helpers/... ./internal/hubapi/...
    ;;
  component)
    # Offline JSON/parser fixtures (−cm/−min). Profile: local_unit.
    run_pkgs ./tests/component/...
    ;;

  api)
    run_pkgs ./tests/api/...
    ;;
  hub-prod)
    # Prod-safe: api only (−cm/min/resilience).
    run_pkgs ./tests/api/...
    ;;
  hub-all)
    # P1: unit + component + api (−cm).
    run_pkgs ./internal/config/... ./internal/helpers/... ./internal/hubapi/... ./tests/component/... ./tests/api/...
    ;;
  *)
    echo "Unknown slice: ${SLICE}" >&2
    echo "Known: unit|component|api|hub-prod|hub-all" >&2
    exit 2
    ;;
esac
