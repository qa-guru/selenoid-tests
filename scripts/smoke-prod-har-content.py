#!/usr/bin/env python
"""Prod smoke: hub enableHAR meta + bodies (WD Chrome + PW Chromium). No dual-writer.

Writes artifacts under build/har-compare/prod-step5/ and prints a JSON summary.
Does not touch android / firefox / webkit / *-min.
"""
from __future__ import annotations

import json
import os
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

BASE = os.environ.get("SELENOID_PUBLIC_URL", "https://selenoid.qa.guru").rstrip("/")
USER = os.environ.get("SELENOID_USER", "qa_engineer")
PASSWORD = os.environ.get("SELENOID_PASSWORD", "aAb_-4gs53FD")
ACCESS_KEY = os.environ.get("SELENOID_ACCESS_KEY", f"{USER}:{PASSWORD}")
OUT = Path(
    os.environ.get(
        "HAR_SMOKE_OUT",
        str(Path(__file__).resolve().parents[1] / "build" / "har-compare" / "prod-step5"),
    )
)
WD_BROWSER = os.environ.get("SMOKE_WD_BROWSER", "chrome")
WD_VERSION = os.environ.get("SMOKE_WD_VERSION", "149.0")
PW_BROWSER = os.environ.get("SMOKE_PW_BROWSER", "playwright-chromium")
PW_VERSION = os.environ.get("SMOKE_PW_VERSION", "1.61.1")
NAV_URL = os.environ.get("SMOKE_NAV_URL", "https://example.com/")


def _ssl_context() -> ssl.SSLContext:
    try:
        import certifi

        return ssl.create_default_context(cafile=certifi.where())
    except Exception:
        return ssl.create_default_context()


SSL_CTX = _ssl_context()


def _basic_auth_header() -> str:
    import base64

    token = base64.b64encode(f"{USER}:{PASSWORD}".encode()).decode()
    return f"Basic {token}"


def http_json(method: str, url: str, body: dict | None = None, auth: bool = True) -> tuple[int, object]:
    data = None if body is None else json.dumps(body).encode()
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    if auth:
        headers["Authorization"] = _basic_auth_header()
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=120, context=SSL_CTX) as resp:
            raw = resp.read()
            code = resp.getcode()
            if not raw:
                return code, None
            return code, json.loads(raw.decode())
    except urllib.error.HTTPError as e:
        raw = e.read()
        try:
            payload = json.loads(raw.decode()) if raw else None
        except Exception:
            payload = raw.decode(errors="replace") if raw else None
        return e.code, payload


def har_stats(har: dict) -> dict:
    entries = (((har or {}).get("log") or {}).get("entries")) or []
    http = 0
    with_text = 0
    with_size = 0
    for e in entries:
        req = e.get("request") or {}
        url = req.get("url") or ""
        if not url.startswith("http"):
            continue
        http += 1
        content = ((e.get("response") or {}).get("content")) or {}
        text = content.get("text")
        if isinstance(text, str) and text != "":
            with_text += 1
        size = content.get("size")
        if isinstance(size, (int, float)) and size > 0:
            with_size += 1
    return {
        "entries": len(entries),
        "http": http,
        "withContentText": with_text,
        "withContentSize": with_size,
    }


