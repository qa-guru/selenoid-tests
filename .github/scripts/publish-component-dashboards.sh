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

# Allure 3 with only the dashboard plugin writes index.html at --output root.
# With awesome+dashboard it nests under --output/dashboard/.
dashboard_src() {
  local work="$1"
  if [ -f "${work}/dashboard/index.html" ]; then
    printf '%s\n' "${work}/dashboard"
    return 0
  fi
  if [ -f "${work}/index.html" ]; then
    printf '%s\n' "${work}"
    return 0
  fi
  return 1
}

published=()
for comp in "${components[@]}"; do
  work="build/component-dashboard/${comp}"
  filtered="build/component-results/${comp}"
  rm -rf "${work}" "${filtered}"
  if ! node .github/scripts/filter-allure-results-by-component.mjs \
    "${RESULTS_DIR}" "${filtered}" "${comp}"; then
    echo "publish-component-dashboards: skip ${comp} (no results)" >&2
    continue
  fi
  echo "publish-component-dashboards: generate ${comp}"
  ALLURE_COMPONENT_DASHBOARD="${comp}" \
    allure_generate "${filtered}" \
      --config .github/scripts/allurerc-component-dashboard.mjs \
      --output "${work}"
  src="$(dashboard_src "${work}" || true)"
  if [ -z "${src}" ]; then
    echo "publish-component-dashboards: skip ${comp} (no index.html under ${work})" >&2
    continue
  fi
  mkdir -p "${DEST_DIR}/dashboards/${comp}"
  cp -R "${src}/." "${DEST_DIR}/dashboards/${comp}/"
  published+=("${comp}")
done

if [ "${#published[@]}" -eq 0 ]; then
  echo "publish-component-dashboards: no component dashboards produced" >&2
  exit 1
fi

echo "publish-component-dashboards: done → ${DEST_DIR}/dashboards/{${published[*]}}"
