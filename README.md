# selenoid-tests

Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **playwright-image** (Java e2e/integration; Go unit — в CI, чат 2).

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

- `push` → `main`: unit/config tests (`config.*Test`, `-DskipHealthCheck=true`)
- `workflow_dispatch`: optional `test_class`, `test_tags`
- `repository_dispatch` (`deploy-smoke`): post-deploy smoke из selenoid/selenoid-ui/cm/playwright-image

## Матрица

| Класс | Сервис | @Tag |
|-------|--------|------|
| ConfigReaderTest | — | unit |
| HubStatusTests | selenoid | api |
| HubSessionTests | selenoid | smoke |
| HubPlaywrightSessionTests | playwright-image | playwright, smoke |
| UiStatusTests | selenoid-ui | api |
| UiSseStreamTests | selenoid-ui | api |
| UiStatusBarTests | selenoid-ui | smoke |
| UiStatusRecoveryTests | selenoid-ui | resilience, local-only |

## Config keys

| Key | Default |
|-----|---------|
| hubUrl | http://127.0.0.1:4444/ |
| uiUrl | http://127.0.0.1:8080/ |
| remoteUrl | http://127.0.0.1:4444/wd/hub |
| playwrightWsEndpoint | ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1 |

Override: `-DhubUrl`, `-DuiUrl`, `-DplaywrightWsEndpoint`, …
