# selenoid-tests

Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **playwright-image** — Go unit (в CI из исходных репо) + Java e2e/integration.

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

## Gradle slices

```bash
./gradlew test --tests 'config.*Test' -DskipHealthCheck=true          # unit (CI bootstrap)
./gradlew test -DincludeTags=api -DskipHealthCheck=true              # HTTP integration
./gradlew test -DincludeTags=smoke -DexcludeTags=playwright,resilience,local-only
./gradlew test -DincludeTags=playwright                               # Playwright WS
./gradlew test -DincludeTags=resilience                               # hub kill/recovery (local)

./gradlew allureReport
```

## CI

Workflow: `.github/workflows/selenoid-tests.yml`

| Job | Что делает |
|-----|------------|
| `go-unit` (matrix) | Checkout `qa-guru/selenoid`, `selenoid-ui`, `cm` → Go unit → Allure |
| `java-e2e` | `./gradlew test` с tags `smoke,api`, `-DskipHealthCheck=true` |
| `report` | Merge `build/allure-results/**` → `allureReport` → gh-pages → TestOps 5271 |

Триггеры:

- `push` → `main`: go-unit + java (smoke,api)
- `workflow_dispatch`: optional `test_class`, `test_tags`, `source_ref` (ref Go-репо)
- `repository_dispatch` (`deploy-smoke`): post-deploy smoke

### Go unit → Allure (вариант A)

**Выбран вариант A:** `gotestsum --junitfile` + `@allurereport/reader` (JUnit XML plugin Allure 3).

Скрипты:

- `scripts/run-go-unit.sh <repo> <epic>` — запуск `ci/test.sh`-эквивалента + конвертация
- `scripts/junit-to-allure.mjs` — JUnit XML → `*-result.json` с labels `epic`, `layer=unit`

На каждый Go-репо проставляются:

| repo | epic | layer |
|------|------|-------|
| selenoid | `selenoid` | `unit` |
| selenoid-ui | `selenoid-ui` | `unit` |
| cm | `cm` | `unit` |

Вариант B (`allurectl watch -- go test -json`) не используется: нет нативного Allure-адаптера для Go, JUnit-путь проще и совместим с merge в Gradle report.

### Java @Epic (lowercase)

| Epic | Тесты |
|------|-------|
| `selenoid` | HubStatusTests, HubSessionTests |
| `selenoid-ui` | Ui*Tests |
| `playwright-image` | HubPlaywrightSessionTests |

### Dashboard (`allurerc.mjs`)

1. Общий `statusDynamics` (limit 30)
2. `currentStatus` с `groupBy: epic`
3. `stabilityDistribution` с `groupBy: epic`
4. Отдельный `statusDynamics` per epic (filter в dynamic config)

## Матрица

| Класс | Сервис | @Epic | @Tag |
|-------|--------|-------|------|
| ConfigReaderTest | — | — | unit |
| HubStatusTests | selenoid | selenoid | api |
| HubSessionTests | selenoid | selenoid | smoke |
| HubPlaywrightSessionTests | playwright-image | playwright-image | playwright, smoke |
| UiStatusTests | selenoid-ui | selenoid-ui | api |
| UiSseStreamTests | selenoid-ui | selenoid-ui | api |
| UiStatusBarTests | selenoid-ui | selenoid-ui | smoke |
| UiStatusRecoveryTests | selenoid-ui | selenoid-ui | resilience, local-only |

## Config keys

| Key | Default |
|-----|---------|
| hubUrl | http://127.0.0.1:4444/ |
| uiUrl | http://127.0.0.1:8080/ |
| remoteUrl | http://127.0.0.1:4444/wd/hub |
| playwrightWsEndpoint | ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1 |

Override: `-DhubUrl`, `-DuiUrl`, `-DplaywrightWsEndpoint`, …
