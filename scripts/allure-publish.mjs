#!/usr/bin/env node
/**
 * Second Allure 3 generate pass — notifications plugin + dashboard theme.
 * Pass 1: `./gradlew allureReport` (report on disk, summary.json flushed).
 */
import { spawnSync } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const ALLURE_VERSION = "3.13.0";
const testsDir = fileURLToPath(new URL("..", import.meta.url));
const resultsDir = path.join(testsDir, "build/allure-results");
const outputDir = path.join(testsDir, "build/reports/allure-report/allureReport");

function run(cmd, args) {
  const result = spawnSync(cmd, args, {
    cwd: testsDir,
    stdio: "inherit",
    env: process.env,
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

run("node", ["scripts/render-notifications-config.mjs"]);

const mode = process.env.NOTIFICATION_MODE || "dry-run";
console.log(`allure-publish: second generate (notifications=${mode})`);

run("npx", [
  "--yes",
  `allure@${ALLURE_VERSION}`,
  "generate",
  resultsDir,
  "-o",
  outputDir,
  "--config",
  "allurerc.publish.mjs",
]);

console.log("allure-publish: done");
