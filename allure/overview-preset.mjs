/**
 * Ethalon overview preset — kit SSOT + profile titles/layers only.
 *
 * Order: Allure + Sonar quality gates, then overview chart quad (indices 2–5).
 * Tiles/gates/renderers come from `@qa-guru/allure-report-kit/presets/overview-preset`.
 */
import { presets } from "@qa-guru/allure-report-kit";
import { OVERVIEW_PRESET as KIT_OVERVIEW_PRESET } from "@qa-guru/allure-report-kit/presets/overview-preset";

import { PYRAMID_LAYERS, TITLES } from "./constants.mjs";
import { buildGatePanels } from "./quality-gate-panels.mjs";

/** @type {import("@qa-guru/allure-report-kit").OverviewPreset} */
export const OVERVIEW_PRESET = {
  ...KIT_OVERVIEW_PRESET,
  titles: {
    currentStatus: TITLES.currentStatus,
    durationDynamics: TITLES.durationDynamics,
    testingPyramid: TITLES.testingPyramid,
    durations: TITLES.durationsByLayer,
  },
  pyramidLayers: [...PYRAMID_LAYERS],
};

/**
 * Overview chart quad only (indices 2–5 when embedded in lead).
 * @param {import("@qa-guru/allure-report-kit").FromOverviewOptions} [options]
 */
export function buildOverviewTiles(options = {}) {
  return presets.fromOverviewCharts({
    preset: OVERVIEW_PRESET,
    layers: [...PYRAMID_LAYERS],
    ...options,
  });
}

/**
 * Lead section: quality gates, then overview charts.
 * @param {import("@qa-guru/allure-report-kit").FromLeadOptions} [options]
 */
export function buildLeadTiles(options = {}) {
  return presets.fromLead({
    preset: OVERVIEW_PRESET,
    layers: [...PYRAMID_LAYERS],
    gatePanels: buildGatePanels({ layout: "2x1" }),
    ...options,
  });
}
