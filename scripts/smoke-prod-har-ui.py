#!/usr/bin/env python
"""Prod UI smoke: harContent control visibility + HarViewer body text on bodies HAR."""
from __future__ import annotations

import json
import os
import sys
from pathlib import Path

BASE = os.environ.get("SELENOID_PUBLIC_URL", "https://selenoid.qa.guru").rstrip("/")
OUT = Path(
    os.environ.get(
        "HAR_SMOKE_OUT",
        str(Path(__file__).resolve().parents[1] / "build" / "har-compare" / "prod-step5"),
    )
)
BODIES_SID = os.environ.get(
    "SMOKE_BODIES_SESSION_ID",
    "",
)


def main() -> int:
    from playwright.sync_api import sync_playwright

    OUT.mkdir(parents=True, exist_ok=True)
    summary_path = OUT / "summary.json"
    bodies_sid = BODIES_SID
    if not bodies_sid and summary_path.exists():
        data = json.loads(summary_path.read_text(encoding="utf-8"))
        for r in data.get("results") or []:
            if r.get("writer") == "wd" and r.get("mode") == "bodies" and r.get("sessionId"):
                bodies_sid = r["sessionId"]
                break
    if not bodies_sid:
        print("FAIL: no bodies session id (set SMOKE_BODIES_SESSION_ID or run HAR smoke first)", flush=True)
        return 2

    result = {
        "harContent_hidden_when_enableHAR_off": False,
        "harContent_visible_when_enableHAR_on": False,
        "harViewer_shows_body_text": False,
        "bodiesSessionId": bodies_sid,
        "artifacts": {},
    }

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(viewport={"width": 1440, "height": 1100})

        page.goto(f"{BASE}/#/new-session", wait_until="load", timeout=60_000)
        page.wait_for_selector('[data-testid="capabilities-browser-select"]', timeout=45_000)
        # Caps stack unlocks only after a browser is selected
        page.locator(
            '[data-testid="capabilities-browser-select"] button[aria-pressed="false"]',
            has_text="chrome: 149.0",
        ).filter(has_not_text="min").first.click(timeout=15_000)
        page.wait_for_selector('[data-testid="caps-enable-har"]', timeout=20_000)

        # default enableHAR off → harContent control absent
        if page.locator('[data-testid="caps-har-content"]').count() == 0:
            result["harContent_hidden_when_enableHAR_off"] = True
        shot_off = OUT / "ui-caps-harContent-hidden-prod.png"
        page.locator('[data-testid="capabilities-caps-flags"]').scroll_into_view_if_needed()
        page.screenshot(path=str(shot_off), full_page=False)
        result["artifacts"]["caps_off"] = str(shot_off)

        # turn enableHAR on — PlaqueFieldSeg buttons use data-value + aria-pressed
        page.locator('[data-testid="caps-enable-har"] button[data-value="true"]').click(timeout=10_000)
        page.wait_for_selector('[data-testid="caps-har-content"]', timeout=10_000)
        result["harContent_visible_when_enableHAR_on"] = page.locator('[data-testid="caps-har-content"]').count() > 0
        shot_on = OUT / "ui-caps-harContent-visible-prod.png"
        page.locator('[data-testid="caps-har-content"]').scroll_into_view_if_needed()
        page.screenshot(path=str(shot_on), full_page=False)
        result["artifacts"]["caps_on"] = str(shot_on)

        # HarViewer on finished bodies session — Response tab must show content.text
        page.goto(f"{BASE}/#/sessions/{bodies_sid}", wait_until="load", timeout=60_000)
        page.wait_for_selector('[data-testid="session-har-row-0"]', timeout=60_000)
        page.locator('[data-testid="session-har-row-0"]').click(timeout=10_000)
        page.locator('[data-testid="session-har-tab-response"]').click(timeout=10_000)
        panel = page.locator('[data-testid="session-har-panel-response"]')
        panel.wait_for(timeout=10_000)
        body = panel.inner_text(timeout=10_000)
        has_example = "Example Domain" in body or "<!doctype html>" in body.lower()
        muted = "Body not captured" in body
        result["harViewer_shows_body_text"] = bool(has_example and not muted)
        result["harViewer_snippet"] = body[:400]
        shot_hv = OUT / "ui-harViewer-bodies-prod.png"
        page.locator('[data-testid="session-har-viewer"]').scroll_into_view_if_needed()
        page.screenshot(path=str(shot_hv), full_page=False)
        result["artifacts"]["har_viewer"] = str(shot_hv)

        browser.close()

    (OUT / "ui-smoke-summary.json").write_text(json.dumps(result, indent=2), encoding="utf-8")
    print(json.dumps(result, indent=2), flush=True)
    ok = (
        result["harContent_hidden_when_enableHAR_off"]
        and result["harContent_visible_when_enableHAR_on"]
        and result["harViewer_shows_body_text"]
    )
    return 0 if ok else 2


if __name__ == "__main__":
    sys.exit(main())
