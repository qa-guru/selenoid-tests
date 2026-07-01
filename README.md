# selenoid-tests

Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **playwright-image** — Go unit (в CI из исходных репо) + Java e2e/integration/api.

- Allure TestOps: проект `selenoid-tests`, **ALLURE_PROJECT_ID=5271**
- Allure 3 GitHub Pages: `https://qa-guru.github.io/selenoid-tests/reports/<run-id>/`
- Dashboard: `.../reports/<run-id>/dashboard/index.html`

## Prerequisite (локально)

```bash
cd ../dev
./scripts/build-selenoid-ui.sh
./scripts/start-selenoid.sh &
./scripts/start-selenoid-ui.sh &
```

## Gradle pyramid slices

```bash
./gradlew testUnit -DskipHealthCheck=true                              # unit
./gradlew testApi -DskipHealthCheck=true                             # @Layer api
./gradlew testIntegration -DskipHealthCheck=true                     # @Layer integration (без local-only)
./gradlew testE2e -DskipHealthCheck=true                             # e2e smoke
./gradlew testPlaywright -DskipHealthCheck=true                      # Playwright WS
./gradlew testResilience -DskipHealthCheck=true                      # hub kill/recovery (local-only)

# CI-эквиваленты
./gradlew test -Denv=selenoid_github_api -DincludeTags=api -DskipHealthCheck=true
./gradlew test -Denv=selenoid_github_integration -DincludeTags=integration -DskipHealthCheck=true
./gradlew test -DincludeTags=smoke,api -DexcludeTags=integration,resilience,local-only,playwright

./gradlew allureReport
```

Stand override: `-DpyramidStand=selenoid_github` → env `selenoid_github_api`, `selenoid_github_integration`, …

## CI

Workflow: `.github/workflows/selenoid-tests.yml`

| Job | Что делает |
|-----|------------|
| `go-unit` (matrix) | Checkout `qa-guru/selenoid`, `selenoid-ui`, `cm` → Go unit → Allure |
| `java-e2e` | `./gradlew test` с tags `smoke,api` (integration excluded) |
| `report` | Merge `build/allure-results/**` → `allureReport` → gh-pages → TestOps 5271 |

`workflow_dispatch`: `test_tags=integration` для integration slice; `env_profile=selenoid_github_api` для api-only.

### Dashboard (`allurerc.mjs`)

Пирамида: `unit → component → integration → api → e2e → manual`.

## Матрица

| Класс | Сервис | @Epic | @Layer | @Tag |
|-------|--------|-------|--------|------|
| ConfigReaderTest | — | — | unit | — |
| HubStatusTests | selenoid | selenoid | api | api |
| HubSessionApiTests | selenoid | selenoid | api | api |
| PlaywrightEndpointTests | playwright-image | playwright-image | api | api |
| UiStatusTests | selenoid-ui | selenoid-ui | api | api |
| UiSseStreamTests | selenoid-ui | selenoid-ui | api | api |
| UiPingTests | selenoid-ui | selenoid-ui | api | api |
| UiHubStatusConsistencyTests | selenoid-ui | selenoid-ui | integration | integration |
| UiStatusWhenHubDownTests | selenoid-ui | selenoid-ui | integration | integration, local-only |
| UiStatusRecoveryTests | selenoid-ui | selenoid-ui | integration | integration, resilience, local-only |
| HubSessionTests | selenoid | selenoid | e2e | smoke |
| HubPlaywrightSessionTests | playwright-image | playwright-image | e2e | playwright, smoke |
| UiStatusBarTests | selenoid-ui | selenoid-ui | e2e | smoke |

## Config keys

| Key | Default |
|-----|---------|
| apiBaseUrl | `""` (→ hubUrl) |
| hubUrl | http://127.0.0.1:4444/ |
| uiUrl | http://127.0.0.1:8080/ |
| remoteUrl | http://127.0.0.1:4444/wd/hub |
| playwrightWsEndpoint | ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1 |

Override: `-DhubUrl`, `-DuiUrl`, `-DapiBaseUrl`, …
