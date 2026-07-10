#!/usr/bin/env bash
# Start CM-managed hub (:4445) + UI (:8081) for @Tag(api,cm) tests.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="${ROOT}/build/ci-bin"
CONFIG_DIR="${ROOT}/build/ci-cm-config"
LOG_DIR="${ROOT}/build/ci-cm-logs"
BROWSERS="${ROOT}/fixtures/ci-browsers.json"
CM="${BIN}/cm"
HUB_URL="${CM_HUB_URL:-http://127.0.0.1:4445/}"
UI_URL="${CM_UI_URL:-http://127.0.0.1:8081/}"

wait_for_url() {
  local url="$1"
  local label="$2"
  local attempts="${3:-60}"
  local i=1
  while (( i <= attempts )); do
    if curl -sf "$url" >/dev/null 2>&1; then
      echo "==> ${label} ready: ${url}"
      return 0
    fi
    sleep 2
    (( i++ )) || true
  done
  echo "${label} not ready after ${attempts} attempts: ${url}" >&2
  if [[ -d "$LOG_DIR" ]]; then
    for log in "${LOG_DIR}"/*.log; do
      [[ -f "$log" ]] || continue
      echo "--- tail ${log} ---" >&2
      tail -50 "$log" >&2 || true
    done
  fi
  "$CM" selenoid status -c "$CONFIG_DIR" -p 4445 >&2 || true
  "$CM" selenoid-ui status -c "$CONFIG_DIR" -p 8081 >&2 || true
  return 1
}

if [[ ! -x "$CM" ]]; then
  echo "Missing cm binary: ${CM} (run prepare-ci-cm-workspace.sh first)" >&2
  exit 1
fi

mkdir -p "${CONFIG_DIR}/bin" "$LOG_DIR"
cp -f "${BIN}/selenoid" "${CONFIG_DIR}/bin/selenoid"
cp -f "${BIN}/selenoid-ui" "${CONFIG_DIR}/bin/selenoid-ui"
chmod +x "${CONFIG_DIR}/bin/selenoid" "${CONFIG_DIR}/bin/selenoid-ui"

echo "==> CM configure"
"$CM" selenoid configure -c "$CONFIG_DIR" -p 4445 -n -j "$BROWSERS" \
  --selenoid-binary "${CONFIG_DIR}/bin/selenoid" \
  --selenoid-ui-binary "${CONFIG_DIR}/bin/selenoid-ui" \
  >"${LOG_DIR}/cm-configure.log" 2>&1

echo "==> CM start hub"
"$CM" selenoid start -f -c "$CONFIG_DIR" -p 4445 -n -j "$BROWSERS" \
  --selenoid-binary "${CONFIG_DIR}/bin/selenoid" \
  --selenoid-ui-binary "${CONFIG_DIR}/bin/selenoid-ui" \
  >"${LOG_DIR}/cm-hub-start.log" 2>&1

wait_for_url "${HUB_URL}status" "CM hub" 90

echo "==> CM start UI"
"$CM" selenoid-ui start -c "$CONFIG_DIR" -p 8081 \
  >"${LOG_DIR}/cm-ui-start.log" 2>&1

wait_for_url "${UI_URL}status" "CM UI" 90

echo "==> CM stack is up (hub ${HUB_URL}, UI ${UI_URL})"
