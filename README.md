# selenoid-tests

Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **playwright-image** — Go unit (в CI из исходных репо) + Java e2e/integration/api.

- Allure TestOps: проект `selenoid-tests`, **ALLURE_PROJECT_ID=5271**
- Test layers: `@Layer` keys → TestOps mapping (`e2e` → E2E Tests) — RAG `test-layers`, sync: `qa-guru-tms-automator/scripts/sync_testops_layer_mappings.py`
- Component filter: `@Component` label → TestOps custom field **Component** (`cm`, `selenoid`, `selenoid-ui`, `playwright-image`); sync: `qa-guru-tms-automator/scripts/sync_testops_component_mappings.py`
- Allure 3 GitHub Pages: `https://qa-guru.github.io/selenoid-tests/reports/<run-id>/`
- Dashboard: `.../reports/<run-id>/dashboard/index.html`

## Prerequisite (локально)

Hub/UI smoke (`testE2e`, `testApi`, …):

```bash
cd ../dev
./scripts/build-selenoid-ui.sh
./scripts/build-selenoid.sh
./scripts/start-selenoid.sh &
./scripts/start-selenoid-ui.sh &
```

CM integration (`testCmIntegration`): Docker + `cm` binary; hub/UI на **:4445/:8081** (dev-стек на :4444/:8080 не конфликтует).

```bash
cd ../cm && go build -o cm .
cd ../dev && ./scripts/build-selenoid-ui.sh   # ui/build для cross-compile selenoid-ui
./gradlew testCmIntegration -DskipHealthCheck=true   # lifecycle only (3 tests)
./gradlew testE2e -DskipHealthCheck=true             # + CmInstallerSessionTests (cm :4445)
```

## Gradle pyramid slices

```bash
./gradlew testUnit -DskipHealthCheck=true                              # unit
./gradlew testComponent -DskipHealthCheck=true                         # @Layer component
./gradlew testApi -DskipHealthCheck=true                             # @Layer api
./gradlew testIntegration -DskipHealthCheck=true                     # @Layer integration (без local-only)
./gradlew testE2e -DskipHealthCheck=true                             # e2e smoke
./gradlew testPlaywright -DskipHealthCheck=true                      # Playwright WS
./gradlew testResilience -DskipHealthCheck=true                      # hub kill/recovery (local-only)
./gradlew testCmIntegration -DskipHealthCheck=true                   # CM lifecycle (local-only, :4445)

# CI-эквиваленты
./gradlew test -Denv=selenoid_github_api -DincludeTags=api -DskipHealthCheck=true
./gradlew test -Denv=selenoid_github_integration -DincludeTags=integration -DskipHealthCheck=true
./gradlew test -DincludeTags=smoke,api -DexcludeTags=integration,resilience,local-only,playwright

./gradlew allureReport
```

Stand override: `-DpyramidStand=selenoid_github` → env `selenoid_github_api`, `selenoid_github_integration`, …

## CI

Workflow: `.github/workflows/selenoid_github-orchestrator.yml` (`name: selenoid-tests Tests`).

| Job | Что делает |
|-----|------------|
| `go-unit` (matrix) | Checkout `qa-guru/selenoid`, `selenoid-ui`, `cm` → Go unit → Allure |
| `java-e2e` | Push: `testUnit` + `testComponent` (gate) + `testApi` → TestOps; dispatch: `test_tags=integration|api|smoke` |
| `report` | Merge `build/allure-results/**` → `allureReport` → gh-pages → TestOps 5271 |

`workflow_dispatch`: `test_tags=integration` для integration slice; `env_profile=selenoid_github_api` для api-only.

### Dashboard (`allurerc.mjs`)

Пирамида: `unit → component → integration → api → e2e → manual`.

## Матрица

