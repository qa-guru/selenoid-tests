#!/usr/bin/env bash
# Full CM pyramid slice on CI: integration → api (pre-started stack) → e2e.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

GRADLE=(./gradlew -DpyramidStand=selenoid_github -DskipHealthCheck=true)

chmod +x scripts/prepare-ci-cm-workspace.sh scripts/start-ci-cm-stack.sh scripts/stop-ci-cm-stack.sh

./scripts/prepare-ci-cm-workspace.sh
"${GRADLE[@]}" testCmIntegration
./scripts/start-ci-cm-stack.sh
"${GRADLE[@]}" testCmApi
./scripts/stop-ci-cm-stack.sh || true
"${GRADLE[@]}" testCmE2e
