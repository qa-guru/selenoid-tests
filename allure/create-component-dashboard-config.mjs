import { REPORT_LANGUAGE, HISTORY_DEFAULTS } from "./constants.mjs";
import { buildComponentReadmeDashboardLayout } from "./dashboard-layout.mjs";

/**
 * Minimal Allure config: dashboard plugin only, filtered by `component` label.
 * Used for README PNG crops (playwright-image / webdriver-image).
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

  return {
    name: `${slug} — ${component}`,
    ...HISTORY_DEFAULTS,
    plugins: {
      dashboard: {
        options: {
          reportName: component,
          reportLanguage: REPORT_LANGUAGE,
          layout: buildComponentReadmeDashboardLayout(component),
        },
      },
    },
  };
}
