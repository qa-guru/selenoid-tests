import { fileURLToPath } from "node:url";

import { withKit, theme, renderers } from "@qa-guru/allure-report-kit";

import { buildAwesomeCharts } from "./awesome-charts.mjs";
import { categoryRules } from "./categories.mjs";
import {
  HISTORY_DEFAULTS,
  QUALITY_GATE_SOURCE,
  REPORT_LANGUAGE,
} from "./constants.mjs";
import { buildDashboardLayout } from "./dashboard-layout.mjs";
import { qualityGateRules } from "./quality-gate.mjs";

const DASHBOARD_THEME_PLUGIN = fileURLToPath(
  new URL("./plugins/dashboard-theme.mjs", import.meta.url),
);

/**
 * Build Allure 3 config from ethalon modules.
 *
 * HTML theme and chart renderers — @qa-guru/allure-report-kit (soft-fork).
 *
 * @param {object} profile
 * @param {string} profile.slug - repo slug → `{slug} Tests`
 * @param {string[]} [profile.epicCharts] - optional per-epic statusDynamics tiles
 * @param {object} [profile.variables] - Allure variables override
 * @param {object} [profile.publish] - CI publish plugins (notifications)
 * @param {object} [profile.publish.notifications] - `@allure-notifications/plugin` options
 */
export function createAllureConfig({
  slug,
  epicCharts = [],
  variables,
  publish,
} = {}) {
  if (!slug || typeof slug !== "string") {
    throw new Error("createAllureConfig: profile.slug is required");
  }

  return withKit({
    softFork: true,
    renderer: renderers.stock(),
    theme: {
      ...theme.qaGuru(),
      header: {
        enabled: true,
        source: "design-system",
        productName: `${slug} Tests`,
      },
    },
    name: `${slug} Tests`,
    ...HISTORY_DEFAULTS,
    variables: variables ?? {
      Framework: "Go + Allure 3",
      Report: "Allure 3",
    },
    qualityGate: {
      rules: qualityGateRules.map((rule) => ({ ...rule })),
      source: { ...QUALITY_GATE_SOURCE },
    },
    categories: {
      rules: categoryRules.map((rule) => structuredClone(rule)),
    },
    plugins: {
      awesome: {
        options: {
          reportLanguage: REPORT_LANGUAGE,
          groupBy: ["parentSuite", "suite", "subSuite"],
          charts: buildAwesomeCharts(),
        },
      },
      dashboard: {
        options: {
          reportName: `${slug} Tests Dashboard`,
          reportLanguage: REPORT_LANGUAGE,
          layout: buildDashboardLayout({ epicCharts }),
        },
      },
      csv: {
        options: {
          fileName: `${slug}.csv`,
        },
      },
      dashboardTheme: {
        import: DASHBOARD_THEME_PLUGIN,
        options: {
          assetsDir: ".github/assets",
        },
      },
      ...(publish?.notifications
        ? {
            notifications: {
              import: "@allure-notifications/plugin",
              options: publish.notifications,
            },
          }
        : {}),
    },
  });
}