| Класс | Сервис | @Epic | @Layer | @Tag |
|-------|--------|-------|--------|------|
| ConfigReaderTest | — | — | unit | — |
| ConfigReaderCmTest | — | — | unit | — |
| ConfigReaderPlaywrightTest | — | — | unit | — |
| ConfigReaderUiTest | — | — | unit | — |
| ConfigReaderUrlTrimTest | — | — | unit | — |
| CreateSessionRequestJsonTest | — | — | unit | — |
| CmInstallerHelperTest | cm | CM | unit | — |
| CmRunResultTest | cm | CM | unit | — |
| ConfigOwnerMergeTest | — | — | unit | — |
| HubStatusJsonTest | selenoid | selenoid | component | — |
| HubStatusBrowsersJsonTest | selenoid | selenoid | component | — |
| UiStatusJsonTest | selenoid-ui | selenoid-ui | component | — |
| UiPingJsonTest | selenoid-ui | selenoid-ui | component | — |
| UiErrorJsonTest | selenoid-ui | selenoid-ui | component | — |
| SseStateJsonTest | selenoid-ui | selenoid-ui | component | — |
| SseErrorsJsonTest | selenoid-ui | selenoid-ui | component | — |
| SessionCreateJsonTest | selenoid | selenoid | component | — |
| BrowsersConfigJsonTest | selenoid-ui | selenoid-ui | component | — |
| HubStatusForwardCompatJsonTest | selenoid | selenoid | component | — |
| HubPingTests | selenoid | selenoid | api | api |
| HubStatusBrowsersTests | selenoid | selenoid | api | api |
| HubSessionDeleteUnknownTests | selenoid | selenoid | api | api, negative |
| HubSessionInvalidBrowserTests | selenoid | selenoid | api | api, negative |
| UiBrowsersConfigTests | selenoid-ui | selenoid-ui | api | api |
| UiStatusTotalTests | selenoid-ui | selenoid-ui | api | api |
| UiPingVersionTests | selenoid-ui | selenoid-ui | api | api |
| UiSseMultipleEventsTests | selenoid-ui | selenoid-ui | api | api |
| PlaywrightUnknownPathTests | playwright-image | playwright-image | api | api, negative |
| HubStatusTests | selenoid | selenoid | api | api |
| HubSessionApiTests | selenoid | selenoid | api | api |
| PlaywrightEndpointTests | playwright-image | playwright-image | api | api |
| UiStatusTests | selenoid-ui | selenoid-ui | api | api |
| UiSseStreamTests | selenoid-ui | selenoid-ui | api | api |
| UiPingTests | selenoid-ui | selenoid-ui | api | api |
| HubPlaywrightSessionTests | playwright-image | playwright-image | integration | integration |
| HubPlaywrightNavigateTests | playwright-image | playwright-image | integration | integration |
| UiStatusWithSessionTests | selenoid-ui | selenoid-ui | integration | integration |
| UiSseWithSessionTests | selenoid-ui | selenoid-ui | integration | integration |
| StackHealthTests | selenoid-ui | selenoid-ui | integration | integration |
| CmInstallerLifecycleTests | cm | CM | integration | integration, local-only, cm |
| CmInstallerSessionTests | cm | CM | e2e | smoke, cm |
| UiHubStatusConsistencyTests | selenoid-ui | selenoid-ui | integration | integration |
| UiStatusWhenHubDownTests | selenoid-ui | selenoid-ui | integration | integration, local-only |
| UiStatusRecoveryTests | selenoid-ui | selenoid-ui | integration | integration, resilience, local-only |
| HubSessionTests | selenoid | selenoid | e2e | smoke |
| HubSessionIdTests | selenoid | selenoid | e2e | smoke |
| HubSessionHeadingTests | selenoid | selenoid | e2e | smoke |
| HubSessionTitleTests | selenoid | selenoid | e2e | smoke |
| HubPlaywrightSessionTests | playwright-image | playwright-image | e2e | playwright, smoke |
| UiStatusBarTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiDashboardLoadTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiSseIndicatorTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiReloadTests | selenoid-ui | selenoid-ui | e2e | smoke |

## Config keys

| Key | Default |
|-----|---------|
| apiBaseUrl | `""` (→ hubUrl) |
| hubUrl | http://127.0.0.1:4444/ |
| uiUrl | http://127.0.0.1:8080/ |
| remoteUrl | http://127.0.0.1:4444/wd/hub |
| cmHubPort | 4445 (CM installer; dev hub stays :4444) |
| cmUiPort | 8081 |
| playwrightWsEndpoint | ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1 |

Override: `-DhubUrl`, `-DuiUrl`, `-DapiBaseUrl`, …
