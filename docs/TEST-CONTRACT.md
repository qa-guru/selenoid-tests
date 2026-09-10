# Session lifecycle test contract

Coverage of **lines** is not the gate. A test that skips, waits three minutes, or `rerender`s as if SSE already caught up does not cover Stop/Delete. A test that creates the session via hub API does not cover **Create Session**.

## Stop / Delete — what must fail

| Layer | Contract |
|-------|----------|
| RTL (`selenoid-ui/ui`) | Click **Stop session** while the parent still passes a live `browser`. VNC unmounts, FINISHED shows, Delete stays disabled until an artifact exists, then enables. Failed hub DELETE restores Stop. Already-closed RFB is not `disconnect()`ed again. |
| Hub unit (`selenoid`) | Playwright `enableHAR` DELETE removes the session without waiting on CDP (pending HAR < 2s in unit). |
| API | `DELETE /wd/hub/session/{id}` returns within **5s**. VNC/video/HAR sessions use desktop `chromeVersion`, not profile `browserVersion` when that is `*-min` (hub v3.0.16 rejects `-min` + those flags with 400). |
| Integration WD/PW/min | Advertised catalog versions **create**. Create failure on `qa_guru` / `github` is **fail**, not skip. PW `enableHAR` hub DELETE within **8s** (including `*-min` without DevTools `:7070`). |
| E2E smoke | Create (API) → session page → Stop → FINISHED and VNC gone in **2s** (do not wait for `/events`) → video → **Delete session** leaves the page. No skip. No `session-kill` testid. |
| E2E hub-prod | HAR layout after Stop. May skip **github CI** archive timing. Must run on `qa_guru`. Must not be tagged `smoke` if it skips when `TEST_TAGS=smoke`. |

## New Session — what must fail

| Layer | Contract |
|-------|----------|
| Unit | `buildSelenoidOptions` / Playwright WS query: string `"false"` is **false** (`Boolean("false")` is not). VNC/video/HAR/log names and `harContent=bodies` only when the flag is on. |
| RTL WD | Click **Create Session** with default VNC+video. POST `/wd/hub/session` has `enableVNC`/`enableVideo` true and `enableHAR` false. HTTP 200 navigates to `/sessions/:id`. HTTP 500 / `Failed to fetch` **stays** on New Session, Create re-enabled, `error-true` title is the hub message. Toggle VNC off → POST `enableVNC: false`. |
| RTL PW | Click Create. WS URL has `enableVNC`/`enableVideo`. HAR toggle → `enableHAR` on the query. |
| RTL Android | Click Create. `appium:*` alwaysMatch. Failed POST stays on New Session. |
| E2E smoke | Open `/#/new-session` → select chrome → click **Create Session** → `/#/sessions/{id}` with `session-stop`. No skip. Fail as soon as the button is `error-true` — do not wait minutes for a URL change. |

## Skip is allowed only when the fixture is absent

- Warm-pool stand down / `warmTotal=0`
- Android not in hub catalog
- `msedge-min` on arm64 (no image manifest)

Forbidden: skip because create failed for an advertised `*-min`; skip Stop/Delete or Create Session to keep smoke green; wait minutes for FINISHED or for Create after `error-true`; the only “after kill” RTL test is `rerender({ browser: undefined })`; the only Create e2e is hub API + Goto session page.

## Testids (UI ↔ e2e SSOT)

Stop / Delete: `session-stop` · `session-delete` · `session-close` · `session-finished` · `vnc-window` · `session-detail-video` · `session-video-waiting` · `session-artifacts-not-found`

New Session: `capabilities-setup` · `capabilities-create-session` · `capabilities-browser-select` · `capabilities-browser-select-playwright` · `caps-enable-vnc` · `caps-enable-video` · `caps-enable-har` · `caps-enable-log` · `caps-playwright-enable-vnc` · `caps-playwright-enable-video` · `caps-playwright-enable-har`

Rename in UI → update `tests/e2e/ui/capabilities_helpers_test.go` in the same change.

## Coverage gate

Vitest thresholds **100% lines** on owner files of these contracts:

- `src/components/Session/index.tsx`
- `src/components/Session/SessionInfo.tsx`
- `src/components/Sessions/service.ts`
- `src/util/playwrightSessions.ts`
- `src/components/VncCard/VncScreen.tsx`
- `src/util/capabilitiesLogic.ts`
- `src/util/capabilitiesPlaywright.ts`
- `src/util/waitForLiveSession.ts`
- `src/components/CapabilitiesLaunchActions/index.tsx`

Not the whole UI (Capabilities god file, Benchmarks, Docs). Those stay other contracts — do not put a 100% gate on `Capabilities/index.tsx`.
