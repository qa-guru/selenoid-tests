# Selenoid performance benchmarks — SSOT

Published runs for the Selenoid UI **Benchmarks** tab (`#/benchmarks`). UI only displays this data — no live runner.

| File | Role |
|------|------|
| [schema.json](schema.json) | JSON Schema for `runs.json` |
| [runs.json](runs.json) | Measured / pending runs (SSOT) |

HAR **completeness** (field coverage) stays in [`../har-benchmark/`](../har-benchmark/) — do not mix with wall time / artifact KB / CPU·RAM.

## Dimensions

`language` · `protocol` · `image_flavor` (warm|min) · `pool` (cold|warm-pool|hot-pool) · `suite_size` (1|few|many → 1/10/100) · `parallel` · `artifacts` (video/log/har off|meta|bodies) · `versions.hub|ui|cm`

## Required metrics

- `wall_time_s`, session create p50/p95 (ms)
- `cpu_avg_pct` / `cpu_peak_pct`, `ram_avg_mb` / `ram_peak_mb`
- Artifact weights in **KB**: `video_kb`, `log_kb`, `har_kb`, `artifacts_total_kb` (0 when off)
- `host` profile on every run (CPU/RAM otherwise incomparable)

## Jenkins login-test

Four Java Selenide login jobs (same test as [warm-pool](https://jenkins.qa.guru/job/autotests-ai-multistack-tests-pipeline-java-warm-pool/)) live in `runs.json` with `jenkins_url`. UI tab `#/benchmarks` section **0. Jenkins login-test**.

| id | pool | Jenkins |
|----|------|---------|
| `jenkins-java-wd-cold-1-p1-headless-none` | cold | [cold-pool](https://jenkins.qa.guru/job/autotests-ai-multistack-tests-pipeline-java-cold-pool/) |
| `jenkins-java-wd-cold-1-p1-full-attachments` | cold | [cold-pool-full-attachments](https://jenkins.qa.guru/job/autotests-ai-multistack-tests-pipeline-java-cold-pool-full-attachments/) |
| `jenkins-java-wd-warm-1-p1-none` | warm-pool | [warm-pool](https://jenkins.qa.guru/job/autotests-ai-multistack-tests-pipeline-java-warm-pool/) |
| `jenkins-java-wd-hot-1-p1-none` | hot-pool | [hot-pool](https://jenkins.qa.guru/job/autotests-ai-multistack-tests-pipeline-java-hot-pool/) (stub) |

Optional field `jenkins_url`. Pool enum includes `hot-pool`.

Rows ship as `status: "pending"` with `null` metrics until a reference host is measured. Sync a copy into selenoid-ui:

```bash
cp docs/perf-benchmark/runs.json ../selenoid-ui/ui/src/perf-benchmark/runs.json
```

(Paths relative to `selenoid-tests/` when both nested repos sit under `selenoid-home/`. Note: selenoid-ui root `.gitignore` ignores `data/`, so the UI copy lives under `ui/src/perf-benchmark/`.)
