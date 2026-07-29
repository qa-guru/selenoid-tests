#!/usr/bin/env bash
# Run Go pyramid slices (ADR-go-pyramid). Mirrors Gradle testApi / testHubProd / testHubAll.
#
# Composite gates (−cm; prod also −min/−resilience):
#   hub-prod  = unit → component → integration → api → ui → webdriver → playwright → har-prod
#               (testHubProd semantics; har-prod = HubHarProd/HarCapture on selenoid_qa_guru_e2e)
#   hub-all   = hub-prod layers on selenoid_github + min + resilience (−cm; cm = slice only)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

SLICE="${1:-api}"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export ALLURE_RESULTS="${ALLURE_RESULTS:-${ROOT}/build/allure-results/go-hub}"
mkdir -p "${ALLURE_RESULTS}"

STAND="${PYRAMID_STAND:-selenoid_github}"

# Profile: SELENOID_TEST_ENV or env=… (Owner-compatible). Default per single slice only.
if [[ -z "${SELENOID_TEST_ENV:-}" && -z "${env:-}" ]]; then
  case "${SLICE}" in
    unit|component) export SELENOID_TEST_ENV=local_unit ;;
    api)            export SELENOID_TEST_ENV="${STAND}_api" ;;
    integration)    export SELENOID_TEST_ENV="${STAND}_integration" ;;
    e2e|webdriver|ui|playwright)
                    export SELENOID_TEST_ENV="${STAND}_e2e" ;;
    har-prod)       export SELENOID_TEST_ENV=selenoid_qa_guru_e2e ;;
    min)            export SELENOID_TEST_ENV=selenoid_github_min_integration ;;
    cm)             export SELENOID_TEST_ENV=selenoid_github_cm_integration ;;
    resilience)     export SELENOID_TEST_ENV="${STAND}_e2e" ;;
    hub-prod)       STAND="${PYRAMID_STAND:-selenoid_qa_guru}" ;;
    hub-all)        STAND="${PYRAMID_STAND:-selenoid_github}" ;;
    *)              export SELENOID_TEST_ENV="${STAND}_api" ;;
  esac
fi

export SELENOID_TEST_SKIP_HEALTH_CHECK="${SELENOID_TEST_SKIP_HEALTH_CHECK:-true}"

run_pkgs() {
  local active_slice="${1:?slice name required}"
  shift
  local pkgs=("$@")
  local test_flags=(-count=1 -timeout=15m)
  if [[ "${active_slice}" == "integration" || "${active_slice}" == "min" || "${active_slice}" == "resilience" || "${active_slice}" == "ui" || "${active_slice}" == "e2e" || "${active_slice}" == "webdriver" || "${active_slice}" == "playwright" || "${active_slice}" == "har-prod" || "${active_slice}" == "cm" ]]; then
    # Hub session tests share capacity counters (Java @ResourceLock hubSessions).
    test_flags+=(-p 1)
  fi
  echo "go pyramid slice=${active_slice} env=${SELENOID_TEST_ENV:-} allure=${ALLURE_RESULTS}"
  go test "${pkgs[@]}" "${test_flags[@]}"
}

run_with_env() {
  local active_slice="$1"
  local env_name="$2"
  shift 2
  local saved_env="${SELENOID_TEST_ENV:-}"
  export SELENOID_TEST_ENV="${env_name}"
  "$@"
  if [[ -n "${saved_env}" ]]; then
    export SELENOID_TEST_ENV="${saved_env}"
  else
    unset SELENOID_TEST_ENV
  fi
}

run_unit_pkgs() {
  run_pkgs unit ./internal/config/... ./internal/helpers/... ./internal/hubapi/...
}

run_component_pkgs() {
  run_pkgs component ./tests/component/...
}

run_api_pkgs() {
  run_pkgs api ./tests/api/...
}

run_integration_pkgs() {
  export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
  run_pkgs integration ./tests/integration/wd/... ./tests/integration/ui/... ./tests/integration/pw/...
}

run_ui_pkgs() {
  unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
  run_pkgs ui ./tests/e2e/ui/...
}

run_webdriver_pkgs() {
  run_pkgs webdriver ./tests/e2e/webdriver/...
}

