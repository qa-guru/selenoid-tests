#!/usr/bin/env bash
# Allure 3 quality gate + report (replaces Gradle allureQualityGate/allureReport after P5).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

ALLURE_VERSION="${ALLURE_VERSION:-3.13.0}"
RESULTS="${ALLURE_RESULTS:-${ROOT}/build/allure-results}"
OUT="${ALLURE_REPORT_DIR:-${ROOT}/build/reports/allure-report/allureReport}"
CONFIG="${ROOT}/allurerc.mjs"

cmd="${1:-generate}"

if [[ ! -d "${RESULTS}" ]] || ! find "${RESULTS}" -maxdepth 3 -name '*-result.json' -print -quit | grep -q .; then
  echo "allure-report: no *-result.json under ${RESULTS}" >&2
  exit 1
fi

case "${cmd}" in
  quality-gate)
    npx --yes "allure@${ALLURE_VERSION}" quality-gate "${RESULTS}" --config "${CONFIG}"
    ;;
  generate)
    mkdir -p "${OUT}"
    npx --yes "allure@${ALLURE_VERSION}" generate "${RESULTS}" --config "${CONFIG}" --output "${OUT}"
    ;;
  *)
    echo "Usage: $0 quality-gate|generate" >&2
    exit 2
    ;;
esac
