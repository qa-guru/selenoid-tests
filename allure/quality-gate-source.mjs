/**
 * Quality-gate popover paths for qa-guru/selenoid-tests (standalone clone).
 * Objects are complete — do not partial-copy in shell.
 */
export const ALLURE_QUALITY_GATE_SOURCE = {
  configFile: "allurerc.mjs",
  rulesFile: "allure/quality-gate.mjs",
  knownIssuesFile: "./known.json",
  hrefBase: "https://github.com/qa-guru/selenoid-tests/blob/main/",
};

export const ALLURE_QUALITY_GATE_LABELS = {
  passed: { ru: "Allure Quality Gate пройден", en: "Allure Quality Gate passed" },
  failed: { ru: "Allure Quality Gate не пройден", en: "Allure Quality Gate failed" },
};

export const SONAR_QUALITY_GATE_SOURCE = {
  configFile: "sonar-project.properties",
  profile: "qa-guru-canon",
  projectKey: "selenoid-tests",
  hrefBase: "https://github.com/qa-guru/selenoid-tests/blob/main/",
};

export const SONAR_QUALITY_GATE_LABELS = {
  passed: { ru: "Sonar Quality Gate пройден", en: "Sonar Quality Gate passed" },
  failed: { ru: "Sonar Quality Gate не пройден", en: "Sonar Quality Gate failed" },
};

/** Fallback when CI has not attached allure/sonar-quality-gate.json. */
export const SONAR_QUALITY_GATE_FIXTURE = {
  status: "OK",
  project_key: "selenoid-tests",
  analysis_id: "AXdemoPassedAnalysis",
  dashboard_url: "https://sonar.qa.guru/dashboard?id=selenoid-tests",
  conditions: [
    {
      status: "OK",
      metricKey: "coverage",
      comparator: "LT",
      errorThreshold: 80,
      actualValue: 100,
    },
    {
      status: "OK",
      metricKey: "bugs",
      comparator: "GT",
      errorThreshold: 0,
      actualValue: 0,
    },
  ],
};

export const SONAR_QUALITY_GATE_PROFILE_CONDITIONS = [
  { metric: "coverage", op: "LT", error: 80, label: "Coverage on Overall Code ≥ 80%" },
  { metric: "bugs", op: "GT", error: 0, label: "Bugs on Overall Code = 0" },
];
