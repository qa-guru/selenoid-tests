#!/usr/bin/env bash
# Full CM pyramid slice on CI (Go; P5 cutover).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

chmod +x scripts/prepare-ci-cm-workspace.sh scripts/start-ci-cm-stack.sh scripts/stop-ci-cm-stack.sh scripts/run-go-pyramid.sh

./scripts/prepare-ci-cm-workspace.sh
./scripts/start-ci-cm-stack.sh
./scripts/run-go-pyramid.sh cm
./scripts/stop-ci-cm-stack.sh || true
