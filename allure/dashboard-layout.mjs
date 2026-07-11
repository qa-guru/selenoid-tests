import {
  PYRAMID_LAYERS,
  STABILITY_SKIP_STATUSES,
  STABILITY_THRESHOLD,
  TITLES,
} from "./constants.mjs";

/** Filter tests by Allure `component` label (hub README dashboard crops). */
export function componentLabelFilter(component) {
  return ({ labels }) =>
    labels.some(({ name, value }) => name === "component" && value === component);
}

/**
 * Compact 2×2 dashboard for README PNG (status + pyramid + dynamics + treemap).
 * Per-widget filter — no sibling components (webdriver vs playwright).
 */
export function buildComponentReadmeDashboardLayout(component) {
  const filter = componentLabelFilter(component);

  return [
    {
      type: "currentStatus",
      title: TITLES.currentStatus,
      filter,
    },
    {
      type: "testingPyramid",
      title: TITLES.testingPyramid,
      layers: [...PYRAMID_LAYERS],
      filter,
    },
    {
      type: "statusDynamics",
      title: TITLES.statusDynamics,
      limit: 20,
      filter,
    },
    {
      type: "successRateDistribution",
      title: TITLES.successRateDistribution,
      filter,
    },
  ];
}

/**
 * Dashboard plugin layout.
 * Invariant: index 0 = currentStatus, index 1 = testingPyramid.
 */
export function buildDashboardLayout({ epicCharts = [] } = {}) {
  const epicStatusDynamics = epicCharts.map((epic) => ({
    type: "statusDynamics",
    title: `Динамика — ${epic}`,
    limit: 20,
    filter: ({ labels }) =>
      labels.some(({ name, value }) => name === "epic" && value === epic),
  }));

  return [
    {
      type: "currentStatus",
      title: TITLES.currentStatus,
    },
    {
      type: "testingPyramid",
      title: TITLES.testingPyramid,
      layers: [...PYRAMID_LAYERS],
    },
    {
      type: "statusDynamics",
      title: TITLES.statusDynamics,
      limit: 20,
    },
    {
      type: "successRateDistribution",
      title: TITLES.successRateDistribution,
    },
    {
      type: "stabilityDistribution",
      title: TITLES.stabilityByComponent,
      threshold: STABILITY_THRESHOLD,
      skipStatuses: [...STABILITY_SKIP_STATUSES],
      groupBy: "label-name:component",
    },
    {
      type: "coverageDiff",
      title: TITLES.coverageDiff,
    },
    {
      type: "statusTransitions",
      title: TITLES.statusTransitions,
      limit: 20,
    },
    {
      type: "testBaseGrowthDynamics",
      title: TITLES.testBaseGrowthDynamics,
      limit: 20,
    },
    {
      type: "durationDynamics",
      title: TITLES.durationDynamics,
      limit: 20,
    },
    {
      type: "problemsDistribution",
      title: TITLES.problemsByEnvironment,
      by: "environment",
    },
    {
      type: "stabilityDistribution",
      title: TITLES.stabilityByFeature,
      threshold: STABILITY_THRESHOLD,
      skipStatuses: [...STABILITY_SKIP_STATUSES],
      groupBy: "feature",
    },
    {
      type: "stabilityDistribution",
      title: TITLES.stabilityByEpic,
      threshold: STABILITY_THRESHOLD,
      skipStatuses: [...STABILITY_SKIP_STATUSES],
      groupBy: "epic",
    },
    {
      type: "stabilityDistribution",
      title: TITLES.stabilityByStory,
      threshold: STABILITY_THRESHOLD,
      skipStatuses: [...STABILITY_SKIP_STATUSES],
      groupBy: "story",
    },
    {
      type: "testResultSeverities",
      title: TITLES.testResultSeverities,
    },
    {
      type: "durations",
      title: TITLES.durations,
      groupBy: "none",
    },
    {
      type: "durations",
      title: TITLES.durationsByLayer,
      groupBy: "layer",
    },
    {
      type: "statusAgePyramid",
      title: TITLES.statusAgePyramid,
      limit: 20,
    },
    ...epicStatusDynamics,
  ];
}
