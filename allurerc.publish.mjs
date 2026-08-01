import { createAllureConfig } from "./allure/create-config.mjs";

const mode = process.env.NOTIFICATION_MODE || "dry-run";

export default createAllureConfig({
  slug: process.env.ALLURE_PUBLISH_SLUG || "reference-app",
  publish: {
    notifications: {
      config: "./notifications/config.runtime.json",
      mode,
      allureFolder: "./build/reports/allure-report/allureReport/awesome",
      allureResultsFolder: "./build/allure-results",
    },
  },
});
