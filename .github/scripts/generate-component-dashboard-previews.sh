#!/usr/bin/env bash
# SSOT: generators/ethalon/readme/generate-component-dashboard-previews.sh
# Per-component README dashboard PNGs (playwright-image, webdriver-image).
#
# Usage: generate-component-dashboard-previews.sh [RESULTS_DIR] [OUTPUT_DIR]
set -euo pipefail

RESULTS_DIR="${1:-build/allure-results}"
OUTPUT_DIR="${2:-pages/readme}"
ALLURE_VERSION="${ALLURE_VERSION:-3.13.0}"
SERVE_PORT="${SERVE_PORT:-8767}"
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

cd "${REPO_ROOT}"

if [ ! -d "${RESULTS_DIR}" ]; then
  echo "generate-component-dashboard-previews: missing ${RESULTS_DIR}" >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"
components=(playwright-image webdriver-image)

for comp in "${components[@]}"; do
  work="build/readme-dashboard/${comp}"
  rm -rf "${work}"

  echo "generate-component-dashboard-previews: allure generate → ${work} (${comp})"
  ALLURE_COMPONENT_DASHBOARD="${comp}" \
    npx --yes "allure@${ALLURE_VERSION}" generate "${RESULTS_DIR}" \
      --config .github/scripts/allurerc-component-dashboard.mjs \
      --output "${work}"

  if [ ! -f "${work}/dashboard/index.html" ]; then
    echo "generate-component-dashboard-previews: no dashboard in ${work}" >&2
    exit 1
  fi

  npx --yes serve "${work}/dashboard" -l "${SERVE_PORT}" >/tmp/serve-readme-${comp}.log 2>&1 &
  serve_pid=$!
  cleanup() {
    kill "${serve_pid}" 2>/dev/null || true
  }
  trap cleanup EXIT

  for _ in $(seq 1 30); do
    if curl -sf "http://127.0.0.1:${SERVE_PORT}/" >/dev/null; then
      break
    fi
    sleep 1
  done

  PREVIEW_URL="http://127.0.0.1:${SERVE_PORT}/" \
    PREVIEW_OUTPUT_DIR="${OUTPUT_DIR}" \
    PREVIEW_SUFFIX="${comp}" \
    node .github/scripts/capture-dashboard-preview.mjs

  trap - EXIT
  cleanup
  wait "${serve_pid}" 2>/dev/null || true
done

echo "generate-component-dashboard-previews: done → ${OUTPUT_DIR}/dashboard-preview-{playwright-image,webdriver-image}*.png"
