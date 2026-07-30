# HAR marketing pack — assets list

**SSOT matrix:** [MATRIX.md](MATRIX.md) · **draft post:** [TELEGRAM-DRAFT.txt](TELEGRAM-DRAFT.txt)  
**Run:** 2026-07-29 · prod hub **v3.0.5** · UI **v3.0.14** · local hub HEAD on `:4447`

## Primary visuals (Telegram / landing — 3 штуки)

| # | Role | File | Caption hint |
|---|------|------|--------------|
| 1 | **Client PW parity** — local vs Selenoid `recordHar` | *(text/table)* — rows 1–2 из MATRIX; опционально collage двух HarViewer/Allure если есть | «Playwright local ≈ Playwright→Selenoid — 100% URL cov, 26/25» |
| 2 | **Client Selenide bodies** — Allure/client path без hub | `3b-selenide-HarCapture-bodies.har` + prod `prod-step5/3b-selenide-HarCapture-bodies-prod.har` | «HarCapture BODIES — one writer, no hub enableHAR; fixture 27/27 text» |
| 3 | **UI HarViewer full-width** — hub bodies на prod | `ui-harViewer-bodies-prod.png` | «Живой network archive рядом с VNC; Response + Size на bodies» |

## Secondary (thread / landing detail)

| Role | File | Notes |
|------|------|-------|
| harContent toggle hidden | `ui-caps-harContent-hidden-prod.png` | enableHAR off → harContent не показывается |
| harContent toggle visible | `ui-caps-harContent-visible-prod.png` | enableHAR on → opt-in bodies |
| Response tab detail | `prod-step5/ui-harViewer-bodies-response-prod.png` | content.text в HarViewer |
| WD manual session (legacy) | `5-manual-wd-session-live.png` | VNC + session page context |
| WD manual HarViewer (legacy) | `5-manual-wd-har-viewer.png` | pre-cutover UX reference |

## Scorecard artifacts (rows 1–5 + 3b/4b/5b)

| Row | HAR | MATRIX metrics (http / urls / text) |
|-----|-----|-------------------------------------|
| 1 | `1-playwright-local.har` | 26 / 25 / 25 |
| 2 | `2-playwright-selenoid.har` | 26 / 25 / 25 |
| 3 | `3-selenide-selenoid.har` | 28 / 27 / 0 |
| 3b | `3b-selenide-HarCapture-bodies.har` | 28 / 27 / 27 |
| 4 | `4-wd-hub-enableHAR.har` | 29 / 27 / 0 |
| 4b | `4b-wd-hub-enableHAR-bodies.har` | 29 / 27 / 28 |
| 5 | `5-pw-hub-enableHAR.har` | 26 / 25 / 0 |
| 5b | `5b-pw-hub-enableHAR-bodies.har` | 26 / 25 / 25 |

## Prod Step 5b smoke (`prod-step5/`)

| Source | HAR | text | size | gate |
|--------|-----|-----:|-----:|------|
| WD meta | `wd-chrome-enableHAR-meta-prod.har` | 0 | 0 | pass |
| WD bodies | `wd-chrome-enableHAR-bodies-prod.har` | 1 | 1 | pass |
| PW meta | `pw-chromium-enableHAR-meta-prod.har` | 0 | 0 | pass |
| PW bodies | `pw-chromium-enableHAR-bodies-prod.har` | 1 | 1 | pass |
| HarCapture 3b | `3b-selenide-HarCapture-bodies-prod.har` | 1 | 1 | pass |

Summaries: `summary.json` · `hub-summary.json` · `prod-step5/summary.json` · `prod-step5/ui-smoke-summary.json` · `harcapture-prod-summary.json`

## Mirrors

- `build/har-prod-e2e/` — prod UI PNG + named HAR (same run)

## Cancelled (phase 6 — Android HAR)

- `7-android-blocker.json` — status **cancelled**, not on roadmap
- Do not claim Android HAR

## Allowed claims (copy-paste guard)

From MATRIX.md § Allowed claims — use verbatim boundaries:

- PW→Selenoid client `recordHar` ≈ local PW (100% URL cov on fixture)
- HarCapture meta: URL cov 100%; `content.text=0` on fixture
- HarCapture bodies: fixture ≥1 text (27/27); prod example.com 1/1 text+size; **not** ≡ recordHar
- Hub meta: URL cov ≥80%; `withContentText==0`
- Hub bodies (≥ v3.0.5): fixture WD 28/28, PW 25/25 text+size; prod smoke green; **not** ≡ recordHar
- One writer per session
- Prod hub HAR for warm WD Chrome **and** PW Chromium — meta + opt-in bodies

## Do not claim

See MATRIX.md § Do not claim (hub ≡ recordHar, Android, dual-writer, firefox/webkit/min PW hub-HAR).
