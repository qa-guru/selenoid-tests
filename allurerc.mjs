import { createAllureConfig } from "./allure/create-config.mjs";

export default createAllureConfig({
  slug: "selenoid-tests",
  epicCharts: [
    "selenoid",
    "selenoid-ui",
    "cm",
    "playwright-image",
    "webdriver-image",
  ],
});
