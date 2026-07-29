# HAR completeness benchmark — SSOT

**Plan:** [har-completeness-benchmark.plan.md](https://github.com/qa-guru/selenoid-tests/blob/main/docs/har-benchmark/MATRIX.md) (monorepo `.cursor/plans/har-completeness-benchmark.plan.md`)

Committed marketing / verdict docs live here. Large run artifacts (`.har`, prod smoke reruns) stay under `build/har-compare/` (gitignored; regenerate via compare e2e + prod smoke scripts).

| File | Role |
|------|------|
| [MATRIX.md](MATRIX.md) | Scorecard matrix + allowed / forbidden claims |
| [matrix.json](matrix.json) | Machine-readable matrix (incl. `6-android` not_claimed) |
| [ASSETS.md](ASSETS.md) | Marketing pack asset list |
| [TELEGRAM-DRAFT.txt](TELEGRAM-DRAFT.txt) | Draft post copy |
| [7-android.NOT_CLAIMED.txt](7-android.NOT_CLAIMED.txt) | Phase 6 human blocker summary |
| [7-android-blocker.json](7-android-blocker.json) | Phase 6 structured verdict |
| [assets/](assets/) | Committed PNGs for landing / Telegram |

## Phases

| Phase | Status |
|-------|--------|
| 0–5 | green (fixture + client + hub WD/PW meta/bodies) |
| 6 Android | **not_claimed** — no CDP `:7070` on `qaguru/android:*` |
| 7 Marketing pack | docs + assets in this directory |

## Regenerate `build/har-compare/`

From repo root (local hub + fixture on `:8080`):

```bash
# client + hub compare rows — see plan phase 1–5 tests/scripts
python scripts/smoke-prod-har-content.py   # prod step 5b → build/har-compare/prod-step5/
python scripts/smoke-prod-har-ui.py        # prod UI PNGs
go test ./tests/e2e/har/...                 # prod HarCapture 3b
```

After a full rerun, sync committed docs if matrix/claims changed:

```bash
cp build/har-compare/{MATRIX.md,ASSETS.md,TELEGRAM-DRAFT.txt,matrix.json,7-android*} docs/har-benchmark/
```

## CI gate

`go test ./internal/helpers/ -run TestHarBenchmarkAndroidNotClaimed` — phase 6 blocker SSOT must stay `not_claimed` until Android image ships `:7070`.