def fetch_har(session_id: str) -> tuple[int, dict | None]:
    code, payload = http_json("GET", f"{BASE}/har/{session_id}.har")
    if code == 200 and isinstance(payload, dict):
        return code, payload
    # some proxies may return text; retry raw
    req = urllib.request.Request(
        f"{BASE}/har/{session_id}.har",
        headers={"Authorization": _basic_auth_header(), "Accept": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=60, context=SSL_CTX) as resp:
            return resp.getcode(), json.loads(resp.read().decode())
    except Exception as e:
        return 0, {"error": str(e), "payload": payload}


def wd_session(har_content: str | None) -> dict:
    options: dict = {
        "enableVNC": True,
        "enableHAR": True,
        "name": f"step5-wd-{'bodies' if har_content == 'bodies' else 'meta'}",
    }
    if har_content:
        options["harContent"] = har_content
    caps = {
        "browserName": WD_BROWSER,
        "browserVersion": WD_VERSION,
        "selenoid:options": options,
    }
    code, body = http_json(
        "POST",
        f"{BASE}/wd/hub/session",
        {"capabilities": {"alwaysMatch": caps}},
    )
    if code >= 400 or not isinstance(body, dict):
        raise RuntimeError(f"WD create failed {code}: {body}")
    value = body.get("value") or {}
    sid = value.get("sessionId") or (body.get("sessionId"))
    if not sid:
        raise RuntimeError(f"WD no sessionId: {body}")
    # navigate
    nav_code, _ = http_json(
        "POST",
        f"{BASE}/wd/hub/session/{sid}/url",
        {"url": NAV_URL},
    )
    if nav_code >= 400:
        http_json("DELETE", f"{BASE}/wd/hub/session/{sid}")
        raise RuntimeError(f"WD navigate failed {nav_code}")
    time.sleep(2.5)
    http_json("DELETE", f"{BASE}/wd/hub/session/{sid}")
    # HAR appears after session delete
    deadline = time.time() + 45
    har_code, har = 0, None
    while time.time() < deadline:
        har_code, har = fetch_har(sid)
        if har_code == 200 and isinstance(har, dict) and (har.get("log") or {}).get("entries") is not None:
            break
        time.sleep(1.0)
    stats = har_stats(har if isinstance(har, dict) else {})
    label = "bodies" if har_content == "bodies" else "meta"
    path = OUT / f"wd-chrome-enableHAR-{label}-prod.har"
    if isinstance(har, dict):
        path.write_text(json.dumps(har, indent=2), encoding="utf-8")
    return {
        "writer": "wd",
        "mode": label,
        "sessionId": sid,
        "harHttp": har_code,
        "stats": stats,
        "artifact": str(path) if path.exists() else None,
        "ok_meta": (
            label == "meta"
            and har_code == 200
            and stats["http"] >= 1
            and stats["withContentText"] == 0
        ),
        "ok_bodies": (
            label == "bodies"
            and har_code == 200
            and stats["http"] >= 1
            and stats["withContentText"] >= 1
            and stats["withContentSize"] >= 1
        ),
    }


def collect_session_ids(status: dict | None) -> set[str]:
    ids: set[str] = set()

    def walk(node):
        if isinstance(node, dict):
            sessions = node.get("sessions")
            if isinstance(sessions, list):
                for sess in sessions:
                    if isinstance(sess, dict):
                        sid = sess.get("id")
                        if isinstance(sid, str) and sid:
                            ids.add(sid)
            for v in node.values():
                walk(v)
        elif isinstance(node, list):
            for v in node:
                walk(v)

    walk(status)
    return ids


def wait_new_session_id(before: set[str], timeout_s: float = 30.0) -> str | None:
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        _, status = http_json("GET", f"{BASE}/hub/status", auth=False)
        now = collect_session_ids(status if isinstance(status, dict) else {})
        delta = now - before
        if delta:
            return next(iter(delta))
        time.sleep(0.25)
    return None


def pw_session(har_content: str | None) -> dict:
    """Minimal Playwright-over-WS smoke via hub CDP is heavy; use browser_info + enableHAR through
    a short Node/Playwright connect if available, else HTTP-less WS handshake is not enough.

    Prefer: selenium-like path is WD-only. For PW we use `playwright` python sync if installed,
    else mark skipped with reason.
    """
    label = "bodies" if har_content == "bodies" else "meta"
    try:
        from playwright.sync_api import sync_playwright  # type: ignore
    except Exception as e:
        return {
            "writer": "pw",
            "mode": label,
            "skipped": True,
            "reason": f"playwright python not installed: {e}",
            "ok_meta": False,
            "ok_bodies": False,
        }

    q = {
        "accessKey": ACCESS_KEY,
        "enableVNC": "true",
        "enableHAR": "true",
        "name": f"step5-pw-{label}",
    }
    if har_content:
        q["harContent"] = har_content
    qs = urllib.parse.urlencode(q)
    ws = f"{BASE.replace('https://', 'wss://').replace('http://', 'ws://')}/playwright/{PW_BROWSER}/{PW_VERSION}?{qs}"

    want_name = f"step5-pw-{label}"
    _, status_before = http_json("GET", f"{BASE}/hub/status", auth=False)
    before_ids = collect_session_ids(status_before if isinstance(status_before, dict) else {})
    sid = None
    with sync_playwright() as p:
        browser = p.chromium.connect(ws, timeout=120_000)
        sid = wait_new_session_id(before_ids, timeout_s=30.0)
        page = browser.new_page()
        # Hub attaches CDP to /page asynchronously; yield before first navigation.
        page.wait_for_timeout(750)
        page.goto(NAV_URL, wait_until="load", timeout=60_000)
        page.wait_for_timeout(1500)
        try:
            browser.close()
        except Exception:
            pass

    if not sid:
        return {
            "writer": "pw",
            "mode": label,
            "skipped": False,
            "error": "could not resolve new Playwright session id from /hub/status",
            "ok_meta": False,
            "ok_bodies": False,
        }

    deadline = time.time() + 45
    har_code, har = 0, None
    while time.time() < deadline:
        har_code, har = fetch_har(sid)
        if har_code == 200 and isinstance(har, dict) and (har.get("log") or {}).get("entries") is not None:
            break
        time.sleep(1.0)
    stats = har_stats(har if isinstance(har, dict) else {})
    path = OUT / f"pw-chromium-enableHAR-{label}-prod.har"
    if isinstance(har, dict):
        path.write_text(json.dumps(har, indent=2), encoding="utf-8")
    return {
        "writer": "pw",
        "mode": label,
        "sessionId": sid,
        "harHttp": har_code,
        "stats": stats,
        "artifact": str(path) if path.exists() else None,
        "ok_meta": (
            label == "meta"
            and har_code == 200
            and stats["http"] >= 1
            and stats["withContentText"] == 0
        ),
        "ok_bodies": (
            label == "bodies"
            and har_code == 200
            and stats["http"] >= 1
            and stats["withContentText"] >= 1
            and stats["withContentSize"] >= 1
        ),
    }


def main() -> int:
    OUT.mkdir(parents=True, exist_ok=True)
    results = []
    # sequential — one writer, used==0 between runs
    for mode in (None, "bodies"):
        print(f"=== WD {WD_BROWSER}/{WD_VERSION} harContent={mode or 'meta(default)'} ===", flush=True)
        r = wd_session(mode)
        results.append(r)
        print(json.dumps(r, indent=2), flush=True)
        time.sleep(2)
        print(f"=== PW {PW_BROWSER}/{PW_VERSION} harContent={mode or 'meta(default)'} ===", flush=True)
        r = pw_session(mode)
        results.append(r)
        print(json.dumps(r, indent=2), flush=True)
        time.sleep(2)

    summary = {
        "base": BASE,
        "nav": NAV_URL,
        "hubVersionNote": "prod size gate requires hub ≥ v3.0.5 (content.size from decoded body)",
        "bodiesMinWithContentText": 1,
        "bodiesMinWithContentSize": 1,
        "results": results,
        "gates": {
            "wd_meta": next((x.get("ok_meta") for x in results if x["writer"] == "wd" and x["mode"] == "meta"), False),
            "wd_bodies": next((x.get("ok_bodies") for x in results if x["writer"] == "wd" and x["mode"] == "bodies"), False),
            "pw_meta": next((x.get("ok_meta") for x in results if x["writer"] == "pw" and x["mode"] == "meta"), False),
            "pw_bodies": next((x.get("ok_bodies") for x in results if x["writer"] == "pw" and x["mode"] == "bodies"), False),
        },
    }
    (OUT / "summary.json").write_text(json.dumps(summary, indent=2), encoding="utf-8")
    print("=== SUMMARY ===", flush=True)
    print(json.dumps(summary["gates"], indent=2), flush=True)
    ok = all(summary["gates"].values())
    return 0 if ok else 2


if __name__ == "__main__":
    sys.exit(main())
