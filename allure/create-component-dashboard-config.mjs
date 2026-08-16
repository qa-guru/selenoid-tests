import { withKit, theme, renderers } from "@qa-guru/allure-report-kit";

import { REPORT_LANGUAGE } from "./constants.mjs";
import { buildComponentReadmeDashboardLayout } from "./dashboard-layout.mjs";

/**
 * Minimal Allure config: dashboard plugin only.
 * Results are sliced by `component` before generate (see filter-allure-results-by-component.mjs).
 *
 * Env: ALLURE_COMPONENT_DASHBOARD (required) — e.g. playwright-image
 */
export function createComponentDashboardConfig({
  slug,
  component = process.env.ALLURE_COMPONENT_DASHBOARD,
} = {}) {
  if (!slug || typeof slug !== "string") {
    throw new Error("createComponentDashboardConfig: slug is required");
  }
  if (!component || typeof component !== "string") {
    throw new Error(
      "createComponentDashboardConfig: component or ALLURE_COMPONENT_DASHBOARD is required",
    );
  }

  return withKit({
    softFork: true,
    renderer: renderers.stock(),
    theme: {
      ...theme.qaGuru(),
      header: {
        enabled: true,
        source: "design-system",
        productName: `${slug} — ${component}`,
      },
    },
    name: `${slug} — ${component}`,
    plugins: {
      dashboard: {
        options: {
          reportName: component,
          reportLanguage: REPORT_LANGUAGE,
          layout: buildComponentReadmeDashboardLayout(component),
        },
      },
    },
  });
}
