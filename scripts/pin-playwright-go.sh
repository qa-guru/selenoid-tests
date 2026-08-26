#!/usr/bin/env bash
# Pin github.com/mxschmitt/playwright-go to the Playwright protocol of SOURCE_VERSION.
#
# playwright-go v0.MM00.0 tracks Playwright 1.MM.x (v0.6100.0 ↔ 1.61, v0.6000.0 ↔ 1.60).
# Protocol match is major.minor: a 1.61.0 client is OK with a 1.61.1 server.
# Without this pin, go.mod's v0.6100.0 client gets 428 from a 1.60.x launchServer.
#
# Trigger: SOURCE_VARIANT=playwright + SOURCE_VERSION (same as pull-published-browser-image.sh).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

playwright_go_modver() {
  local raw="${1#v}"
  raw="${raw%-min}"
  if [[ ! "${raw}" =~ ^([0-9]+)\.([0-9]+)(\.[0-9]+)?$ ]]; then
    echo "pin-playwright-go: unparseable Playwright version: ${1}" >&2
    return 1
  fi
  local maj="${BASH_REMATCH[1]}"
  local min="${BASH_REMATCH[2]}"
  if [[ "${maj}" != "1" ]]; then
    echo "pin-playwright-go: unsupported Playwright major ${maj} (want 1.MM.x)" >&2
    return 1
  fi
  printf 'v0.%d.0\n' "$((10#${min} * 100))"
}

self_test() {
  local got
  got="$(playwright_go_modver 1.60.0)"
  [[ "${got}" == "v0.6000.0" ]] || { echo "fail 1.60.0 → ${got}" >&2; exit 1; }
  got="$(playwright_go_modver 1.61.1)"
  [[ "${got}" == "v0.6100.0" ]] || { echo "fail 1.61.1 → ${got}" >&2; exit 1; }
  got="$(playwright_go_modver v1.61.0-min)"
  [[ "${got}" == "v0.6100.0" ]] || { echo "fail v1.61.0-min → ${got}" >&2; exit 1; }
  got="$(playwright_go_modver 1.52.0)"
  [[ "${got}" == "v0.5200.0" ]] || { echo "fail 1.52.0 → ${got}" >&2; exit 1; }
  if playwright_go_modver not-a-version >/dev/null 2>&1; then
    echo "fail: unparseable version should error" >&2
    exit 1
  fi
  echo "pin-playwright-go: self-test ok"
}

if [[ "${1:-}" == "--self-test" ]]; then
  self_test
  exit 0
fi

SOURCE_VERSION="${SOURCE_VERSION:-}"
SOURCE_VARIANT="${SOURCE_VARIANT:-}"

if [[ "${SOURCE_VARIANT}" != "playwright" || -z "${SOURCE_VERSION}" ]]; then
  echo "pin-playwright-go: skip (need SOURCE_VARIANT=playwright and SOURCE_VERSION)"
  exit 0
fi

mod_ver="$(playwright_go_modver "${SOURCE_VERSION}")"

if [[ "${1:-}" == "--print" ]]; then
  echo "${mod_ver}"
  exit 0
fi

echo "pin-playwright-go: Playwright ${SOURCE_VERSION} → github.com/mxschmitt/playwright-go@${mod_ver}"
cd "${ROOT}"
go get "github.com/mxschmitt/playwright-go@${mod_ver}"
go mod tidy
echo "pin-playwright-go: $(go list -m github.com/mxschmitt/playwright-go)"
