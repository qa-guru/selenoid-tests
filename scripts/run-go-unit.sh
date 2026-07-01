#!/usr/bin/env bash
# Run Go unit tests for one service repo and export Allure results (Variant A: gotestsum + JUnit reader).
set -euo pipefail

REPO="${1:?repo name: selenoid|selenoid-ui|cm}"
EPIC="${2:?epic label, e.g. selenoid}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPOS_DIR="${ROOT}/repos"
JUNIT_DIR="${ROOT}/build/junit"
ALLURE_DIR="${ROOT}/build/allure-results/go-${REPO}"

mkdir -p "${JUNIT_DIR}" "${ALLURE_DIR}"
JUNIT_FILE="${JUNIT_DIR}/go-${REPO}.xml"

export GO111MODULE=on
export PATH="${PATH}:$(go env GOPATH)/bin"

if ! command -v gotestsum >/dev/null 2>&1; then
  go install gotest.tools/gotestsum@v1.12.3
fi

cd "${REPOS_DIR}/${REPO}"

set +e
case "${REPO}" in
  selenoid)
    gotestsum --junitfile "${JUNIT_FILE}" -- \
      -tags 's3 metadata' -race -coverprofile=coverage.txt -covermode=atomic \
        -coverpkg github.com/aerokube/selenoid,github.com/aerokube/selenoid/session,github.com/aerokube/selenoid/config,github.com/aerokube/selenoid/protect,github.com/aerokube/selenoid/service,github.com/aerokube/selenoid/upload,github.com/aerokube/selenoid/info,github.com/aerokube/selenoid/jsonerror \
        ./...
    ;;
  selenoid-ui)
    test -f ui/build/index.html || {
      echo "ui/build/index.html missing — run yarn build in ui/" >&2
      exit 1
    }
    go install github.com/rakyll/statik@latest
    go generate github.com/aerokube/selenoid-ui
    gotestsum --junitfile "${JUNIT_FILE}" -- \
      -race -coverprofile=coverage.txt -covermode=atomic ./...
    ;;
  cm)
    gotestsum --junitfile "${JUNIT_FILE}" -- \
      -race -coverprofile=coverage.txt -covermode=atomic \
        -coverpkg github.com/aerokube/cm/selenoid \
        github.com/aerokube/cm/selenoid
    ;;
  *)
    echo "Unknown repo: ${REPO}" >&2
    exit 2
    ;;
esac
TEST_EXIT=$?
set -e

if [ ! -f "${JUNIT_FILE}" ]; then
  echo "JUnit file not found: ${JUNIT_FILE}" >&2
  exit "${TEST_EXIT:-1}"
fi

cd "${ROOT}"
npm install --no-save @allurereport/reader@3.14.0 @allurereport/reader-api@3.14.0
node scripts/junit-to-allure.mjs \
  --input "${JUNIT_FILE}" \
  --output "${ALLURE_DIR}" \
  --epic "${EPIC}" \
  --layer unit

exit "${TEST_EXIT}"
