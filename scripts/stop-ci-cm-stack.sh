#!/usr/bin/env bash
# Stop CM-managed hub/UI started by start-ci-cm-stack.sh.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="${ROOT}/build/ci-bin"
CONFIG_DIR="${ROOT}/build/ci-cm-config"
CM="${BIN}/cm"

if [[ ! -x "$CM" || ! -d "$CONFIG_DIR" ]]; then
  exit 0
fi

"$CM" selenoid-ui stop -c "$CONFIG_DIR" -p 8081 2>/dev/null || true
"$CM" selenoid stop -c "$CONFIG_DIR" -p 4445 2>/dev/null || true
for port in 4445 8081; do
  ids="$(docker ps -q --filter "publish=${port}" 2>/dev/null || true)"
  if [[ -n "$ids" ]]; then
    # shellcheck disable=SC2086
    docker stop $ids >/dev/null 2>&1 || true
  fi
done
echo "==> CM stack stopped"
