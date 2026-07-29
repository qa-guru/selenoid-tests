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
    har-prod)       export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_qa_guru}_e2e" ;;
    min)            export SELENOID_TEST_ENV=selenoid_github_min_integration ;;
    *)              export SELENOID_TEST_ENV="${PYRAMID_STAND:-selenoid_github}_api" ;;
  esac
fi

export SELENOID_TEST_SKIP_HEALTH_CHECK="${SELENOID_TEST_SKIP_HEALTH_CHECK:-true}"

run_pkgs() {
  local pkgs=("$@")
  local test_flags=(-count=1 -timeout=15m)
  if [[ "${SLICE}" == "integration" || "${SLICE}" == "min" || "${SLICE}" == "ui" || "${SLICE}" == "e2e" || "${SLICE}" == "webdriver" || "${SLICE}" == "playwright" || "${SLICE}" == "har-prod" ]]; then
    # Hub session tests share capacity counters (Java @ResourceLock hubSessions).
    test_flags+=(-p 1)
  fi
  echo "go pyramid slice=${SLICE} env=${SELENOID_TEST_ENV:-} allure=${ALLURE_RESULTS}"
  go test "${pkgs[@]}" "${test_flags[@]}"
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
  har-prod)
    # P3 prod HAR smoke: hub enableHAR + HarCapture bodies on qa_guru_e2e (−cm/min/resilience).
    unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
    run_pkgs ./tests/e2e/har/...
    ;;
  hub-all)
    # P1: unit + component + api (−cm).
    run_pkgs ./internal/config/... ./internal/helpers/... ./internal/hubapi/... ./tests/component/... ./tests/api/...
    ;;
  integration)
    # P2 WD warm sessions + UI status/SSE + PW sessions (−min/−cm).
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
    run_pkgs ./tests/integration/wd/... ./tests/integration/ui/... ./tests/integration/pw/...
    ;;
  ui)
    # P3 UI e2e smoke via local playwright-go Chromium (not remote WS connect).
    unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
    run_pkgs ./tests/e2e/ui/...
    ;;
  webdriver)
    # P3 WebDriver e2e smoke via raw WD HTTP (HubSession* parity).
    run_pkgs ./tests/e2e/webdriver/...
    ;;
  playwright)
    # P3 Playwright e2e smoke via remote WS connect (HubPlaywrightSessionTests parity, −min).
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
    run_pkgs ./tests/e2e/playwright/...
    ;;
  e2e)
    # Umbrella: Playwright WS + UI local Chromium + WebDriver e2e (+ har when present).
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
    run_pkgs ./tests/e2e/playwright/...
    unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
    run_pkgs ./tests/e2e/ui/... ./tests/e2e/webdriver/... ./tests/e2e/har/...
    ;;
  min)
    # P4 local-only: offline min catalogs + live min WD/PW sessions (github min stack).
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
    run_pkgs ./tests/component/min/... ./tests/integration/min/...
    ;;
  *)
    echo "Unknown slice: ${SLICE}" >&2
    echo "Known: unit|component|api|integration|ui|webdriver|playwright|e2e|har-prod|hub-prod|hub-all|min" >&2
    exit 2
    ;;
esac
