#!/usr/bin/env node
/**
 * Render notifications/config.runtime.json from template + env (ADR 008).
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const testsDir = fileURLToPath(new URL("..", import.meta.url));
const templatePath = path.join(testsDir, "notifications/config.json");
const outPath = path.join(testsDir, "notifications/config.runtime.json");

const cfg = JSON.parse(fs.readFileSync(templatePath, "utf8"));

cfg.base.project = process.env.NOTIFICATION_PROJECT || cfg.base.project;
cfg.base.links = {
  report: process.env.ALLURE_REPORT_URL || "",
  dashboard: process.env.ALLURE_DASHBOARD_URL || "",
  testops: process.env.TESTOPS_LAUNCH_URL || "",
  build: process.env.BUILD_URL || "",
};

cfg.telegram.token = process.env.TELEGRAM_BOT_TOKEN || process.env.TELEGRAM_TOKEN || "";
cfg.telegram.chat = process.env.TELEGRAM_CHAT_ID || "-1004381150566";
cfg.telegram.topic = process.env.TELEGRAM_TOPIC_ID || "";

fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, `${JSON.stringify(cfg, null, 2)}\n`);
console.log(`render-notifications-config: ${outPath}`);
