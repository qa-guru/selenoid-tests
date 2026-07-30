# HAR completeness benchmark — SSOT

**Plan:** monorepo `.cursor/plans/har-completeness-benchmark.plan.md`

Committed marketing / verdict docs live here. Large run artifacts (`.har`, prod smoke reruns) stay under `build/har-compare/` (gitignored).

| File | Role |
|------|------|
| [MATRIX.md](MATRIX.md) | Scorecard matrix + allowed / forbidden claims |
| [matrix.json](matrix.json) | Machine-readable matrix |
| [ASSETS.md](ASSETS.md) | Marketing pack asset list |
| [TELEGRAM-DRAFT.txt](TELEGRAM-DRAFT.txt) | Draft post copy |
| [7-android-blocker.json](7-android-blocker.json) | Phase 6 cancelled record |
| [assets/](assets/) | Committed PNGs for landing / Telegram |

## Phases

| Phase | Status |
|-------|--------|
| 0–5 | green |
| 6 Android | **cancelled** (out of scope) |
| 7 Marketing pack | done — this directory |

## Regenerate `build/har-compare/`

```bash
python scripts/smoke-prod-har-content.py
python scripts/smoke-prod-har-ui.py
go test ./tests/e2e/har/... -count=1
```

Sync if matrix changed:

```bash
cp build/har-compare/{MATRIX.md,ASSETS.md,TELEGRAM-DRAFT.txt,matrix.json} docs/har-benchmark/
```

## CI gate

`go test ./internal/helpers/ -run TestHarBenchmarkAndroidCancelled` — phase 6 must stay **cancelled**.