run_playwright_pkgs() {
  export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
  run_pkgs playwright ./tests/e2e/playwright/...
}

run_har_prod_pkgs() {
  unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
  run_pkgs har-prod ./tests/e2e/har/...
}

run_min_pkgs() {
  export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
  run_pkgs min ./tests/component/min/... ./tests/integration/min/...
}

run_resilience_pkgs() {
  unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
  run_pkgs resilience ./tests/integration/resilience/...
}

# hub-prod: testHubProd (−cm/−min/−resilience) + har-prod on qa_guru_e2e.
run_hub_prod() {
  local stand="${PYRAMID_STAND:-selenoid_qa_guru}"
  echo "go pyramid composite=hub-prod stand=${stand} (−cm/−min/−resilience + har-prod)"
  run_with_env unit local_unit run_unit_pkgs
  run_with_env component local_unit run_component_pkgs
  run_with_env integration "${stand}_integration" run_integration_pkgs
  run_with_env api "${stand}_api" run_api_pkgs
  run_with_env ui "${stand}_e2e" run_ui_pkgs
  run_with_env webdriver "${stand}_e2e" run_webdriver_pkgs
  run_with_env playwright "${stand}_e2e" run_playwright_pkgs
  run_with_env har-prod selenoid_qa_guru_e2e run_har_prod_pkgs
}

# hub-all: testHubAll on github stand (−cm; cm = slice `cm` only).
run_hub_all() {
  local stand="${PYRAMID_STAND:-selenoid_github}"
  echo "go pyramid composite=hub-all stand=${stand} (−cm; +min/−resilience on github)"
  run_with_env unit local_unit run_unit_pkgs
  run_with_env component local_unit run_component_pkgs
  run_with_env integration "${stand}_integration" run_integration_pkgs
  run_with_env api "${stand}_api" run_api_pkgs
  run_with_env ui "${stand}_e2e" run_ui_pkgs
  run_with_env webdriver "${stand}_e2e" run_webdriver_pkgs
  run_with_env playwright "${stand}_e2e" run_playwright_pkgs
  run_with_env min selenoid_github_min_integration run_min_pkgs
  run_with_env resilience "${stand}_e2e" run_resilience_pkgs
}

case "${SLICE}" in
  unit)
    run_unit_pkgs
    ;;
  component)
    run_component_pkgs
    ;;

  api)
    run_api_pkgs
    ;;
  hub-prod)
    run_hub_prod
    ;;
  har-prod)
    run_har_prod_pkgs
    ;;
  hub-all)
    run_hub_all
    ;;
  integration)
    run_integration_pkgs
    ;;
  ui)
    run_ui_pkgs
    ;;
  webdriver)
    run_webdriver_pkgs
    ;;
  playwright)
    run_playwright_pkgs
    ;;
  e2e)
    # Umbrella: Playwright WS + UI local Chromium + WebDriver e2e (+ har when present).
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD="${PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD:-1}"
    run_playwright_pkgs
    unset PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD
    run_pkgs e2e ./tests/e2e/ui/... ./tests/e2e/webdriver/... ./tests/e2e/har/...
    ;;
  min)
    run_min_pkgs
    ;;
  resilience)
    run_resilience_pkgs
    ;;
  cm)
    # P4 local-only: CM pyramid on :4445/:8081 (github CM stack; −prod). Prerequisite: start-ci-cm-stack.sh for api.
    echo "go pyramid slice=${SLICE} env=${SELENOID_TEST_ENV:-} allure=${ALLURE_RESULTS}"
    go test ./internal/config/ -run 'TestConfigReader_ResolveCm' -count=1 -timeout=15m
    go test ./internal/cm/... -count=1 -timeout=15m
    go test ./tests/cm/component/... -count=1 -timeout=15m
    go test ./tests/cm/api/... -count=1 -timeout=15m -p 1
    go test ./tests/cm/integration/... -count=1 -timeout=25m -p 1
    go test ./tests/cm/e2e/... -count=1 -timeout=25m -p 1
    ;;
  *)
    echo "Unknown slice: ${SLICE}" >&2
    echo "Known: unit|component|api|integration|ui|webdriver|playwright|e2e|har-prod|hub-prod|hub-all|min|resilience|cm" >&2
    exit 2
    ;;
esac
