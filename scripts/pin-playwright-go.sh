#!/usr/bin/env bash
# Pin github.com/mxschmitt/playwright-go to the Playwright protocol of SOURCE_VERSION.
#
# Bindings are tagged v0.MMxx.y (v0.6100.0 ↔ 1.61, v0.6000.0 ↔ 1.60, v0.6201.1 ↔ 1.62).
# Protocol match is major.minor: a 1.61.0 client is OK with a 1.61.1 server, not with 1.60/1.62.
# Without this pin, go.mod's v0.6100.0 client gets 428 from a non-1.61 launchServer.
#
# Trigger: SOURCE_VARIANT=playwright + SOURCE_VERSION (same as pull-published-browser-image.sh).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MODULE="github.com/mxschmitt/playwright-go"

# Latest playwright-go tag for Playwright 1.${min}.x from a space-separated version list.
pick_from_versions() {
  local min="$1"
  local versions_str="$2"
  local match=""
  local v
  for v in ${versions_str}; do
    if [[ "${v}" =~ ^v0\.${min}[0-9]{2}\.[0-9]+$ ]]; then
      match="${v}"
    fi
  done
  if [[ -z "${match}" ]]; then
    echo "pin-playwright-go: no ${MODULE} tag for Playwright 1.${min}.x in: ${versions_str}" >&2
    return 1
  fi
  printf '%s\n' "${match}"
}

playwright_minor() {
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
  printf '%s\n' "${min}"
}

playwright_go_modver() {
  local min
  min="$(playwright_minor "$1")"
  local versions
  versions="$(go list -m -versions "${MODULE}")"
  pick_from_versions "${min}" "${versions}"
}

# v0.6000.0 (PW 1.60) declares module github.com/playwright-community/playwright-go, so
# `go get mxschmitt@v0.6000.0` fails. Copy it locally and rewrite imports to the path
# selenoid-tests already uses. 1.61+ tags declare mxschmitt again and go get works.
rewrite_community_to_mxschmitt() {
  local ver="$1"
  local community="github.com/playwright-community/playwright-go"
  local dest="${ROOT}/.ci-pin/playwright-go"
  echo "pin-playwright-go: ${MODULE}@${ver} path mismatch; rewriting ${community}@${ver} → ${dest}"
  local dir
  dir="$(go mod download -json "${community}@${ver}" | python3 -c 'import json,sys; print(json.load(sys.stdin)["Dir"])')"
  rm -rf "${dest}"
  mkdir -p "${ROOT}/.ci-pin"
  cp -a "${dir}" "${dest}"
  chmod -R u+w "${dest}"
  local sed_inplace=(sed -i)
  if ! sed --version >/dev/null 2>&1; then
    sed_inplace=(sed -i '')
  fi
  find "${dest}" \( -name '*.go' -o -name 'go.mod' \) -exec "${sed_inplace[@]}" \
    's|github.com/playwright-community/playwright-go|github.com/mxschmitt/playwright-go|g' {} +
  overlay_npm_driver_installer "${dest}"
  go mod edit -replace "${MODULE}=./.ci-pin/playwright-go"
}

# v0.6000.0 still fetches playwright-1.60.0-linux.zip from retired azureedge CDNs (404).
# Overlay the v0.6100.0 npm/nodejs installer and retarget playwright-core to this tag's CLI version.
overlay_npm_driver_installer() {
  local dest="$1"
  local pw_ver
  pw_ver="$(grep -Eo 'playwrightCliVersion = "[0-9]+\.[0-9]+\.[0-9]+"' "${dest}/run.go" | head -1 | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')"
  if [[ -z "${pw_ver}" ]]; then
    echo "pin-playwright-go: could not read playwrightCliVersion from ${dest}/run.go" >&2
    return 1
  fi
  echo "pin-playwright-go: overlay npm driver installer, playwright-core@${pw_ver}"
  curl -fsSL "https://raw.githubusercontent.com/playwright-community/playwright-go/v0.6100.0/run.go" -o "${dest}/run.go"
  if sed --version >/dev/null 2>&1; then
    sed -i "s/playwrightCliVersion = \"1.61.1\"/playwrightCliVersion = \"${pw_ver}\"/" "${dest}/run.go"
  else
    sed -i '' "s/playwrightCliVersion = \"1.61.1\"/playwrightCliVersion = \"${pw_ver}\"/" "${dest}/run.go"
  fi
  if ! grep -q "playwrightCliVersion = \"${pw_ver}\"" "${dest}/run.go"; then
    echo "pin-playwright-go: failed to retarget npm installer to ${pw_ver}" >&2
    return 1
  fi
}

pin_module() {
  local ver="$1"
  if go get "${MODULE}@${ver}"; then
    go mod tidy
    return 0
  fi
  rewrite_community_to_mxschmitt "${ver}"
  go mod tidy
}

self_test() {
  local fixture="v0.5200.0 v0.5700.0 v0.5700.1 v0.6000.0 v0.6100.0 v0.6201.0 v0.6201.1"
  local got
  got="$(pick_from_versions 60 "${fixture}")"
  [[ "${got}" == "v0.6000.0" ]] || { echo "fail 1.60 → ${got}" >&2; exit 1; }
  got="$(pick_from_versions 61 "${fixture}")"
  [[ "${got}" == "v0.6100.0" ]] || { echo "fail 1.61 → ${got}" >&2; exit 1; }
  got="$(pick_from_versions 62 "${fixture}")"
  [[ "${got}" == "v0.6201.1" ]] || { echo "fail 1.62 → ${got}" >&2; exit 1; }
  got="$(pick_from_versions 57 "${fixture}")"
  [[ "${got}" == "v0.5700.1" ]] || { echo "fail 1.57 → ${got}" >&2; exit 1; }
  got="$(playwright_minor 1.60.0)"
  [[ "${got}" == "60" ]] || { echo "fail minor 1.60.0 → ${got}" >&2; exit 1; }
  got="$(playwright_minor v1.62.1-min)"
  [[ "${got}" == "62" ]] || { echo "fail minor v1.62.1-min → ${got}" >&2; exit 1; }
  if playwright_minor not-a-version >/dev/null 2>&1; then
    echo "fail: unparseable version should error" >&2
    exit 1
  fi
  if pick_from_versions 99 "${fixture}" >/dev/null 2>&1; then
    echo "fail: missing series should error" >&2
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

echo "pin-playwright-go: Playwright ${SOURCE_VERSION} → ${MODULE}@${mod_ver}"
cd "${ROOT}"
pin_module "${mod_ver}"
echo "pin-playwright-go: $(go list -m "${MODULE}")"
