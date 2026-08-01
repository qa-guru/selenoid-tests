# Stack pin: main / v3 (Selenoid 3)

**Репозиторий:** E2E orchestrator (qa-guru/selenoid-tests)

Этот файл на **`main`** описывает живой CI/test toolchain. Pin-ветки 2.x — отдельные frozen `STACK-PIN.md`.

| Поле | Значение |
|------|----------|
| Линия | Selenoid 3 |
| Stack semver | hub/cm/UI **v3.0.0+** (prod pin — deploy-чат) |
| Go | 1.26.5+ |
| Go (примечание) | Root module `github.com/qa-guru/selenoid-tests`; ADR [`docs/ADR-go-pyramid.md`](docs/ADR-go-pyramid.md) |
| Prod | [selenoid.qa.guru](https://selenoid.qa.guru) — smoke gates `hub-prod` / deploy-smoke |
| Git anchor | `main` |
| Matrix | Go pyramid: `unit` + `integration` + `api` + `e2e` + `cm` slices |
| CI gate | `go-hub` + `go-cm` + `go-unit` matrix → Allure 3 gh-pages + TestOps 5271 |

## Selenoid 2 maintenance pin (не путать)

Rollback / maintenance **v2.3.0** / React **18** — только pin-ветка
[`selenoid2-1.55-engine29.6-go1.26-react18`](https://github.com/qa-guru/selenoid-tests/tree/selenoid2-1.55-engine29.6-go1.26-react18).
Rollback **v2.2.1** / React 16 —
[`selenoid2-1.45-engine26.1-go1.26-react16`](https://github.com/qa-guru/selenoid-tests/tree/selenoid2-1.45-engine26.1-go1.26-react16).

См. также: [`projects/selenoid-home/README.md`](https://github.com/qa-guru/zero-design-system/blob/master/projects/selenoid-home/README.md) (monorepo SSOT).
