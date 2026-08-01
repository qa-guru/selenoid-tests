import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { injectDashboardOverrides } from "../../.github/scripts/lib/inject-dashboard-overrides.mjs";

const DEFAULT_ASSETS = fileURLToPath(
  new URL("../../.github/assets", import.meta.url),
);

/**
 * Allure 3 plugin — Palette A dashboard/awesome chart overrides (pyramid, durations).
 * Etalon: @allurereport/plugin-slack `done` hook.
 */
export default class DashboardThemePlugin {
  /** @param {{ assetsDir?: string, waitMs?: number }} options */
  constructor(options = {}) {
    this.options = options;
  }

  done = async (context) => {
    const reportRoot = context.output;
    const assetsDir = path.resolve(
      process.cwd(),
      this.options.assetsDir ?? DEFAULT_ASSETS,
    );
    const waitMs = this.options.waitMs ?? 8000;

    await waitForReportHtml(reportRoot, waitMs);

    injectDashboardOverrides({
      reportRoot,
      assetsRoot: assetsDir,
      copyAssets: true,
    });
  };
}

async function waitForReportHtml(reportRoot, waitMs) {
  const marker = path.join(reportRoot, "awesome/index.html");
  const started = Date.now();
  while (Date.now() - started < waitMs) {
    if (fs.existsSync(marker)) {
      return;
    }
    await new Promise((resolve) => setTimeout(resolve, 50));
  }
  throw new Error(
    `dashboard-theme: timeout waiting for ${marker} (${waitMs}ms)`,
  );
}
