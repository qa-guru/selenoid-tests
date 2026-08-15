/**
 * Lead Allure + Sonar quality-gate overrides for presets.fromLead.
 * Live CI JSON (`allure/sonar-quality-gate.json` or SONAR_QUALITY_GATE_JSON)
 * wins; otherwise the SSOT fixture.
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { sonarProjectStatusToQualityGateOptions } from "@qa-guru/allure-report-kit/runtime";

import {
  QUALITY_GATE_LABELS,
  REPORT_LANGUAGE,
  SONAR_QUALITY_GATE_FIXTURE,
  SONAR_QUALITY_GATE_LABELS,
  SONAR_QUALITY_GATE_PROFILE_CONDITIONS,
  SONAR_QUALITY_GATE_SOURCE,
  TITLES,
} from "./constants.mjs";

const LIVE_SONAR_GATE_JSON = path.join(
  path.dirname(fileURLToPath(import.meta.url)),
  "sonar-quality-gate.json",
);

function loadSonarProjectStatus() {
  const filePath = process.env.SONAR_QUALITY_GATE_JSON || LIVE_SONAR_GATE_JSON;
  try {
    if (!fs.existsSync(filePath)) {
      return SONAR_QUALITY_GATE_FIXTURE;
    }
    const parsed = JSON.parse(fs.readFileSync(filePath, "utf8"));
    if (parsed && typeof parsed === "object" && (parsed.status || Array.isArray(parsed.conditions))) {
      return parsed;
    }
  } catch {
    /* keep fixture */
  }
  return SONAR_QUALITY_GATE_FIXTURE;
}

function buildSonarGateData() {
  const projectStatus = loadSonarProjectStatus();
  const projectKey =
    projectStatus.project_key || projectStatus.projectKey || SONAR_QUALITY_GATE_SOURCE.projectKey;
  return sonarProjectStatusToQualityGateOptions(projectStatus, {
    lang: REPORT_LANGUAGE,
    profile: SONAR_QUALITY_GATE_SOURCE.profile,
    profileConditions: SONAR_QUALITY_GATE_PROFILE_CONDITIONS.map((row) => ({ ...row })),
    source: { ...SONAR_QUALITY_GATE_SOURCE, projectKey },
    labels: SONAR_QUALITY_GATE_LABELS,
    barTitle: TITLES.sonarQualityGate,
  });
}

/**
 * @param {{ layout?: string }} [options]
 * @returns {NonNullable<import("@qa-guru/allure-report-kit").FromLeadOptions["gatePanels"]>}
 */
export function buildGatePanels({ layout = "2x1" } = {}) {
  return {
    allureQualityGate: {
      title: TITLES.qualityGate,
      layout,
      labels: QUALITY_GATE_LABELS,
    },
    sonarQualityGate: {
      title: TITLES.sonarQualityGate,
      layout,
      dots: false,
      data: buildSonarGateData(),
    },
  };
}
