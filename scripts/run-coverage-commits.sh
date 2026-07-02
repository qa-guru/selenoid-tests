#!/usr/bin/env bash
# Sequential commits for pyramid coverage expansion (50 commits, ~2.5 min apart).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

INTERVAL="${COMMIT_INTERVAL_SEC:-150}"
LOG="$ROOT/build/coverage-commit-series.log"
mkdir -p "$(dirname "$LOG")"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG"
}

commit_push() {
  local msg="$1"
  shift
  if git diff --cached --quiet && [[ "$#" -eq 0 ]]; then
    log "SKIP (nothing staged): $msg"
    return 0
  fi
  if [[ "$#" -gt 0 ]]; then
    git add "$@"
  fi
  if git diff --cached --quiet; then
    log "SKIP (empty stage): $msg"
    return 0
  fi
  git commit -m "$msg"
  git push origin HEAD
  log "DONE: $msg"
  log "Sleep ${INTERVAL}s..."
  sleep "$INTERVAL"
}

log "Starting coverage commit series from $(pwd)"

# --- unit (10) ---
commit_push "test(unit): ConfigReader CM hub/ui/remote URLs" \
  src/test/java/config/ConfigReader.java \
  src/test/java/config/TestConfig.java \
  src/test/resources/config/default.properties \
  src/test/java/config/ConfigReaderCmTest.java

commit_push "test(unit): ConfigReader playwright query builder" \
  src/test/java/config/ConfigReaderPlaywrightTest.java

commit_push "test(unit): ConfigReader UI URL fail-fast and slash" \
  src/test/java/config/ConfigReaderUiTest.java

commit_push "test(unit): ConfigReader URL trim normalization" \
  src/test/java/config/ConfigReaderUrlTrimTest.java

commit_push "test(unit): CreateSessionRequest JSON serialization" \
  src/test/java/config/CreateSessionRequestJsonTest.java

commit_push "test(unit): CmInstallerHelper implementation" \
  src/test/java/helpers/CmInstallerHelper.java

commit_push "test(unit): CmInstallerHelper temp config directory" \
  src/test/java/helpers/CmInstallerHelperTest.java

commit_push "test(unit): CmRunResult requireSuccess contract" \
  src/test/java/helpers/CmRunResultTest.java

commit_push "test(unit): Owner system property overrides" \
  src/test/java/config/ConfigOwnerMergeTest.java

commit_push "test(unit): ConfigReader baseline hub and api URLs" \
  src/test/java/config/ConfigReaderTest.java

commit_push "test(unit): Gradle testUnit slice for helpers" \
  build.gradle

# --- component (10) ---
commit_push "test(component): fixture loader and testComponent slice" \
  build.gradle \
  src/test/java/tests/component/FixtureJson.java

commit_push "test(component): parse idle hub status fixture" \
  src/test/resources/fixtures/hub/status-idle.json \
  src/test/java/tests/component/HubStatusJsonTest.java

commit_push "test(component): parse hub status browsers map fixture" \
  src/test/resources/fixtures/hub/status-with-browsers.json \
  src/test/java/tests/component/HubStatusBrowsersJsonTest.java

commit_push "test(component): parse UI status wrapper fixture" \
  src/test/resources/fixtures/ui/status.json \
  src/test/java/tests/component/UiStatusJsonTest.java

commit_push "test(component): parse UI ping fixture" \
  src/test/resources/fixtures/ui/ping.json \
  src/test/java/tests/component/UiPingJsonTest.java

commit_push "test(component): parse UI error fixture" \
  src/test/resources/fixtures/ui/error.json \
  src/test/java/tests/component/UiErrorJsonTest.java

commit_push "test(component): parse SSE state fixture" \
  src/test/resources/fixtures/sse/state.json \
  src/test/java/tests/component/SseStateJsonTest.java

commit_push "test(component): parse SSE errors fixture" \
  src/test/resources/fixtures/sse/errors.json \
  src/test/java/tests/component/SseErrorsJsonTest.java

commit_push "test(component): parse session create response fixture" \
  src/test/resources/fixtures/hub/session-create.json \
  src/test/java/tests/component/SessionCreateJsonTest.java

commit_push "test(component): browsers config and forward-compatible status fixtures" \
  src/test/resources/fixtures/ui/browsers-config.json \
  src/test/resources/fixtures/hub/status-unknown-field.json \
  src/test/java/tests/component/BrowsersConfigJsonTest.java \
  src/test/java/tests/component/HubStatusForwardCompatJsonTest.java

