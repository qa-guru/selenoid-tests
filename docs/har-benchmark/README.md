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

Hub enableHAR meta/bodies + HarCapture bodies are Go `har-prod` (writes `prod-step5/*.har`). UI harContent toggle + HarViewer bodies: `tests/e2e/ui/ui_har_content_test.go` (hub-prod / ui slice).

```bash
PYRAMID_STAND=selenoid_qa_guru ./scripts/run-go-pyramid.sh har-prod
# optional UI: PYRAMID_STAND=selenoid_qa_guru ./scripts/run-go-pyramid.sh ui
```

Sync if matrix changed:

```bash
cp build/har-compare/{MATRIX.md,ASSETS.md,TELEGRAM-DRAFT.txt,matrix.json} docs/har-benchmark/
```

## CI gate

`go test ./internal/helpers/ -run TestHarBenchmarkAndroidCancelled` — phase 6 must stay **cancelled**.
