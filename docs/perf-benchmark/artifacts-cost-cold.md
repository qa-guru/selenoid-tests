# Artifacts cost — Java cold one-shot (Benchmarks §4)

**Host:** box1 `java-jdk21-jenkins-warm-agent-1` · workspace `java-cold-pool` · login-test.  
**Method:** sequential `./gradlew test` with attach/capability flips. Not a Jenkins job; no permanent sampler.  
**Measured:** 2026-08-21 ~00:05 UTC.  
**IDs:** `art-java-wd-warm-1-p1-*` in [`runs.json`](runs.json).

| preset | wall_s | video_kb | log_kb | har_kb | total |
|--------|-------:|---------:|-------:|-------:|------:|
| none | 2.729 | 0 | 0 | 0 | 0 |
| log | 3.939 | 0 | 73 | 0 | 73 |
| video | 5.259 | 44 | 0 | 0 | 44 |
| har meta | 3.995 | 0 | 0 | 49 | 49 |
| har bodies | 3.995 | 0 | 0 | 49 | 49 |
| video+log | 5.163 | 44 | 73 | 0 | 117 |
| video+log+har meta | 5.296 | 44 | 73 | 49 | 166 |
| video+log+har bodies | 5.296 | 44 | 73 | 49 | 166 |

`har bodies` mirrors `har meta`: Java `HarCapture` is CDP performance-log meta (no `Network.getResponseBody` yet).

Jenkins cold none/full pins stay in §0; §4 reads only `art-java-*`.
