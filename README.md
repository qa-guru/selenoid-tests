# selenoid-tests

Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **browser-image** (`playwright/` + `webdriver/`) — Go unit (в CI из исходных репо) + Java e2e/integration/api.

**Scope:** `warm-pool-orchestrator/` — out of scope (deferred), не в матрице и не в CI.

- Allure TestOps: проект `selenoid-tests`, **ALLURE_PROJECT_ID=5271**
- Test layers: `@Layer` keys → TestOps mapping (`e2e` → E2E Tests) — RAG `test-layers`, sync: `qa-guru-tms-automator/scripts/sync_testops_layer_mappings.py`
- Component filter: `@Component` label → TestOps custom field **Component** (`cm`, `selenoid`, `selenoid-ui`, `playwright-image`, `webdriver-image`); sync: `qa-guru-tms-automator/scripts/sync_testops_component_mappings.py`
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
./gradlew testHubAll -DskipHealthCheck=true                             # full hub/UI pyramid (CI push gate)
./gradlew testUnit -DskipHealthCheck=true                              # unit only
./gradlew testComponent -DskipHealthCheck=true                         # @Layer component
./gradlew testApi -DskipHealthCheck=true                             # @Layer api
./gradlew testIntegration -DskipHealthCheck=true                     # @Layer integration (incl. local-only)
./gradlew testE2e -DskipHealthCheck=true                             # e2e smoke (hub + UI)
./gradlew testWebdriverE2e -DskipHealthCheck=true                    # webdriver-image smoke (HubSession*)
./gradlew testUiE2e -DskipHealthCheck=true                           # selenoid-ui smoke (Ui*)
./gradlew testPlaywright -DskipHealthCheck=true                      # Playwright WS (incl. firefox/webkit local-only)
./gradlew testResilience -DskipHealthCheck=true                      # hub restart recovery
./gradlew testMin -DskipHealthCheck=true                             # chromium-min + chrome-min
./gradlew testCmIntegration -DskipHealthCheck=true                   # CM lifecycle (CI: java-cm job; local: ports :4445/:8081)

# CI-эквиваленты (slice-only dispatch)
./gradlew test -Denv=selenoid_github_api -DincludeTags=api -DskipHealthCheck=true
./gradlew test -Denv=selenoid_github_integration -DincludeTags=integration -DskipHealthCheck=true
./gradlew test -Denv=selenoid_github_min_integration -DincludeTags=min -DskipHealthCheck=true

./gradlew allureReport
```

Stand override: `-DpyramidStand=selenoid_github` → env `selenoid_github_api`, `selenoid_github_integration`, …

### Prod hub (`selenoid.autotests.cloud`)

Profiles: `selenoid_autotests_cloud_api`, `selenoid_autotests_cloud_e2e` — remote hub `https://selenoid.autotests.cloud` (auth `user1:1234` in properties, same as deploy smoke).

```bash
./gradlew testApi -DpyramidStand=selenoid_autotests_cloud -DskipHealthCheck=true
./gradlew testE2e -DpyramidStand=selenoid_autotests_cloud -DskipHealthCheck=true
```

Post-deploy: `selenoid.autotests.cloud` → Actions → `trigger-deploy-smoke` → `repository_dispatch deploy-smoke` → this repo (`skip_go_unit`, `env_profile=selenoid_autotests_cloud_api`).

Prod caveats (nginx): `HubStatusApi` uses raw `GET /hub/status` (not UI `/status` with `.state`). `GET /logs/{id}` — nginx → hub (auth); UI uses `/ws/logs/{id}`.

### `testPlaywright` prerequisite

Hub на `:4444` и Docker-образ из `fixtures/ci-browsers.json` / `dev/browsers.json`:

```bash
cd dev && ./scripts/start-selenoid.sh &
docker pull qaguru/playwright-chromium:1.61.1   # или ./scripts/pull-browser-images.sh
./gradlew testPlaywright -DskipHealthCheck=true
```

В CI `scripts/start-ci-selenoid-stack.sh` тянет образы из `fixtures/ci-browsers.json` (chrome + playwright-chromium).

### Playwright-chromium-min (`1.61.1-min`)

Образ в `fixtures/ci-browsers.json`; endpoint — `selenoid_github_min_integration.properties` (VNC/video off). Входит в `testHubAll` (`testMin`).

