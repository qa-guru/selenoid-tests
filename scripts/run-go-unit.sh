#!/usr/bin/env bash
# Run Go unit tests for one service repo and export Allure results via native
# `go test -json` → cmd/gotest2allure (no Node/JUnit bridge).
set -euo pipefail

REPO="${1:?repo name: selenoid|selenoid-ui|cm}"
EPIC="${2:?epic label, e.g. selenoid}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPOS_DIR="${ROOT}/repos"
JSON_DIR="${ROOT}/build/gotest-json"
ALLURE_DIR="${ROOT}/build/allure-results/go-${REPO}"

mkdir -p "${JSON_DIR}" "${ALLURE_DIR}"
JSON_FILE="${JSON_DIR}/go-${REPO}.jsonl"

export GO111MODULE=on
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export PATH="${PATH}:$(go env GOPATH)/bin"

cd "${REPOS_DIR}/${REPO}"

set +e
case "${REPO}" in
  selenoid)
    go test -json -tags 's3 metadata' -race -coverprofile=coverage.txt -covermode=atomic \
      -coverpkg github.com/qa-guru/selenoid,github.com/qa-guru/selenoid/session,github.com/qa-guru/selenoid/config,github.com/qa-guru/selenoid/protect,github.com/qa-guru/selenoid/service,github.com/qa-guru/selenoid/upload,github.com/qa-guru/selenoid/info,github.com/qa-guru/selenoid/jsonerror \
      ./... >"${JSON_FILE}"
    ;;
  selenoid-ui)
    test -f ui/build/index.html || {
      echo "ui/build/index.html missing — run yarn build in ui/" >&2
      exit 1
    }
    if [ -d ui/allure-results ] && [ "$(ls -A ui/allure-results 2>/dev/null)" ]; then
      cp -R ui/allure-results/. "${ALLURE_DIR}/"
    fi
    go install github.com/rakyll/statik@latest
    go generate github.com/qa-guru/selenoid-ui
    go test -json -race -coverprofile=coverage.txt -covermode=atomic ./... >"${JSON_FILE}"
    ;;
  cm)
    go test -json -race -coverprofile=coverage.txt -covermode=atomic \
      -coverpkg github.com/qa-guru/cm/...,github.com/qa-guru/cm/cmd \
      ./... >"${JSON_FILE}"
    ;;
  *)
    echo "Unknown repo: ${REPO}" >&2
    exit 2
    ;;
esac
TEST_EXIT=$?
set -e

if [[ ! -s "${JSON_FILE}" ]]; then
  echo "go test -json output missing/empty: ${JSON_FILE}" >&2
  exit "${TEST_EXIT:-1}"
fi

cd "${ROOT}"
go run ./cmd/gotest2allure \
  --input "${JSON_FILE}" \
  --output "${ALLURE_DIR}" \
  --epic "${EPIC}" \
  --component "${EPIC}" \
  --layer unit

exit "${TEST_EXIT}"
