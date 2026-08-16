#!/usr/bin/env bash
# Live per-component Allure dashboards for GitHub Pages.
# Filter is Allure `component` label (same grain as epicCharts).
#
# Usage: publish-component-dashboards.sh [RESULTS_DIR] [DEST_REPORT_DIR]
# DEST_REPORT_DIR is pages/reports/<run-id> (also copied to latest).
set -euo pipefail

RESULTS_DIR="${1:-build/allure-results}"
DEST_DIR="${2:?destination report dir required}"
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

cd "${REPO_ROOT}"

if [ ! -d "${RESULTS_DIR}" ]; then
  echo "publish-component-dashboards: missing ${RESULTS_DIR}" >&2
  exit 1
fi

components=(
  selenoid
  selenoid-ui
  cm
  playwright-image
  webdriver-image
  video-recorder
  android
  ios
)

allure_generate() {
  if npx --no-install allure --version >/dev/null 2>&1; then
    npx --no-install allure generate "$@"
  else
    npx --yes "allure@${ALLURE_VERSION:-3.14.3}" generate "$@"
  fi
}

for comp in "${components[@]}"; do
  work="build/component-dashboard/${comp}"
  rm -rf "${work}"
  echo "publish-component-dashboards: generate ${comp}"
  ALLURE_COMPONENT_DASHBOARD="${comp}" \
    allure_generate "${RESULTS_DIR}" \
      --config .github/scripts/allurerc-component-dashboard.mjs \
      --output "${work}"
  if [ ! -f "${work}/dashboard/index.html" ]; then
    echo "publish-component-dashboards: no dashboard in ${work}" >&2
    exit 1
  fi
  mkdir -p "${DEST_DIR}/dashboards/${comp}"
  cp -r "${work}/dashboard/." "${DEST_DIR}/dashboards/${comp}/"
done

echo "publish-component-dashboards: done → ${DEST_DIR}/dashboards/{${components[*]}}"
