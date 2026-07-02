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

export DOCKER_API_VERSION="${DOCKER_API_VERSION:-1.45}"
export GO111MODULE=on
export PATH="$(go env GOPATH)/bin:${PATH:-}"

mkdir -p "$BIN" "$LOG_DIR"

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
  echo "==> Building selenoid-ui frontend"
  local ui="${REPOS}/selenoid-ui/ui"
  export CI=false
  # GHA blocks --openssl-legacy-provider in NODE_OPTIONS; pass as CLI flag on Node 17+.
  unset NODE_OPTIONS
  yarn --cwd "$ui" install --frozen-lockfile
  local node_major
  node_major="$(node -p "process.versions.node.split('.')[0]")"
  if (( node_major >= 17 )); then
    (cd "$ui" && node --openssl-legacy-provider ./node_modules/.bin/react-scripts build)
  else
    yarn --cwd "$ui" build
  fi

  echo "==> Building selenoid-ui binary"
  go install github.com/rakyll/statik@latest
  (cd "${REPOS}/selenoid-ui" && go generate . && go build -o "${BIN}/selenoid-ui" .)
}

pull_browser_images() {
  echo "==> Pulling browser images from ${BROWSERS}"
  if [[ ! -f "$BROWSERS" ]]; then
    echo "Missing browsers.json: ${BROWSERS}" >&2
    exit 1
  fi
  while IFS= read -r img; do
    [[ -n "$img" ]] || continue
    echo "    docker pull ${img}"
    docker pull "$img"
  done < <(jq -r '.. | objects | select(has("image")) | .image' "$BROWSERS" | sort -u)
}

start_stack() {
  echo "==> Starting Selenoid hub"
  nohup "${BIN}/selenoid" -conf "$BROWSERS" -limit 3 >"${LOG_DIR}/selenoid.log" 2>&1 &
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
build_selenoid
build_selenoid_ui
pull_browser_images
start_stack

echo "==> CI Selenoid stack is up (hub ${HUB_URL}, UI ${UI_URL})"