```bash
./gradlew testMin -DskipHealthCheck=true
```

| Класс | @Layer | @Tag |
|-------|--------|------|
| PlaywrightMinCatalogJsonTest | component | min |
| HubPlaywrightMinSessionTests (`tests.integration`) | integration | min |
| HubPlaywrightMinSessionTests (`tests`) | e2e | min |

## CI

Workflow: `.github/workflows/selenoid_github-orchestrator.yml` (`name: selenoid-tests Tests`).

| Job | Что делает |
|-----|------------|
| `go-unit` (matrix) | Checkout `qa-guru/selenoid`, `selenoid-ui`, `cm` → Go unit → Allure |
| `java-e2e` | Push/default: `testHubAll` (all hub slices incl. playwright, resilience, local-only, min); slice: `test_tags=…` |
| `java-cm` | Push: `testCmIntegration` + `testCmApi` + `testCmE2e` (CM :4445/:8081); dispatch `test_tags=cm` |
| `report` | Merge `build/allure-results/**` → `allureReport` → gh-pages → TestOps 5271 |

`workflow_dispatch`: `test_tags=integration` для integration slice; `env_profile=selenoid_github_api` для api-only.

### Component × Layer × CI (push `main`)

Пирамида: `unit → component → integration → api → e2e → manual`. Числа — Java-классы в матрице ниже; **Go unit** — отдельно в `go-unit`.

| Component | unit | component | integration | api | e2e | manual | CI push |
|-----------|:----:|:---------:|:-----------:|:---:|:---:|:------:|---------|
| **selenoid** | Go | 8 | 1 | 10 | △¹ | — | `go-unit` + `testHubAll` |
| **selenoid-ui** | Go + 1 | 6 | 7 | 7 | 5 | △ | `go-unit` + `testHubAll` |
| **cm** | Go + 3 | 4 | 2 | 3 | 1 | — | `go-unit` + `java-cm` |
| **playwright-image** | 1 | 3 | 5 | 2 | 2 | — | `testHubAll` |
| **webdriver-image** | — | 2 | 2 | 2 | 4 | — | `testHubAll` |
| **dev** | — | △² | △³ | — | — | ✓ | — |
| **selenoid-autotests-cloud** | — | — | — | △⁴ | △⁵ | ✓ | deploy-smoke dispatch |

¹ **selenoid e2e:** нет `@Component("selenoid")` e2e-класса; hub-path — `webdriver-image` / `playwright-image` e2e в `testHubAll`.  
² `browsers.json` — косвенно через component JSON fixtures.  
³ `start-ci-selenoid-stack.sh` — оркестрация CI, без отдельного test-class.  
⁴ Post-deploy: `selenoid_autotests_cloud_api` (default `trigger-deploy-smoke`).  
⁵ Профиль `selenoid_autotests_cloud_e2e` — вручную / расширенный deploy-smoke.

### Manual (runbook)

| Сценарий | Где | Как |
|----------|-----|-----|
| Локальный стек hub/UI | `../dev/README.md` | `build-selenoid*.sh` + `start-selenoid*.sh` |
| Prod hub smoke | `selenoid-autotests-cloud` | `./deploy/smoke-remote.sh https://selenoid.autotests.cloud` |
| VNC viewer в UI | selenoid-ui | Сессия с `enableVNC` → открыть VNC в dashboard |
| Video playback | selenoid-ui | Сессия с `enableVideo` → `/video/` в UI |
| CM install на чистый хост | cm + autotests-cloud | `deploy/deploy.sh` / Actions deploy |
| Полный hub pyramid локально | этот репо | `./gradlew testHubAll -DskipHealthCheck=true` |

### Deploy triggers (`repository_dispatch`)

После docker-push в `release.yml`:

