# Warm isolation — one option at a time (one-shot)

**Host:** box1 warm agent `java-jdk21-jenkins-warm-agent-1` · workspace `java-warm-pool` · login-test only.  
**Method:** sequential `./gradlew test` with a single `-D` attach/capability flip. Not a Jenkins job; no permanent sampler.  
**Measured:** 2026-08-20 ~23:47 UTC.

Common baseline flags: `remoteUrl=http://host.docker.internal:4444/wd/hub`, `headless=false`, warm reuse, VNC/video/HAR off unless the variant enables them.

## Walls

| variant | wall_ms | Δ vs none | Δ vs allure3-empty | att | weight |
|---------|--------:|----------:|-------------------:|----:|--------|
| **none** (`allureReportMode=none`) | 2558 | 0 | — | 0 | 0 |
| **allure3-empty** (allure3, no attaches) | 3745 | +1187 | 0 | 0 | 0 |
| **screenshot** only | 3978 | +1420 | **+233** | 1 | **52 KB** |
| **pagesource** only | 3830 | +1272 | **+85** | 1 | **21 KB** |
| **console** only | 3858 | +1300 | **+113** | 1 | **2 KB** |
| **video** only (`enableVideo`+`attachVideo`) | 4973 | +2415 | **+1228** | 1 stub | **44 KB** mp4 |
| **har** only (`enableHar`+`attachHarLogs`) | 3951 | +1393 | **+206** | 2* | **49 KB** `.har` |

\* HAR row also writes a viewer HTML in results; weight column is **raw `.har` only** (viewer excluded).

Video file: `/opt/selenoid/video/7dd226c825dfeb96393549fe772c95a9.mp4` → `du -k` = **44**.

## Takeaways

1. **Allure3 tax** (~+1.2s vs `none`) dominates lite attach walls — compare options to `allure3-empty`, not to `none`.
2. **Weight (isolated):** screenshot 52 · HAR 49 · video 44 · page source 21 · console 2 (KB).
3. **Wall (isolated, vs allure3-empty):** video is the expensive one (~+1.2s); screenshot/HAR ~+0.2s; page source/console ~+0.1s (noise band — single run).
4. Lite pack (screenshot+source+console) ≈ 52+21+2 = **75 KB** — matches earlier warm-lite pin ~73 KB.
5. Do **not** leave this runner in Jenkins XML; re-run only when re-pinning.

## Δ tax summary (vs allure3-empty)

```
console     +0.1s   2 KB
pagesource  +0.1s  21 KB
screenshot  +0.2s  52 KB
har         +0.2s  49 KB
video       +1.2s  44 KB
allure3     +1.2s   0 KB   (vs none)
```
