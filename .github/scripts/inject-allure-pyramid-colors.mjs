#!/usr/bin/env node
/**
 * Inject Palette A pyramid colors into Allure 3 awesome/dashboard HTML.
 * Requires dashboard-overrides.js/css at site root (pages/).
 *
 * Usage:
 *   node inject-allure-pyramid-colors.mjs <pagesDir> <reportRoot>
 *   node inject-allure-pyramid-colors.mjs pages pages/reports/latest
 */
import fs from "node:fs";
import path from "node:path";

const pagesDir = process.argv[2];
const reportRoot = process.argv[3];

if (!pagesDir || !reportRoot) {
  console.error("usage: inject-allure-pyramid-colors.mjs <pagesDir> <reportRoot>");
  process.exit(1);
}

const ABSOLUTE_CSS = /href="\/dashboard-overrides\.css"/g;
const ABSOLUTE_JS = /src="\/dashboard-overrides\.js"/g;
const RELATIVE_CSS = /href="[^"]*dashboard-overrides\.css"/g;
const RELATIVE_JS = /src="[^"]*dashboard-overrides\.js"/g;

function assetHref(htmlDir, assetName) {
  const rel = path.relative(htmlDir, pagesDir);
  const prefix = rel ? `${rel.split(path.sep).join("/")}/` : "";
  return `${prefix}${assetName}`;
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

function needsUpgrade(html) {
  return ABSOLUTE_CSS.test(html) || ABSOLUTE_JS.test(html);
}

const targets = ["awesome/index.html", "dashboard/index.html"];

for (const rel of targets) {
  const filePath = path.join(reportRoot, rel);
  if (!fs.existsSync(filePath)) {
    console.log(`inject-pyramid-colors: skip missing ${filePath}`);
    continue;
  }

  const htmlDir = path.dirname(filePath);
  const cssHref = assetHref(htmlDir, "dashboard-overrides.css");
  const jsHref = assetHref(htmlDir, "dashboard-overrides.js");
  let html = fs.readFileSync(filePath, "utf8");

  if (hasCorrectOverrides(html, cssHref, jsHref)) {
    console.log(`inject-pyramid-colors: already patched ${filePath}`);
    continue;
  }

  if (html.includes("dashboard-overrides.js") || html.includes("dashboard-overrides.css")) {
    html = stripExistingOverrides(html);
    console.log(`inject-pyramid-colors: upgraded absolute paths in ${filePath}`);
  }

  const inject = buildInject(cssHref, jsHref);
  html = html.replace("</head>", `${inject}\n</head>`);
  fs.writeFileSync(filePath, html);
  console.log(`inject-pyramid-colors: patched ${filePath} (${cssHref})`);
}