| Репо | Secret | event-type |
|------|--------|------------|
| [qa-guru/selenoid](https://github.com/qa-guru/selenoid) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api,smoke` |
| [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api,smoke` |
| [qa-guru/cm](https://github.com/qa-guru/cm) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api` |
| [qa-guru/browser-image](https://github.com/qa-guru/browser-image) `publish.yml` | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `playwright` (`source_variant=playwright`) |
| [qa-guru/browser-image](https://github.com/qa-guru/browser-image) `publish-webdriver.yml` | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `smoke` → `testWebdriverE2e` (`source_variant=webdriver`) |

Payload: `source_repo`, `source_ref`, `source_version`, `test_tags`, опционально `source_variant` (`playwright` \| `webdriver`).  
TestOps launch name: `Deploy smoke — {source_repo} {source_version} #{run}`.

Ручная проверка:

```bash
gh api repos/qa-guru/selenoid-tests/dispatches --input - <<'EOF'
{"event_type":"deploy-smoke","client_payload":{"source_repo":"qa-guru/selenoid","source_version":"manual","test_tags":"api"}}
EOF
```

### Dashboard (`allurerc.mjs`)

Пирамида: `unit → component → integration → api → e2e → manual`.

## Пробелы → закрыто (2026-07)

| Сервис | unit | component | integration | api | e2e | Добавлено |
|--------|------|-----------|-------------|-----|-----|-----------|
| **cm** | ✓ | **+4** | **+1** | **+3** | ✓ | version/help fixtures; CI job `java-cm` |
| **playwright-image** (`browser-image/playwright/`) | ✓ | **+4** | **+3** | ✓ | ✓ | +firefox/webkit WS (local-only) |
| **webdriver-image** (`browser-image/webdriver/`) | — | ✓ (+min) | ✓ | ✓ | ✓ | WD status + session API; chrome + chrome-min integration |
| **selenoid** | ✓ | **+2** | **+1** | **+2** | ✓¹ | logs, status+session, HubStatusParserTest |
| **selenoid-ui** | ✓ | ✓ | **+1** | ✓ | **+1** | browsers-config integration, sessions list e2e |

CM api: `./gradlew testCmApi -DpyramidStand=selenoid_github -DskipHealthCheck=true` (after `scripts/start-ci-cm-stack.sh`).

¹ **selenoid e2e:** `@Component("selenoid")` — нет отдельного e2e-класса; сквозной путь hub покрыт `HubSessionTests` / `HubPlaywrightSessionTests` (`webdriver-image` / `playwright-image`, `@Layer e2e`).

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
| CmBrowsersConfigJsonTest | cm | cm | component | — |
| CmStatusOutputTest | cm | cm | component | — |
| CmHubStatusApiTests | cm | cm | api | api, cm |
| CmHubSessionApiTests | cm | cm | api | api, cm |
| CmUiStatusApiTests | cm | cm | api | api, cm |
| ConfigOwnerMergeTest | — | — | unit | — |
| HubStatusJsonTest | selenoid | selenoid | component | — |
| HubStatusParserTest | selenoid | — | component | — |
| HubStatusBrowsersJsonTest | selenoid | selenoid | component | — |
| UiStatusJsonTest | selenoid-ui | selenoid-ui | component | — |
| UiPingJsonTest | selenoid-ui | selenoid-ui | component | — |
| UiErrorJsonTest | selenoid-ui | selenoid-ui | component | — |
| SseStateJsonTest | selenoid-ui | selenoid-ui | component | — |
| SseErrorsJsonTest | selenoid-ui | selenoid-ui | component | — |
| SessionCreateJsonTest | selenoid | selenoid | component | — |
| BrowsersConfigJsonTest | selenoid-ui | selenoid-ui | component | — |
| HubStatusForwardCompatJsonTest | selenoid | selenoid | component | — |
| HubLogsListJsonTest | selenoid | selenoid | component | — |
| HubPingTests | selenoid | selenoid | api | api |
| HubStatusBrowsersTests | selenoid | selenoid | api | api |
| HubSessionDeleteUnknownTests | selenoid | selenoid | api | api, negative |
| HubSessionInvalidBrowserTests | selenoid | selenoid | api | api, negative |
| UiBrowsersConfigTests | selenoid-ui | selenoid-ui | api | api |
| UiStatusTotalTests | selenoid-ui | selenoid-ui | api | api |
| UiPingVersionTests | selenoid-ui | selenoid-ui | api | api |
| UiSseMultipleEventsTests | selenoid-ui | selenoid-ui | api | api |
| PlaywrightUnknownPathTests | playwright-image | playwright-image | api | api, negative |
| PlaywrightWsPathJsonTest | playwright-image | playwright-image | component | — |
| PlaywrightBrowserCapsJsonTest | playwright-image | playwright-image | component | — |
| PlaywrightMinCatalogJsonTest | playwright-image | playwright-image | component | min |
| HubStatusTests | selenoid | selenoid | api | api |
| HubLogsListApiTests | selenoid | selenoid | api | api, negative |
| HubStatusSessionApiTests | selenoid | selenoid | api | api |
| HubSessionApiTests | selenoid | selenoid | api | api |
| PlaywrightEndpointTests | playwright-image | playwright-image | api | api |
| UiStatusTests | selenoid-ui | selenoid-ui | api | api |
| UiSseStreamTests | selenoid-ui | selenoid-ui | api | api |
| UiPingTests | selenoid-ui | selenoid-ui | api | api |
| HubPlaywrightSessionTests | playwright-image | playwright-image | integration | integration |
| HubPlaywrightMinSessionTests | playwright-image | playwright-image | integration | integration, min |
| HubPlaywrightNavigateTests | playwright-image | playwright-image | integration | integration |
| UiStatusWithSessionTests | selenoid-ui | selenoid-ui | integration | integration |
| UiSseWithSessionTests | selenoid-ui | selenoid-ui | integration | integration |
| StackHealthTests | selenoid-ui | selenoid-ui | integration | integration |
| CmInstallerLifecycleTests | cm | CM | integration | integration, cm |
| CmInstallerSessionTests | cm | CM | e2e | smoke, cm |
| UiHubStatusConsistencyTests | selenoid-ui | selenoid-ui | integration | integration |
| UiStatusWhenHubDownTests | selenoid-ui | selenoid-ui | integration | integration, local-only |
| UiBrowsersConfigIntegrationTests | selenoid-ui | selenoid-ui | integration | integration |
| HubStatusSessionIntegrationTests | selenoid | selenoid | integration | integration |
| HubChromeSessionIntegrationTests | webdriver-image | webdriver-image | integration | integration |
| HubChromeMinSessionTests | webdriver-image | webdriver-image | integration | integration, min |
| UiStatusRecoveryTests | selenoid-ui | selenoid-ui | integration | integration, resilience, local-only |
| HubSessionTests | webdriver-image | webdriver-image | e2e | smoke |
| HubSessionIdTests | webdriver-image | webdriver-image | e2e | smoke |
| HubSessionHeadingTests | webdriver-image | webdriver-image | e2e | smoke |
| HubSessionTitleTests | webdriver-image | webdriver-image | e2e | smoke |
| HubCapabilitiesApiTests | selenoid | selenoid | api | api |
| HubWebDriverStatusApiTests | selenoid | selenoid | api | api |
| WebDriverStatusApiTests | webdriver-image | webdriver-image | api | api |
| WebDriverSessionApiTests | webdriver-image | webdriver-image | api | api |
| HubWebDriverStatusJsonTest | selenoid | selenoid | component | — |
| ChromeMinCatalogJsonTest | webdriver-image | webdriver-image | component | min |
| CmVersionOutputTest | cm | cm | component | — |
| CmHelpOutputTest | cm | cm | component | — |
| CmCliVersionTests | cm | cm | integration | cm |
| HubPlaywrightFirefoxSessionTests | playwright-image | playwright-image | integration | playwright, local-only |
| HubPlaywrightWebkitSessionTests | playwright-image | playwright-image | integration | playwright, local-only |
| HubPlaywrightSessionTests | playwright-image | playwright-image | e2e | playwright, smoke |
| HubPlaywrightMinSessionTests | playwright-image | playwright-image | e2e | min |
| UiStatusBarTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiDashboardLoadTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiSseIndicatorTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiReloadTests | selenoid-ui | selenoid-ui | e2e | smoke |
| UiSessionsListTests | selenoid-ui | selenoid-ui | e2e | smoke |

## Config keys

| Key | Default |
|-----|---------|
| apiBaseUrl | `""` (→ hubUrl) |
| hubStatusPath | `/status` (prod cloud: `/hub/status` — raw hub via nginx) |
| hubUrl | http://127.0.0.1:4444/ |
| uiUrl | http://127.0.0.1:8080/ |
| remoteUrl | http://127.0.0.1:4444/wd/hub |
| cmHubPort | 4445 (CM installer; dev hub stays :4444) |
| cmUiPort | 8081 |
| playwrightWsEndpoint | ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1 |

Override: `-DhubUrl`, `-DuiUrl`, `-DapiBaseUrl`, …
