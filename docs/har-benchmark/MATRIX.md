# HAR completeness matrix

**Plan:** [har-completeness-benchmark.plan.md](../../../../../.cursor/plans/har-completeness-benchmark.plan.md) — **done 2026-07-29** (phases 0–7; phase 6 Android not claimed).

**Committed SSOT:** `selenoid-tests/docs/har-benchmark/` (this file synced from last green run).

Fixture: `http://127.0.0.1:8080` (remote browser: `http://host.docker.internal:8080`)

Re-run: **2026-07-29** · local hub **HEAD** (fresh binary on `:4447` for this run; settle **4s** after load) · prod cutover hub **v3.0.5** / UI **v3.0.14** — hub bodies + **HarCapture 3b bodies** smoke on [selenoid.qa.guru](https://selenoid.qa.guru) · `qaguru/playwright-chromium:1.61.1` (`EXPOSE 7070`) · `qaguru/webdriver-chrome:149` (warm `:7070`) · selenoid-ui on `:8080` (fixture)

| Source | Writer | entries | http | urls | status | reqHdr | resHdr | size>0 | content.text | URL cov vs PW-local |
|--------|--------|--------:|-----:|-----:|-------:|-------:|-------:|-------:|-------------:|--------------------:|
| 1-pw-local-recordHar | Playwright | 26 | 26 | 25 | 26 | 26 | 26 | 25 | 25 | baseline |
| 2-pw-selenoid-recordHar | Playwright | 26 | 26 | 25 | 26 | 26 | 26 | 25 | 25 | 100% |
| 3-selenide-HarCapture | zero-design-system HarCapture (meta default) | 28 | 28 | 27 | 28 | 28 | 28 | 27 | 0 | 100% |
| 4-wd-hub-enableHAR | selenoid (meta default) | 29 | 29 | 27 | 7 | 29 | 7 | 0 | 0 | 100% |
| 5-pw-hub-enableHAR | selenoid (meta default) | 26 | 26 | 25 | 20 | 26 | 20 | 0 | 0 | 100% |

### Supplementary (same run, gates in compare tests)

| Source | Writer | entries | http | urls | status | reqHdr | resHdr | size>0 | content.text | URL cov vs PW-local |
|--------|--------|--------:|-----:|-----:|-------:|-------:|-------:|-------:|-------------:|--------------------:|
| 3b-selenide-HarCapture-bodies | HarCapture `BODIES` (no hub `enableHAR`) | 28 | 28 | 27 | 28 | 28 | 28 | 27 | 27 | 100% |
| 4b-wd-hub-enableHAR-bodies | selenoid (`harContent=bodies`, hub ≥ v3.0.5) | 29 | 29 | 27 | 29 | 29 | 29 | 28 | 28 | 100% |
| 5b-pw-hub-enableHAR-bodies | selenoid (`harContent=bodies`, hub ≥ v3.0.5) | 26 | 26 | 25 | 25 | 26 | 25 | 25 | 25 | 100% |

## Diff vs previous MATRIX (2026-07-29 earlier / prod cutover row)

| Row | Was (http / urls / text) | Now | Why |
|-----|--------------------------|-----|-----|
| 1 baseline | 9 / 9 / 9 | 26 / 25 / 25 | Fixture selenoid-ui on `:8080` serves current embedded UI (more CSS/assets); settle 4s |
| 2 PW→Selenoid | 9 / 9 / 9 | 26 / 25 / 25 | Same fixture + remote URL; parity with 1 unchanged (100% URL cov) |
| 3 HarCapture meta | 11 / 11 / 0 | 28 / 27 / 0 | Extra client-side requests on same page; `content.text=0` unchanged |
| 3b HarCapture bodies | 11 / 11 / 11 | 28 / 27 / 27 | Same fixture growth; bodies gate `withContentText>=1` still green |
| 4 WD hub meta | 11 / 11 / 0, status 6 | 29 / 27 / 0, status 7 | Larger fixture; hub meta still omits text/size; partial status at quit |
| 4b WD hub bodies | 11 / 11 / 11 | 29 / 27 / 28 | Hub ≥ v3.0.5 size-from-body; text+size gates green |
| 5 PW hub meta | 9 / 9 / 0, status 5 | 26 / 25 / 0, status 20 | Larger fixture; **requires hub binary with PW `enableHAR` + `:7070`** |
| 5b PW hub bodies | 9 / 9 / 9 | 26 / 25 / 25 | Hub ≥ v3.0.5; text+size gates green on expanded fixture |

**Infra note:** registry hub `:4444` (Jul 27 process) does not emit PW hub-HAR; green run used freshly built hub on `:4447`. Restart `:4444` after `./scripts/build-selenoid.sh` to align registry stand with test defaults.

## Prod Step 5b smoke (example.com, short)

| Source | Writer | har HTTP | http | withContentText | withContentSize | gate |
|--------|--------|---------:|-----:|----------------:|----------------:|------|
| wd-chrome-enableHAR-meta-prod | hub meta (omit harContent) | 200 | 1 | 0 | 0 | pass |
| wd-chrome-enableHAR-bodies-prod | hub `harContent=bodies` | 200 | 1 | 1 | 1 | pass |
| pw-chromium-enableHAR-meta-prod | hub meta | 200 | 1 | 0 | 0 | pass |
| pw-chromium-enableHAR-bodies-prod | hub `harContent=bodies` | 200 | 1 | 1 | 1 | pass |
| 3b-selenide-HarCapture-bodies-prod | zero-design-system HarCapture (`HarContentMode.BODIES`, no hub `enableHAR`) | — | 1 | 1 | 1 | pass |

UI: `harContent` control only when `enableHAR`; HarViewer Response tab shows `content.text` on bodies session; table/detail **Size** not «—» when `content.size>0` (WD bodies smoke: 559 B).

## Verdict

- **A (client):** PW-native ≈ PW→Selenoid (`recordHar`); 100% URL coverage on fixture
- **A′ (client bodies):** Selenide HarCapture `HarContentMode.BODIES` on same fixture (separate session, no hub `enableHAR`) — URL cov ≥ meta (this run 100%); best-effort gate `withContentText >= 1` (this run: 27). **Not** ≡ `recordHar`
- **A′ prod (client bodies):** HarCapture `HarContentMode.BODIES` on prod [selenoid.qa.guru](https://selenoid.qa.guru) warm WD Chrome — short example.com smoke: `withContentText >= 1` and `withContentSize >= 1` best-effort (this run 1/1). **Not** ≡ `recordHar`; one writer (no hub `enableHAR`)
- **B (hub meta):** hub `enableHAR` default (`harContent=meta` / omit) — URL cov 100%; gate `withContentText==0`; partial `status` / no `content.size` (known hub CDP gap vs `recordHar`)
- **B′ (hub bodies):** hub `enableHAR` + `harContent=bodies` on hub **≥ v3.0.5** — URL cov unchanged; fixture gates `withContentText >= 1` **and** `withContentSize >= 1` (this run: WD 28/28, PW 25/25). **Not** ≡ `recordHar`
- **5:** PW hub HAR meta green on hub ≥ HEAD with PW `:7070` mapping (one writer; no client `recordHar`)
- **Prod:** hub **v3.0.5** + UI **v3.0.14** on [selenoid.qa.guru](https://selenoid.qa.guru); Step 5b smoke green for WD/PW meta + bodies (`withContentText` **and** `withContentSize` gates on hub ≥ v3.0.5) + UI wire/viewer/size + HarCapture 3b prod — **bodies text+size claimed on prod (best-effort, not ≡ recordHar)**
- **One writer per session** — observed on all rows
- **Android:** not claimed (phase 6 verified 2026-07-29 — see below)

## Phase 6 — Android (not claimed)

| Check | Result |
|-------|--------|
| CDP `:7070` on `qaguru/android:*` | **no** — `EXPOSE 4444 5900` only; Appium/UiAutomator2 |
| Hub `harCaptureEnabled` | **blocked** — needs `devtoolsHostPort`; Android has none |
| UI `enableHAR` cap | **hidden** — `buildAndroidSelenoidOptions` omits HAR; toggle WebDriver-only |
| Chrome Mobile in image | **deferred** — `browser-image/android/README.md` |
| Scorecard / `.har` | **not run** — no path to capture |
| Verdict | **not_claimed** — blocker doc: `7-android-blocker.json` |

Roadmap (not this phase): Chrome Mobile + chromedriver + `:7070` proxy in android image, then UI cap + fixture scorecard. MITM sidecar explicitly out of scope per `har.adoc`.

## Allowed claims

- Playwright→Selenoid client `recordHar` ≈ native Playwright local (100% URL coverage on fixture)
- Selenide HarCapture **meta default** on Selenoid covers baseline URLs (100%); `content.text=0` on this fixture
- Selenide HarCapture **bodies** on Selenoid (local fixture): `withContentText >= 1` best-effort (this run 27/27); URL cov not weaker than meta; one writer (no hub `enableHAR` on that session)
- Selenide HarCapture **bodies** on prod [selenoid.qa.guru](https://selenoid.qa.guru) (warm WD Chrome, example.com smoke): `withContentText >= 1` best-effort (this run 1/1); `withContentSize >= 1` best-effort; one writer; **not** ≡ `recordHar`
- Hub `enableHAR` default meta (WD warm Chrome + PW Chromium-family `:7070`): URL coverage ≥80% of PW baseline; `withContentText==0`
- Hub `enableHAR` + `harContent=bodies` on **local hub ≥ v3.0.5**: `withContentText >= 1` **and** `withContentSize >= 1` on fixture for part or all of http entries (best-effort; this run WD 28/28, PW 25/25); URL cov gate unchanged
- Hub `enableHAR` + `harContent=bodies` on prod [selenoid.qa.guru](https://selenoid.qa.guru) (hub ≥ **v3.0.5**, UI ≥ **v3.0.14**, images with `:7070`): WD warm Chrome + PW Chromium-family — `withContentText >= 1` **and** `withContentSize >= 1` on short example.com smoke (best-effort); meta omit still `withContentText==0`
- One writer per session
- Prod [selenoid.qa.guru](https://selenoid.qa.guru) hub HAR for warm WD Chrome **and** Playwright Chromium-family — meta/URL path **and** opt-in bodies text+size (as above)

## Do not claim

- Hub HAR (meta or bodies) identical to Playwright `recordHar` (status/size/text completeness differs)
- Hub `bodies` ≡ `recordHar` text/status/size counts
- Client HarCapture `bodies` ≡ `recordHar` (status/size/text parity not promised; gate is best-effort ≥1)
- Playwright hub `enableHAR` on firefox / webkit / `*-min` (no `:7070` by design)
- Android HAR
- Dual-writer «for completeness» (hub `enableHAR` + client `recordHar` / `HarCapture` on one session)

## Known gaps (hub CDP vs recordHar)

- default **meta** omits `content.text`; bodies opt-in best-effort (not absolute “never text”)
- hub **meta** path may still have `content.size=0` and partial `status` (in-flight at quit)
- hub **bodies** on ≥ v3.0.5 sets `content.size` from decoded body length when `getResponseBody` succeeds; still not ≡ `recordHar`
- client HarCapture meta: `content.text=0` on fixture; bodies opt-in separate session
- do not combine hub `enableHAR` with client `recordHar` / `HarCapture` on one session
- prod smoke uses short example.com; local fixture uses selenoid-ui — counts differ; gates only

## Artifacts

- `summary.json` · `hub-summary.json`
- `1-playwright-local.har` · `2-playwright-selenoid.har` · `3-selenide-selenoid.har`
- `4-wd-hub-enableHAR.har` · `5-pw-hub-enableHAR.har`
- supplementary: `3b-selenide-HarCapture-bodies.har` · `4b-wd-hub-enableHAR-bodies.har` · `5b-pw-hub-enableHAR-bodies.har`
- prod_step5_dir: `prod-step5/` (`summary.json` · `ui-smoke-summary.json` · `harcapture-prod-summary.json`)
- prod hub: `prod-step5/wd-chrome-enableHAR-meta-prod.har` · `prod-step5/wd-chrome-enableHAR-bodies-prod.har` · `prod-step5/pw-chromium-enableHAR-meta-prod.har` · `prod-step5/pw-chromium-enableHAR-bodies-prod.har`
- prod HarCapture: `3b-selenide-HarCapture-bodies-prod.har` · `prod-step5/3b-selenide-HarCapture-bodies-prod.har`
- prod UI: `ui-caps-harContent-hidden-prod.png` · `ui-caps-harContent-visible-prod.png` · `ui-harViewer-bodies-prod.png` (also under `prod-step5/`)
- also mirrored under `build/har-prod-e2e/`
- manual: `5-manual-wd-har-viewer.png` · `5-manual-wd-session-live.png`
- android: `7-android.NOT_CLAIMED.txt` · `7-android-blocker.json`
