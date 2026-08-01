# ADR: Go test pyramid (replace Java)

**Status:** Implemented (P5 cutover 2026-07-29)  
**Date:** 2026-07-29  
**Repo:** [qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests)

## Context

Stack v3 (`projects/selenoid-home/README.md`): autotests = **Go**.  
Hub pyramid: root module `github.com/qa-guru/selenoid-tests` + product Go unit via `scripts/run-go-unit.sh`.  
Prod gate: `hub-prod` on `selenoid_qa_guru_*` (unit→integration→api→e2e→webdriver→ui→playwright→har-prod; no resilience/min/cm).

## Decision

1. **Root Go module** `github.com/qa-guru/selenoid-tests` — `go.mod` at repo root; packages `internal/…`, `tests/…`.
2. **Java removed (P5)** — hub pyramid gate is Go only; product Go unit via `scripts/run-go-unit.sh` unchanged.
3. **Config SSOT** — `src/test/resources/config/*.properties` (`SELENOID_TEST_ENV` / `env`, same keys as former Owner `TestConfig`).
4. **HTTP** — stdlib `net/http` (+ thin helpers). No RestAssured/resty for hub/UI JSON.
5. **Browser** — `playwright-go` for UI e2e and Playwright WS slices (P3). WebDriver sessions stay raw WD HTTP (as today).
6. **Allure** — native `*-result.json` writer (`internal/allurex`) with labels `layer`, `component`, `epic`, `feature`, `story`, `tag`; `framework=go`, `language=go`. Merge into existing `report` job / TestOps **5271** / gh-pages.
7. **Slices** — `scripts/run-go-pyramid.sh`: `unit` (`internal/*` + `tests/unit/fixture/`), `api`, `integration`, `e2e`, `webdriver`, `ui`, `playwright`, `min`, `resilience`, `cm`, `hub-prod`, `hub-all`. **`@Layer("component")`** — только RTL (Vitest); Go JSON parsers = **unit**.
8. **CI gate (P5)** — jobs `go-hub` + `go-cm` (+ `go-unit` from product repos). No Java dual-run. CM stays off prod tag lists (`hub-prod` / `selenoid_qa_guru_*`).

## Layout

```
selenoid-tests/
  go.mod
  internal/config|httpx|hubapi|uiapi|allurex|health/
  tests/unit/fixture/
  tests/{api,integration,e2e,webdriver,ui,playwright,cm}/...
  scripts/run-go-pyramid.sh
  src/test/resources/config/*.properties
  src/test/resources/fixtures/
```

## Non-goals

- Vitest/RTL in `selenoid-ui/ui` (stays).
- typescript-go for UI unit (separate track).
- Moving product hub/ui/cm `*_test.go` into this repo.

## Phases

| Phase | Scope |
|-------|--------|
| P0 | Config, health, allurex, CI `go-hub`, api vertical: Hub/UI status + UI ping |
| P1 | unit (+ JSON parser fixtures) + full api (−cm) |
| P2 | integration (WD/PW sessions, UI status/SSE) on github + prod profiles |
| P3 | e2e UI + webdriver/playwright + HAR prod smoke |
| P4 | cm + resilience + min (local-only) |
| P5 | Gate → Go; remove Java; README «Автотесты» = Go — **done** |

## Consequences

- Properties path stays under `src/test/resources/config` (same keys; Go loader in `internal/config`).
- Quality gate: api/integration/e2e need ≥1 Allure step (`allure-reporting-requirements`).
