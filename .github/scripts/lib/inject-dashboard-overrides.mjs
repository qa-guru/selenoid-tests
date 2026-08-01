import fs from "node:fs";
import path from "node:path";

const TARGETS = ["awesome/index.html", "dashboard/index.html"];
const ASSET_NAMES = ["dashboard-overrides.css", "dashboard-overrides.js"];

/**
 * Copy Palette A overrides into an Allure 3 report and patch HTML heads.
 * Assets land at report root; awesome/dashboard link via `../`.
 *
 * @param {object} opts
 * @param {string} opts.reportRoot
 * @param {string} opts.assetsRoot
 * @param {boolean} [opts.copyAssets=true]
 */
export function injectDashboardOverrides({
  reportRoot,
  assetsRoot,
  copyAssets = true,
}) {
  if (copyAssets) {
    for (const name of ASSET_NAMES) {
      const src = path.join(assetsRoot, name);
      const dest = path.join(reportRoot, name);
      if (!fs.existsSync(src)) {
        throw new Error(`inject-dashboard-overrides: missing ${src}`);
      }
      fs.mkdirSync(reportRoot, { recursive: true });
      fs.copyFileSync(src, dest);
    }
  }

  for (const rel of TARGETS) {
    patchHtmlFile({
      filePath: path.join(reportRoot, rel),
      cssHref: "../dashboard-overrides.css",
      jsHref: "../dashboard-overrides.js",
    });
  }
}

function patchHtmlFile({ filePath, cssHref, jsHref }) {
  if (!fs.existsSync(filePath)) {
    console.log(`inject-dashboard-overrides: skip missing ${filePath}`);
    return;
  }

  let html = fs.readFileSync(filePath, "utf8");

  if (hasCorrectOverrides(html, cssHref, jsHref)) {
    console.log(`inject-dashboard-overrides: already patched ${filePath}`);
    return;
  }

  if (html.includes("dashboard-overrides.js") || html.includes("dashboard-overrides.css")) {
    html = stripExistingOverrides(html);
    console.log(`inject-dashboard-overrides: upgraded paths in ${filePath}`);
  }

  const inject = buildInject(cssHref, jsHref);
  if (!html.includes("</head>")) {
    throw new Error(`inject-dashboard-overrides: no </head> in ${filePath}`);
  }
  html = html.replace("</head>", `${inject}\n</head>`);
  fs.writeFileSync(filePath, html);
  console.log(`inject-dashboard-overrides: patched ${filePath} (${cssHref})`);
}

function buildInject(cssHref, jsHref) {
  return [
    `    <link rel="stylesheet" type="text/css" href="${cssHref}" data-dashboard-overrides>`,
    `    <script src="${jsHref}" defer data-dashboard-overrides></script>`,
  ].join("\n");
}

function stripExistingOverrides(html) {
  return html
    .replace(/^\s*<link[^>]*data-dashboard-overrides[^>]*>\n?/gm, "")
    .replace(/^\s*<script[^>]*data-dashboard-overrides[^>]*><\/script>\n?/gm, "");
}

function hasCorrectOverrides(html, cssHref, jsHref) {
  return html.includes(`href="${cssHref}"`) && html.includes(`src="${jsHref}"`);
}
