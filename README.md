# selenoid-tests

<p align="center">
  <img src="docs/logo.svg" alt="Selenoid Stack logo" width="120">
</p>

Go e2e/integration orchestrator for the Selenoid stack — merged Allure pyramid across hub, UI, cm and browser nodes (Go unit + Go hub/cm slices). Prod smoke target: [selenoid.qa.guru](https://selenoid.qa.guru) (Selenoid 3).

[![Selenoid Stack](https://qa-guru.github.io/selenoid-tests/readme/badge.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![Selenoid Stack stats](https://qa-guru.github.io/selenoid-tests/readme/stats.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![Selenoid Stack metrics](https://qa-guru.github.io/selenoid-tests/readme/metrics-panel.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

## Automated Tests Dashboard

Live SVG metrics + Allure 3 dashboard (pyramid tile **testingPyramid**), updated after each orchestrator run on `main`:

<a href="https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://qa-guru.github.io/selenoid-tests/readme/dashboard-preview-dark.png">
    <img
      src="https://qa-guru.github.io/selenoid-tests/readme/dashboard-preview.png"
      alt="Allure 3 dashboard — pyramid, status dynamics, success distribution"
      width="800"
    />
  </picture>
</a>

Dashboard PNG updates after each orchestrator run on `main` (Playwright screenshot of Allure 3 dashboard).

| Link | Description |
|------|-------------|
| [Dashboard](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/) | Pyramid: Go unit ×3 + Go hub + Go CM |
| [Awesome](https://qa-guru.github.io/selenoid-tests/reports/latest/awesome/) | Epic drill-down: selenoid, selenoid-ui, cm, webdriver-image, playwright-image, video-recorder |
| [TestOps project](https://allure.autotests.cloud/project/5271) | Cloud launches |
| [CI workflow](https://github.com/qa-guru/selenoid-tests/actions/workflows/selenoid_github-orchestrator.yml) | `workflow_dispatch` max run: `env_profile=selenoid_github_e2e` |

Per-component badges: `readme/badge-{selenoid,selenoid-ui,cm,webdriver-image,playwright-image,video-recorder}.svg` — for hub repo READMEs.

<!-- stack-branches-note:start -->
> ## Стабильные билды
>
> **Prod:** [selenoid.qa.guru](https://selenoid.qa.guru) — **Selenoid 3** (hub/cm/UI v3.x на `main`). Go pyramid на `main` — gates `hub-prod` / `hub-all`; pin-ветки 2.x — frozen rollback.
>
> | Ветка | Semver | Назначение |
> |-------|--------|------------|
> | **`main`** | **v3.0.0+** | Активная prod-линия + CI orchestrator (Go) |
> | `selenoid2-1.55-…-react18` | v2.3.0 | frozen maintenance pin |
> | `selenoid2-1.45-…-react16` | v2.2.1 | frozen rollback reference |
>
> Точные версии pin-веток — `STACK-PIN.md` на соответствующей ветке. Monorepo SSOT: [`projects/selenoid-home/README.md`](https://github.com/qa-guru/zero-design-system/blob/master/projects/selenoid-home/README.md).
<!-- stack-branches-note:end -->


Центральный репозиторий автотестов Selenoid-стека: [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests).

Покрывает **selenoid**, **selenoid-ui**, **cm**, **browser-image** (`playwright/` + `webdriver/`) — **Go autotests** (hub pyramid + product unit из исходных репо).

**Автотесты = Go:** root module `github.com/qa-guru/selenoid-tests` — ADR [`docs/ADR-go-pyramid.md`](docs/ADR-go-pyramid.md). Gates: `hub-prod` (−cm/min/resilience + `har-prod`), `hub-all` (github full −cm). Slices: `unit|api|integration|ui|webdriver|playwright|e2e|har-prod|min|resilience|cm|warm-pool`. CI gate: `go-hub` + `go-cm` (+ `go-unit` matrix). **`@Layer("unit")`:** `internal/*` + `tests/unit/fixture/` (JSON parsers). **`@Layer("component")`:** только RTL (`selenoid-ui/ui`).

**Warm-pool:** live orchestrator + hub-attach — slice `warm-pool` (`tests/e2e/warmpool`, local stand `:9090`). **Не** в `hub-prod` / `hub-all` (CI/prod без loopback attach). Unit: `internal/warmpool` + fixture contract.

## Экосистема qa-guru Selenoid

| Ресурс | Ссылка | Роль |
|--------|--------|------|
| selenoid | [github.com/qa-guru/selenoid](https://github.com/qa-guru/selenoid) | Hub |
| selenoid-ui | [github.com/qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui) | Web UI |
| cm | [github.com/qa-guru/cm](https://github.com/qa-guru/cm) | Установщик |
| browser-image | [github.com/qa-guru/browser-image](https://github.com/qa-guru/browser-image) | Docker browser nodes |
| **selenoid-tests** (этот) | [github.com/qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests) | E2e/integration ethalon |
| Docker Hub | [hub.docker.com/u/qaguru](https://hub.docker.com/u/qaguru) | Образы `qaguru/*` |

**Stack:** prod [selenoid.qa.guru](https://selenoid.qa.guru) — **Selenoid 3** (hub/cm/UI v3.x on `main`). Pin-ветки 2.x — frozen rollback in git only. Browser-image — image tags, не git semver.

- Allure TestOps: проект `selenoid-tests`, **ALLURE_PROJECT_ID=5271**
- Test layers: `@Layer` keys → TestOps mapping (`e2e` → E2E Tests) — RAG `test-layers`, sync: `qa-guru-tms-automator/scripts/sync_testops_layer_mappings.py`
- Component filter: `@Component` label → TestOps custom field **Component** (`cm`, `selenoid`, `selenoid-ui`, `playwright-image`, `webdriver-image`, `video-recorder`); sync: `qa-guru-tms-automator/scripts/sync_testops_component_mappings.py`
- Allure 3 GitHub Pages: `https://qa-guru.github.io/selenoid-tests/reports/<run-id>/` (per-run URLs retained on gh-pages — CI `keep_files: true`)
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

CM integration (`cm` slice): Docker + `cm` binary; hub/UI на **:4445/:8081** (dev-стек на :4444/:8080 не конфликтует).

```bash
cd ../cm && go build -o cm .
cd ../dev && ./scripts/build-selenoid-ui.sh   # ui/build для cross-compile selenoid-ui
./scripts/prepare-ci-cm-workspace.sh
./scripts/start-ci-cm-stack.sh &
SELENOID_TEST_ENV=selenoid_github_cm_integration ./scripts/run-go-pyramid.sh cm
```

## Go pyramid slices

```bash
# Offline
SELENOID_TEST_ENV=local_unit ./scripts/run-go-pyramid.sh unit
SELENOID_TEST_ENV=local_unit ./scripts/run-go-pyramid.sh unit

# CI push gate (github stack, −cm; +min/resilience)
PYRAMID_STAND=selenoid_github ./scripts/run-go-pyramid.sh hub-all

# Prod cloud gate (−cm/−min/−resilience + har-prod)
PYRAMID_STAND=selenoid_qa_guru ./scripts/run-go-pyramid.sh hub-prod

# Single slices (stand from PYRAMID_STAND or SELENOID_TEST_ENV)
./scripts/run-go-pyramid.sh api
./scripts/run-go-pyramid.sh integration
./scripts/run-go-pyramid.sh ui
./scripts/run-go-pyramid.sh webdriver
./scripts/run-go-pyramid.sh playwright
./scripts/run-go-pyramid.sh e2e
./scripts/run-go-pyramid.sh min
./scripts/run-go-pyramid.sh resilience
./scripts/run-go-pyramid.sh cm
./scripts/run-go-pyramid.sh har-prod

# Allure report (local, after a run)
./scripts/allure-report.sh generate
```

Stand override: `PYRAMID_STAND=selenoid_github` → env `selenoid_github_api`, `selenoid_github_integration`, …  
Profile override: `SELENOID_TEST_ENV=selenoid_qa_guru_api ./scripts/run-go-pyramid.sh api`

### Prod hub (`selenoid.qa.guru`)

Profiles: `selenoid_qa_guru_api`, `selenoid_qa_guru_e2e` — remote hub `https://selenoid.qa.guru` (auth `user1:1234` in properties; e2e `uiUrl` embeds credentials for Capabilities create-session XHR).

```bash
SELENOID_TEST_ENV=selenoid_qa_guru_api ./scripts/run-go-pyramid.sh api
PYRAMID_STAND=selenoid_qa_guru ./scripts/run-go-pyramid.sh hub-prod
SELENOID_TEST_ENV=selenoid_qa_guru_e2e ./scripts/run-go-pyramid.sh har-prod
# UI smoke subset:
SELENOID_TEST_ENV=selenoid_qa_guru_e2e ./scripts/run-go-pyramid.sh ui
```

Post-deploy: `selenoid.qa.guru` → Actions → `trigger-deploy-smoke` → `repository_dispatch deploy-smoke` → this repo (`skip_go_unit`, `env_profile=selenoid_qa_guru_api`, default tags **`api` only**).

**Release smoke vs deploy smoke:** `selenoid-ui` `release.yml` dispatches `api,smoke` **after** polling prod `/ui/status` for the new UI pin (`wait_for_prod_ui_version`). Image publish alone is not a deploy — pre-pin smoke against transitional/old prod is a race, not a product regression. Post-deploy trigger from `selenoid.qa.guru` remains the hub/api gate (tags `api` by default; UI e2e gate is the waited release-smoke).

Prod caveats (nginx): hub API uses raw `GET /hub/status` via `hubStatusPath` (not UI `/status` with `.state`). `GET /logs/{id}` — nginx → hub (auth); UI uses `/ws/logs/{id}`.

UI e2e canon (v3): root → `#/statistics`; Sessions archive (ex-Videos) → `#/sessions` + `.archive__list`; New Session (ex-Capabilities) → `#/new-session`; status tiles → `Connected` / `Issue` / `Unknown`; VNC → `[data-testid=vnc-window].vnc-window--connected`.

### `playwright` slice prerequisite

Hub на `:4444` и Docker-образ из `fixtures/ci-browsers.json` / `dev/browsers.json`:

```bash
cd dev && ./scripts/start-selenoid.sh &
docker pull qaguru/playwright-chromium:1.61.1   # или ./scripts/pull-browser-images.sh
PYRAMID_STAND=selenoid_github ./scripts/run-go-pyramid.sh playwright
```

В CI `scripts/start-ci-selenoid-stack.sh` тянет образы из `fixtures/ci-browsers.json` (chrome + firefox + msedge warm + playwright-chromium).

### Playwright-chromium-min (`1.61.1-min`)

Образ в `fixtures/ci-browsers.json`; endpoint — `selenoid_github_min_integration.properties` (VNC/video off). Входит в `hub-all` (`min` slice).

```bash
./scripts/run-go-pyramid.sh min
```

| Package | Layer | Tag |
|---------|-------|-----|
| `tests/unit/fixture/min` | unit | min |
| `tests/integration/min` | integration | min |

## CI

Workflow: `.github/workflows/selenoid_github-orchestrator.yml` (`name: selenoid-tests Tests`).

| Job | Что делает |
|-----|------------|
| `go-unit` (matrix) | Checkout `qa-guru/selenoid`, `selenoid-ui`, `cm` → Go unit → Allure |
| `go-hub` | **CI gate** — `run-go-pyramid.sh`; push/github → `hub-all`; prod profile → `hub-prod`; dispatch `test_tags` → slice |
| `go-cm` | Push (non-prod): `run-go-pyramid.sh cm` (CM :4445/:8081) |
| `report` | Merge `build/allure-results/**` → Allure 3 → gh-pages → TestOps 5271 |

```bash
# Go unit+unit JSON parsers (offline, −cm)
SELENOID_TEST_ENV=local_unit ./scripts/run-go-pyramid.sh unit
SELENOID_TEST_ENV=local_unit ./scripts/run-go-pyramid.sh unit
# Go api (prod cloud)
SELENOID_TEST_ENV=selenoid_qa_guru_api ./scripts/run-go-pyramid.sh api
# Go prod gate (hub-prod semantics + har-prod)
PYRAMID_STAND=selenoid_qa_guru ./scripts/run-go-pyramid.sh hub-prod
# Go CI push gate (github stack, −cm)
PYRAMID_STAND=selenoid_github ./scripts/run-go-pyramid.sh hub-all
# Go HAR prod smoke
SELENOID_TEST_ENV=selenoid_qa_guru_e2e ./scripts/run-go-pyramid.sh har-prod
# Warm-pool live stand :9090 (hub-attach skips unless slots + -warm-pool-url)
./scripts/run-go-pyramid.sh warm-pool
```

`workflow_dispatch`: `test_tags=integration|api|smoke|playwright` → Go slice; `env_profile=selenoid_qa_guru_*` → `hub-prod`.

### Component × Layer × CI (push `main`)

Пирамида: `unit → component → integration → api → e2e → manual`. **Go hub** — пакеты ниже; **Go unit** — отдельно в `go-unit`.  
**Матрица (Selenoid 3 / `main`):** hub/ui/cm/browser-image покрыты Go pyramid + product Go unit. `warm-pool` — local slice, не CI push.

| Component | unit | component | integration | api | e2e | manual | CI push |
|-----------|:----:|:---------:|:-----------:|:---:|:---:|:------:|---------|
| **selenoid** | Go | — | Go | Go | Go (wd/pw) | — | `go-unit` + `go-hub` |
| **selenoid-ui** | Go | RTL | Go | Go | Go | —⁶ | `go-unit` + `go-hub` |
| **cm** | Go | — | Go | Go | Go | — | `go-unit` + `go-cm` |
| **playwright-image** | Go¹ | — | Go | Go | Go | — | `go-hub` |
| **webdriver-image** | Go | — | Go | Go | Go | — | `go-hub` |
| **video-recorder** | — | Go | — | Go | — | — | `go-hub` (api slice / dispatch) |
| **warm-pool** | Go⁸ | — | — | Go⁸ | Go⁸ | — | — (slice `warm-pool`, skip if stand down) |
| **dev** | — | —² | —³ | — | — | ✓ | — |
| **selenoid-qa-guru** | — | — | — | —⁴ | —⁵ | ✓ | deploy-smoke dispatch |

¹ **selenoid e2e:** отдельного `@Component("selenoid")` e2e-пакета нет; сквозной hub-path покрыт `tests/e2e/webdriver` и `tests/e2e/playwright` в `hub-all`. JSON fixture parsers — **`@Layer("unit")`**, `tests/unit/fixture/`, не component.  
² **dev unit (JSON fixtures):** `browsers.json` SSOT — `tests/unit/fixture/*CatalogJsonTest`, `BrowsersConfigJsonTest`.  
³ **dev integration:** `start-ci-selenoid-stack.sh` — оркестрация CI, не test-class.  
⁴ **cloud api:** post-deploy `selenoid_qa_guru_api` через `trigger-deploy-smoke` / `repository_dispatch` — не локальный класс в этой матрице.  
⁵ **cloud e2e:** профиль `selenoid_qa_guru_e2e` — manual / расширенный deploy-smoke.  
⁶ **selenoid-ui manual:** video playback — runbook (ниже); VNC viewer — Go `tests/e2e/ui` (prod profile `selenoid_qa_guru_e2e`).  
⁷ **playwright-image / webdriver-image unit:** catalog JSON + session body — `tests/unit/fixture/` + `internal/config` (`@Layer("unit")`).  
⁸ **warm-pool:** `internal/warmpool` + fixture contract (unit); live GET/reserve/release (api) and hub-attach (e2e) in `tests/e2e/warmpool`. Stand `python scripts/stands/ensure.py selenoid-warm-pool`. Hub-attach skips unless hub has `warmTotal>0` and ChromeDriver on loopback is dialable.

### Manual (runbook)

| Сценарий | Где | Как |
|----------|-----|-----|
| Локальный стек hub/UI | `../dev/README.md` | `build-selenoid*.sh` + `start-selenoid*.sh` |
| Prod hub smoke | `selenoid-qa-guru` | `./deploy/smoke-remote.sh https://selenoid.qa.guru` |
| VNC viewer в UI | selenoid-ui | Go `tests/e2e/ui` (`run-go-pyramid.sh ui`) |
| Video playback | selenoid-ui | Сессия с `enableVideo` → `/video/` в UI |
| CM install на чистый хост | cm + autotests-cloud | `deploy/deploy.sh` / Actions deploy |
| Полный hub pyramid локально | этот репо | `PYRAMID_STAND=selenoid_github ./scripts/run-go-pyramid.sh hub-all` |

### Deploy triggers (`repository_dispatch`)

После docker-push в `release.yml`:

| Репо | Secret | event-type |
|------|--------|------------|
| [qa-guru/selenoid](https://github.com/qa-guru/selenoid) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api,smoke` |
| [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api,smoke` |
| [qa-guru/cm](https://github.com/qa-guru/cm) | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `api` |
| [qa-guru/browser-image](https://github.com/qa-guru/browser-image) `publish.yml` | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `playwright` (`source_variant=playwright`) |
| [qa-guru/browser-image](https://github.com/qa-guru/browser-image) `publish-webdriver.yml` | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `smoke` → browser slice (`source_variant=webdriver`, `source_browser`, `webdriver_variant`) |
| [qa-guru/browser-image](https://github.com/qa-guru/browser-image) `publish-video-recorder.yml` | `SELENOID_TESTS_DISPATCH_TOKEN` | `deploy-smoke` | `smoke` → `testVideoRecorder` (`source_variant=video-recorder`) |

Payload: `source_repo`, `source_ref`, `source_version`, `test_tags`, опционально `source_variant` (`playwright` \| `webdriver` \| `video-recorder`), `source_browser` (`chrome` \| `firefox` \| `msedge`), `webdriver_variant` (`warm` \| `min`).  
WebDriver dispatch: chrome warm → Go `webdriver`; chrome min → Go `integration/min`; firefox/msedge → Go `tests/integration/wd` browser slice.
TestOps launch name: `Deploy smoke — {source_repo} {source_version} #{run}`.

Ручная проверка:

```bash
gh api repos/qa-guru/selenoid-tests/dispatches --input - <<'EOF'
{"event_type":"deploy-smoke","client_payload":{"source_repo":"qa-guru/selenoid","source_version":"manual","test_tags":"api"}}
EOF
```

### Dashboard (`allurerc.mjs`)

Пирамида: `unit → component → integration → api → e2e → manual`.

## Go package map

| Package / slice | Layer | Notes |
|-----------------|-------|-------|
| `internal/config`, `internal/helpers`, `internal/hubapi` | unit | offline pure logic |
| `tests/unit/fixture`, `tests/unit/fixture/min` | unit | JSON fixture parsers (`@Tag("min")` on min catalog) |
| `tests/integration/wd`, `tests/integration/ui`, `tests/integration/pw`, `tests/integration/min`, `tests/integration/resilience` | integration | hub stack |
| `tests/api` | api | hub + UI + video-recorder |
| `tests/e2e/ui`, `tests/e2e/webdriver`, `tests/e2e/playwright`, `tests/e2e/har` | e2e | browser / HAR |
| `tests/cm/...` | cm pyramid | `:4445/:8081`; CI `go-cm` |
| `scripts/run-go-unit.sh` | product unit | selenoid / selenoid-ui / cm repos |

CM api locally: `./scripts/start-ci-cm-stack.sh` then `SELENOID_TEST_ENV=selenoid_github_cm_integration ./scripts/run-go-pyramid.sh cm`.

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
| browser / browserVersion | chrome / **149.0** (warm chrome API + e2e) |
| chromeVersion | 149.0 |
| chromeMinVersion | 149.0-min |
| firefoxVersion | 151.0 |
| firefoxMinVersion | 151.0-min |
| msedgeVersion | 145.0 |
| msedgeMinVersion | 145.0-min |

Override: env `SELENOID_TEST_*`, plain `hubUrl`/`uiUrl`/… in process env, or keys in `src/test/resources/config/*.properties`.
