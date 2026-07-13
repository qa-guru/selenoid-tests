#!/usr/bin/env bash
# Bootstrap Selenoid hub (:4444) + UI (:8080) on GitHub Actions / CI runner.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPOS="${ROOT}/repos"
BIN="${ROOT}/build/ci-bin"
LOG_DIR="${ROOT}/build/ci-logs"
BROWSERS="${ROOT}/fixtures/ci-browsers.json"
HUB_URL="${HUB_URL:-http://127.0.0.1:4444/}"
UI_URL="${UI_URL:-http://127.0.0.1:8080/}"

# Do NOT pin DOCKER_API_VERSION: the CI runner's Docker daemon maxes at API 1.48,
# while the v2.3.0 hub canon is 1.55. moby client + docker CLI auto-negotiate the
# API version when it is unset. Honor an explicit caller override if provided.
if [[ -n "${DOCKER_API_VERSION:-}" ]]; then
  export DOCKER_API_VERSION
fi
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GO111MODULE=on
export PATH="$(go env GOPATH)/bin:${PATH:-}"

mkdir -p "$BIN" "$LOG_DIR"

stop_existing_stack() {
  for name in selenoid selenoid-ui; do
    local pid_file="${LOG_DIR}/${name}.pid"
    if [[ -f "$pid_file" ]]; then
      local pid
      pid="$(cat "$pid_file")"
      if kill -0 "$pid" 2>/dev/null; then
        echo "==> Stopping existing ${name} (pid ${pid})"
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
      fi
      rm -f "$pid_file"
    fi
  done
}

require_repo() {
  local name="$1"
  if [[ ! -d "${REPOS}/${name}" ]]; then
    echo "Missing repo checkout: ${REPOS}/${name}" >&2
    exit 1
  fi
}

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
  [[ -f "${LOG_DIR}/selenoid.log" ]] && tail -50 "${LOG_DIR}/selenoid.log" >&2 || true
  [[ -f "${LOG_DIR}/selenoid-ui.log" ]] && tail -50 "${LOG_DIR}/selenoid-ui.log" >&2 || true
  return 1
}

build_selenoid() {
  echo "==> Building selenoid hub"
  (cd "${REPOS}/selenoid" && go build -o "${BIN}/selenoid" .)
}

build_selenoid_ui() {
  echo "==> Building selenoid-ui frontend (Vite 6 / React 18)"
  local ui="${REPOS}/selenoid-ui/ui"
  export CI=false
  unset NODE_OPTIONS
  if [[ -f "${ui}/build/index.html" ]]; then
    echo "    reusing existing ui/build (skip yarn build)"
  else
    yarn --cwd "$ui" install --frozen-lockfile
    # v2.3.0 UI is Vite (yarn build → ui/build/); works on modern Node (CI Node 24).
    yarn --cwd "$ui" build
  fi

  echo "==> Building selenoid-ui binary"
  go install github.com/rakyll/statik@latest
  (cd "${REPOS}/selenoid-ui" && go generate . && go build -o "${BIN}/selenoid-ui" .)
}

pull_browser_images() {
  if [[ "${SKIP_BROWSER_PULL:-}" == "1" ]]; then
    echo "==> SKIP_BROWSER_PULL=1 — skipping docker pull"
    return 0
  fi
  echo "==> Pulling CI smoke images (chrome 149 + firefox/msedge min + playwright) from ${BROWSERS}"
  if [[ ! -f "$BROWSERS" ]]; then
    echo "Missing browsers.json: ${BROWSERS}" >&2
    exit 1
  fi
  while IFS= read -r img; do
    [[ -n "$img" ]] || continue
    echo "    docker pull ${img}"
    docker pull "$img"
  done < <(jq -r '
    .chrome.versions["149.0"].image // empty,
    .chrome.versions["149.0-min"].image // empty,
    .firefox.versions[.firefox.default].image // empty,
    .firefox.versions[(.firefox.default + "-min")].image // empty,
    .msedge.versions[.msedge.default].image // empty,
    .msedge.versions[(.msedge.default + "-min")].image // empty,
    .["playwright-chromium"].versions["1.61.1"].image // empty,
    .["playwright-chromium"].versions["1.61.1-min"].image // empty,
    .["playwright-firefox"].versions["1.61.1"].image // empty,
    .["playwright-webkit"].versions["1.61.1"].image // empty
  ' "$BROWSERS" | sort -u)
  echo "    docker pull qaguru/video-recorder:latest"
  docker pull qaguru/video-recorder:latest
}

start_stack() {
  echo "==> Starting Selenoid hub"
  nohup "${BIN}/selenoid" \
    -conf "$BROWSERS" \
    -limit 3 \
    -video-recorder-image qaguru/video-recorder:latest \
    >"${LOG_DIR}/selenoid.log" 2>&1 &
  echo $! >"${LOG_DIR}/selenoid.pid"

  echo "==> Starting Selenoid UI"
  nohup "${BIN}/selenoid-ui" \
    -listen :8080 \
    -selenoid-uri "${HUB_URL%/}" \
    -browsers-conf "$BROWSERS" \
    -period 4s >"${LOG_DIR}/selenoid-ui.log" 2>&1 &
  echo $! >"${LOG_DIR}/selenoid-ui.pid"

  wait_for_url "${HUB_URL}status" "Hub"
  wait_for_url "${UI_URL}status" "UI"
}

require_repo selenoid
require_repo selenoid-ui
stop_existing_stack
build_selenoid
build_selenoid_ui
pull_browser_images
start_stack

echo "==> CI Selenoid stack is up (hub ${HUB_URL}, UI ${UI_URL})"
