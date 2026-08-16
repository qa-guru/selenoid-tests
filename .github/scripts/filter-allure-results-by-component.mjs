#!/usr/bin/env node
// Copy Allure results whose `component` label matches. Widget `filter`
// functions do not survive Allure 3 dashboard generate (JSON plugin options).
//
// Usage: filter-allure-results-by-component.mjs SRC DEST COMPONENT
import fs from "node:fs";
import path from "node:path";

const [src, dest, component] = process.argv.slice(2);
if (!src || !dest || !component) {
  console.error(
    "usage: filter-allure-results-by-component.mjs SRC DEST COMPONENT",
  );
  process.exit(2);
}

function* walkFiles(dir) {
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, ent.name);
    if (ent.isDirectory()) yield* walkFiles(p);
    else yield p;
  }
}

function collectAttachmentSources(node, out) {
  if (!node || typeof node !== "object") return;
  for (const att of node.attachments ?? []) {
    if (att.source) out.add(path.basename(att.source));
  }
  for (const step of node.steps ?? []) collectAttachmentSources(step, out);
}

fs.rmSync(dest, { recursive: true, force: true });
fs.mkdirSync(dest, { recursive: true });

const attachments = new Set();
let kept = 0;
const byName = new Map();
for (const file of walkFiles(src)) {
  byName.set(path.basename(file), file);
}

for (const file of walkFiles(src)) {
  if (!file.endsWith("-result.json")) continue;
  let data;
  try {
    data = JSON.parse(fs.readFileSync(file, "utf8"));
  } catch {
    continue;
  }
  const labels = Array.isArray(data.labels) ? data.labels : [];
  const match = labels.some(
    (l) => l && l.name === "component" && l.value === component,
  );
  if (!match) continue;
  fs.copyFileSync(file, path.join(dest, path.basename(file)));
  collectAttachmentSources(data, attachments);
  kept += 1;
}

for (const name of attachments) {
  const file = byName.get(name);
  if (file) fs.copyFileSync(file, path.join(dest, name));
}

for (const extra of ["executor.json", "environment.properties", "categories.json"]) {
  const file = byName.get(extra);
  if (file) fs.copyFileSync(file, path.join(dest, extra));
}

console.log(`filter-allure-results: ${component} kept ${kept}`);
if (kept === 0) process.exit(1);
