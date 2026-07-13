#!/usr/bin/env bash
# CM integration CI: sibling repo layout + build cm/selenoid/selenoid-ui + pull images.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PARENT="$(cd "${ROOT}/.." && pwd)"
REPOS="${ROOT}/repos"
BIN="${ROOT}/build/ci-bin"
BROWSERS="${ROOT}/fixtures/ci-browsers.json"

# Do NOT pin DOCKER_API_VERSION: the CI runner's Docker daemon maxes at API 1.48,
# while the v2.3.0 hub canon is 1.55. moby client + docker CLI auto-negotiate when
# unset. Honor an explicit caller override if provided.
if [[ -n "${DOCKER_API_VERSION:-}" ]]; then
  export DOCKER_API_VERSION
fi
export GO111MODULE=on
export PATH="$(go env GOPATH)/bin:${PATH:-}"

mkdir -p "$BIN" "${PARENT}/dev"

require_repo() {
  local name="$1"
  if [[ ! -d "${REPOS}/${name}" ]]; then
    echo "Missing repo checkout: ${REPOS}/${name}" >&2
    exit 1
  fi
}

link_repo() {
  local name="$1"
  local target="${PARENT}/${name}"
  if [[ ! -e "$target" ]]; then
    ln -s "${REPOS}/${name}" "$target"
  fi
}

build_selenoid_ui_frontend() {
  local ui="${REPOS}/selenoid-ui/ui"
  export CI=false
  unset NODE_OPTIONS
  if [[ -f "${ui}/build/index.html" ]]; then
    echo "==> reusing existing selenoid-ui ui/build"
    return 0
  fi
  echo "==> Building selenoid-ui frontend (Vite 6 / React 18)"
  yarn --cwd "$ui" install --frozen-lockfile
  # v2.3.0 UI is Vite (yarn build → ui/build/); works on modern Node (CI Node 24).
  yarn --cwd "$ui" build
}

pull_browser_images() {
  echo "==> Pulling CI smoke images (chrome 149 + playwright-chromium) from ${BROWSERS}"
  while IFS= read -r img; do
    [[ -n "$img" ]] || continue
    echo "    docker pull ${img}"
    docker pull "$img"
  done < <(jq -r '
    .chrome.versions["149.0"].image // empty,
    .chrome.versions["149.0-min"].image // empty,
    .["playwright-chromium"].versions["1.61.1"].image // empty,
    .["playwright-chromium"].versions["1.61.1-min"].image // empty
  ' "$BROWSERS" | sort -u)
}

require_repo cm
require_repo selenoid
require_repo selenoid-ui

link_repo cm
link_repo selenoid
link_repo selenoid-ui
cp -f "$BROWSERS" "${PARENT}/dev/browsers.json"

echo "==> Building cm"
(cd "${REPOS}/cm" && go build -o "${BIN}/cm" .)

echo "==> Building selenoid hub binary (linux/amd64)"
(cd "${REPOS}/selenoid" && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "${BIN}/selenoid" .)

build_selenoid_ui_frontend
echo "==> Building selenoid-ui binary (linux/amd64)"
go install github.com/rakyll/statik@latest
(cd "${REPOS}/selenoid-ui" && go generate . && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "${BIN}/selenoid-ui" .)

pull_browser_images

echo "    docker pull qaguru/selenoid:latest-release"
docker pull qaguru/selenoid:latest-release
echo "    docker pull qaguru/selenoid-ui:latest-release"
docker pull qaguru/selenoid-ui:latest-release
echo "    docker pull qaguru/video-recorder:latest"
docker pull qaguru/video-recorder:latest

echo "==> CM workspace ready (binaries in ${BIN}, repos linked under ${PARENT})"