# --- api (10) ---
commit_push "test(api): selenoid hub GET /ping" \
  src/test/java/api/hub/HubPingResponse.java \
  src/test/java/api/hub/HubPingApi.java \
  src/test/java/tests/api/HubPingTests.java

commit_push "test(api): selenoid hub status browser families" \
  src/test/java/tests/api/HubStatusBrowsersTests.java

commit_push "test(api): selenoid DELETE unknown session" \
  src/test/java/api/hub/HubSessionApi.java \
  src/test/java/tests/api/HubSessionDeleteUnknownTests.java

commit_push "test(api): selenoid POST session invalid browser" \
  src/test/java/tests/api/HubSessionInvalidBrowserTests.java

commit_push "test(api): selenoid-ui GET /browsers-config" \
  src/test/java/api/ui/UiBrowsersConfigApi.java \
  src/test/java/tests/api/UiBrowsersConfigTests.java

commit_push "test(api): selenoid-ui status total counter" \
  src/test/java/tests/api/UiStatusTotalTests.java

commit_push "test(api): selenoid-ui ping version metadata" \
  src/test/java/tests/api/UiPingVersionTests.java

commit_push "test(api): selenoid-ui SSE two consecutive events" \
  src/test/java/api/ui/SseStreamApi.java \
  src/test/java/tests/api/UiSseMultipleEventsTests.java

commit_push "test(api): playwright unknown path returns 400" \
  src/test/java/api/hub/PlaywrightEndpointApi.java \
  src/test/java/tests/api/PlaywrightUnknownPathTests.java

commit_push "test(api): playwright endpoint upgrade required smoke" \
  src/test/java/tests/api/PlaywrightEndpointTests.java

# --- integration (10) ---
commit_push "test(integration): ui status mirrors hub counters" \
  src/test/java/tests/integration/UiHubStatusConsistencyTests.java

commit_push "test(integration): ui status used increments with hub session" \
  src/test/java/tests/integration/UiStatusWithSessionTests.java

commit_push "test(integration): sse payload with active hub session" \
  src/test/java/tests/integration/UiSseWithSessionTests.java

commit_push "test(integration): ui status when hub is down" \
  src/test/java/tests/integration/UiStatusWhenHubDownTests.java \
  src/test/resources/config/local_integration.properties \
  src/test/resources/config/selenoid_github_integration.properties

commit_push "test(integration): ui status recovery after hub restart" \
  src/test/java/tests/UiStatusRecoveryTests.java

commit_push "test(integration): hub playwright session api helper" \
  src/test/java/api/hub/PlaywrightSessionApi.java

commit_push "test(integration): playwright ws session lifecycle" \
  src/test/java/tests/integration/HubPlaywrightSessionTests.java

commit_push "test(integration): playwright navigate example.com" \
  src/test/java/tests/integration/HubPlaywrightNavigateTests.java

commit_push "test(integration): cm installer lifecycle via cm helper" \
  src/test/java/tests/integration/CmInstallerLifecycleTests.java \
  src/test/resources/config/local_cm_integration.properties

commit_push "test(integration): stack health hub and ui ping" \
  src/test/java/tests/integration/StackHealthTests.java

# --- e2e (10) ---
commit_push "test(e2e): hub remote chrome opens example.com" \
  src/test/java/tests/HubSessionTests.java

commit_push "test(e2e): ui status bar stays connected" \
  src/test/java/tests/UiStatusBarTests.java

commit_push "test(e2e): hub session assigns session id" \
  src/test/java/tests/HubSessionIdTests.java

commit_push "test(e2e): ui dashboard opens root URL" \
  src/test/java/tests/UiDashboardLoadTests.java

commit_push "test(e2e): ui sse indicator ok styling" \
  src/test/java/tests/UiSseIndicatorTests.java

commit_push "test(e2e): hub session renders Example Domain heading" \
  src/test/java/tests/HubSessionHeadingTests.java

commit_push "test(e2e): playwright ws opens example.com" \
  src/test/java/tests/HubPlaywrightSessionTests.java

commit_push "test(e2e): hub session page title smoke" \
  src/test/java/tests/HubSessionTitleTests.java

commit_push "test(e2e): ui dashboard reload keeps connected" \
  src/test/java/tests/UiReloadTests.java

commit_push "test(e2e): cm installed stack remote session" \
  src/test/java/tests/CmInstallerSessionTests.java

commit_push "docs: update README test matrix for pyramid expansion" \
  README.md

log "Coverage commit series finished."
