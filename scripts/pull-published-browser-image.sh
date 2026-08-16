#!/usr/bin/env bash
# Pull a browser-image dispatch tag and retarget fixtures/ci-browsers.json so
# Go smoke/slice hits that image, not the prod pin already on the runner.
#
# Trigger: SOURCE_VERSION + SOURCE_VARIANT (webdriver | playwright | video-recorder).
# Tags: qaguru/webdriver-{browser}:{major}[-min]
#       qaguru/playwright-{browser}:{version}[-min]
#       qaguru/video-recorder:{version}
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BROWSERS="${BROWSERS:-${ROOT}/fixtures/ci-browsers.json}"

SOURCE_VERSION="${SOURCE_VERSION:-}"
SOURCE_VARIANT="${SOURCE_VARIANT:-}"
SOURCE_BROWSER="${SOURCE_BROWSER:-}"
WEBDRIVER_VARIANT="${WEBDRIVER_VARIANT:-}"

if [[ -z "${SOURCE_VERSION}" || -z "${SOURCE_VARIANT}" ]]; then
  echo "pull-published-browser-image: skip (need SOURCE_VERSION + SOURCE_VARIANT)"
  exit 0
fi

emit() {
  local key="$1"
  local val="$2"
  echo "${key}=${val}"
  if [[ -n "${GITHUB_ENV:-}" ]]; then
    echo "${key}=${val}" >> "${GITHUB_ENV}"
  fi
}

raw="${SOURCE_VERSION#v}"
min_suffix=""
if [[ "${WEBDRIVER_VARIANT}" == "min" ]]; then
  min_suffix="-min"
elif [[ "${WEBDRIVER_VARIANT}" == "warm" ]]; then
  min_suffix=""
elif [[ "${raw}" == *-min ]]; then
  min_suffix="-min"
fi
ver="${raw%-min}"

image=""
browser_key=""
session_key=""

case "${SOURCE_VARIANT}" in
  webdriver)
    if [[ -z "${SOURCE_BROWSER}" ]]; then
      echo "error: SOURCE_BROWSER required for source_variant=webdriver" >&2
      exit 1
    fi
    major="${ver%%.*}"
    if [[ "${ver}" == *.* ]]; then
      session_base="${ver}"
    else
      session_base="${ver}.0"
    fi
    session_key="${session_base}${min_suffix}"
    image="qaguru/webdriver-${SOURCE_BROWSER}:${major}${min_suffix}"
    browser_key="${SOURCE_BROWSER}"
    ;;
  playwright)
    pw_browser="${SOURCE_BROWSER:-chromium}"
    pw_browser="${pw_browser#playwright-}"
    session_key="${ver}${min_suffix}"
    image="qaguru/playwright-${pw_browser}:${ver}${min_suffix}"
    browser_key="playwright-${pw_browser}"
    ;;
  video-recorder)
    image="qaguru/video-recorder:${ver}"
    ;;
  *)
    echo "pull-published-browser-image: skip unknown SOURCE_VARIANT=${SOURCE_VARIANT}"
    exit 0
    ;;
esac

echo "==> docker pull published ${image} (variant=${SOURCE_VARIANT} browser=${SOURCE_BROWSER:--} webdriver_variant=${WEBDRIVER_VARIANT:-})"
docker pull "${image}"
emit PUBLISHED_BROWSER_IMAGE "${image}"

if [[ "${SOURCE_VARIANT}" == "video-recorder" ]]; then
  emit VIDEO_RECORDER_IMAGE "${image}"
  exit 0
fi

if [[ ! -f "${BROWSERS}" ]]; then
  echo "Missing browsers.json: ${BROWSERS}" >&2
  exit 1
fi

tmp="$(mktemp)"
jq --arg bk "${browser_key}" --arg vk "${session_key}" --arg img "${image}" --arg pver "${ver}" '
  def ismin: ($vk | endswith("-min"));
  if .[$bk] == null then
    error("unknown browsers.json key: \($bk)")
  else
    (.[$bk].versions | to_entries) as $ents
    | (
        ($ents | map(select((.key | endswith("-min")) == ismin)) | .[0].value)
        // $ents[0].value
      ) as $tpl
    | .[$bk].versions[$vk] = (
        $tpl
        | .image = $img
        | if has("playwrightVersion") then .playwrightVersion = $pver else . end
      )
  end
' "${BROWSERS}" > "${tmp}"
mv "${tmp}" "${BROWSERS}"
echo "==> retargeted ${BROWSERS} ${browser_key} ${session_key} → ${image}"

case "${SOURCE_VARIANT}" in
  webdriver)
    if [[ "${min_suffix}" == "-min" ]]; then
      case "${SOURCE_BROWSER}" in
        chrome) emit chromeMinVersion "${session_key}" ;;
        firefox) emit firefoxMinVersion "${session_key}" ;;
        msedge) emit msedgeMinVersion "${session_key}" ;;
      esac
    else
      case "${SOURCE_BROWSER}" in
        chrome)
          emit chromeVersion "${session_key}"
          emit browserVersion "${session_key}"
          ;;
        firefox) emit firefoxVersion "${session_key}" ;;
        msedge) emit msedgeVersion "${session_key}" ;;
      esac
    fi
    ;;
  playwright)
    emit playwrightWsEndpoint "ws://127.0.0.1:4444/playwright/${browser_key}/${session_key}"
    ;;
esac
