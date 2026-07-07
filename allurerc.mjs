/** @type {import('@allurereport/core-api').Config} */
const EPICS = ["selenoid", "selenoid-ui", "cm", "playwright-image", "webdriver-image"];

const epicStatusDynamics = EPICS.map((epic) => ({
  type: "statusDynamics",
  title: `Динамика — ${epic}`,
  limit: 30,
  filter: ({ labels }) => labels.some(({ name, value }) => name === "epic" && value === epic),
}));

export default {
  name: "selenoid-tests Tests",
  historyPath: "./history.jsonl",
  appendHistory: true,
  historyLimit: 20,
  knownIssuesPath: "./known.json",
  variables: {
    Framework: "JUnit 5 + Selenide + Go",
    Report: "Allure 3",
  },
  qualityGate: {
    rules: [{ maxFailures: 0 }],
  },
  categories: {
    rules: [
      {
        name: "Таймауты",
        matchers: {
          statuses: ["failed", "broken"],
          message: "timeout|Timeout|timed out",
        },
        groupBy: ["severity"],
        groupByMessage: true,
      },
      {
        name: "Ошибки проверок",
        matchers: {
          statuses: ["failed", "broken"],
          message: "AssertionError|assertion|expected",
        },
        groupBy: ["severity"],
        groupByMessage: true,
      },
    ],
  },
  plugins: {
    awesome: {
      options: {
        reportLanguage: "ru",
        groupBy: ["parentSuite", "suite", "subSuite"],
        charts: [
          { type: "currentStatus", title: "Текущий статус" },
          { type: "testResultSeverities", title: "Результаты по severity" },
          { type: "statusDynamics", title: "Динамика статусов", limit: 20 },
          { type: "statusTransitions", title: "Переходы статусов", limit: 20 },
          { type: "durationDynamics", title: "Динамика длительности", limit: 20 },
          { type: "durations", title: "Гистограмма длительностей", groupBy: "none" },
        ],
      },
    },
    dashboard: {
      options: {
        reportName: "selenoid-tests Tests Dashboard",
        reportLanguage: "ru",
        layout: [
          {
            type: "currentStatus",
            title: "Текущий статус по сервисам",
            groupBy: "epic",
          },
          {
            type: "testingPyramid",
            title: "Пирамида тестирования",
            layers: ["unit", "component", "integration", "api", "e2e", "manual"],
          },
          {
            type: "stabilityDistribution",
            title: "Стабильность по сервисам",
            threshold: 90,
            skipStatuses: ["skipped", "unknown"],
            groupBy: "epic",
          },
          {
            type: "successRateDistribution",
            title: "Распределение успешности",
          },
          {
            type: "statusDynamics",
            title: "Динамика статусов (все сервисы)",
            limit: 30,
          },
          {
            type: "statusTransitions",
            title: "Переходы статусов",
            limit: 20,
          },
          {
            type: "durationDynamics",
            title: "Динамика длительности",
            limit: 20,
          },
          ...epicStatusDynamics,
          {
            type: "currentStatus",
            title: "Текущий статус",
          },
          {
            type: "stabilityDistribution",
            title: "Стабильность по feature",
            threshold: 90,
            skipStatuses: ["skipped", "unknown"],
            groupBy: "feature",
          },
        ],
      },
    },
    csv: {
      options: {
        fileName: "selenoid-tests.csv",
      },
    },
  },
};
