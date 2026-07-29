# ADR: Go test pyramid (replace Java)

**Status:** Accepted  
**Date:** 2026-07-29  
**Repo:** [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests)

## Context

Stack v3 target (`projects/selenoid-home/README.md`): autotests = **Go**.  
Current pyramid is Java/JUnit5/Selenide/RestAssured/Playwright-Java (~121 classes) + product Go unit via `scripts/run-go-unit.sh`.  
Prod gate: `testHubProd` on `selenoid_qa_guru_*` (unit→component→integration→api→e2e→webdriver→ui→playwright; no resilience/min/cm).

## Decision

1. **Root Go module** `github.com/qa-guru/selenoid-tests` — `go.mod` at repo root; packages `internal/…`, `tests/…`.
2. **Java is temporary** — keep `src/test` + Gradle until Go closes parity (`testHubProd` + local `testHubAll` / cm), then delete Java/Gradle from the gate and tree.
3. **Config SSOT** — reuse `src/test/resources/config/*.properties` (`env` / `SELENOID_TEST_ENV`, same keys as Owner `TestConfig`).
4. **HTTP** — stdlib `net/http` (+ thin helpers). No RestAssured/resty for hub/UI JSON.
5. **Browser** — `playwright-go` for UI e2e and Playwright WS slices (P3). WebDriver sessions stay raw WD HTTP (as today).
6. **Allure** — native `*-result.json` writer (`internal/allurex`) with labels `layer`, `component`, `epic`, `feature`, `story`, `tag`; `framework=go`, `language=go`. Merge into existing `report` job / TestOps **5271** / gh-pages.
7. **Slices** — `scripts/run-go-pyramid.sh` mirrors Gradle: `api`, `unit`, `component`, `integration`, `e2e`, `webdriver`, `ui`, `playwright`, `min`, `resilience`, `cm`, `hub-prod`, `hub-all`.
8. **CI dual-run** — job `go-hub` next to `java-e2e` until cutover; then Go becomes the gate and Java jobs are removed. `go-unit` (product repos) unchanged. CM stays local-only (not on prod tag lists).

## Layout

```
selenoid-tests/
  go.mod
  internal/config|httpx|hubapi|uiapi|allurex|health/
  tests/{unit,component,api,integration,e2e,webdriver,ui,playwright,cm}/
  scripts/run-go-pyramid.sh
  src/test/...          # Java until P5 delete
  build.gradle          # Java until P5 delete
```

## Non-goals

- Vitest/RTL in `selenoid-ui/ui` (stays).
- typescript-go for UI unit (separate track).
- Moving product hub/ui/cm `*_test.go` into this repo.

## Phases

| Phase | Scope |
|-------|--------|
| P0 | Config, health, allurex, CI `go-hub`, api vertical: Hub/UI status + UI ping |
| P1 | unit + component + full api (−cm) |
| P2 | integration (WD/PW sessions, UI status/SSE) on github + prod profiles |
| P3 | e2e UI + webdriver/playwright + HAR prod smoke |
| P4 | cm + resilience + min (local-only) |
| P5 | Gate → Go; remove Java; README «Автотесты» = Go |

## Consequences

- Dual-run until P5 (temporary cost).
- Properties path stays under `src/test/resources/config` until Java removal; then move to `config/` if desired (same keys).
- Quality gate: api/integration/e2e need ≥1 Allure step (`allure-reporting-requirements`).
