#!/usr/bin/env node
/**
 * Convert JUnit XML (gotestsum --junitfile) to Allure 2 JSON results.
 * Uses @allurereport/reader junitXml reader (Allure 3 JUnit plugin).
 */
import { createHash, randomUUID } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { basename, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { junitXml } from "@allurereport/reader";
import { BufferResultFile } from "@allurereport/reader-api";

function parseArgs(argv) {
  const args = { input: "", output: "", epic: "", component: "", layer: "unit" };
  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--input" || arg === "-i") args.input = argv[++i];
    else if (arg === "--output" || arg === "-o") args.output = argv[++i];
    else if (arg === "--epic") args.epic = argv[++i];
    else if (arg === "--component") args.component = argv[++i];
    else if (arg === "--layer") args.layer = argv[++i];
  }
  if (!args.input || !args.output || !args.epic) {
    console.error(
      "Usage: junit-to-allure.mjs --input FILE.xml --output DIR --epic NAME [--component NAME] [--layer unit]",
    );
    process.exit(2);
  }
  return args;
}

function buildHistoryId(fullName, name) {
  return createHash("md5").update(fullName || name || randomUUID()).digest("hex");
}

/** Go t.Run subtest name after "TestFunc/" — readable title in Allure/TestOps. */
export function resolveDisplayName(name, fullName) {
  if (name) {
    const slash = name.indexOf("/");
    if (slash >= 0 && slash < name.length - 1) {
      return name.slice(slash + 1).trim();
    }
  }
  if (fullName) {
    const match = fullName.match(/\.Test[^/]+\/(.*)$/s);
    if (match?.[1]) {
      return match[1].trim();
    }
  }
  return name;
}

async function convert({ input, output, epic, component, layer }) {
  const componentLabel = component || epic;
  const outputDir = resolve(output);
  await mkdir(outputDir, { recursive: true });

  const attachmentIndex = new Map();
  let count = 0;

  const visitor = {
    async visitTestResult(result) {
      const uuid = randomUUID();
      const rawName = result.name ?? "unnamed test";
      const fullName = result.fullName ?? rawName;
      const name = resolveDisplayName(rawName, fullName);
      const labels = [
        ...(result.labels ?? []),
        { name: "epic", value: epic },
        { name: "component", value: componentLabel },
        { name: "layer", value: layer },
        { name: "framework", value: "go" },
        { name: "language", value: "go" },
      ];

      const steps = (result.steps ?? []).map((step) => {
        if (step.type !== "attachment") return step;
        const stored = attachmentIndex.get(step.originalFileName);
        if (!stored) return step;
        return {
          ...step,
          source: stored.fileName,
        };
      });

      const payload = {
        uuid,
        historyId: buildHistoryId(fullName, name),
        name,
        fullName,
        status: result.status ?? "unknown",
        stage: "finished",
        description: "",
        steps,
        attachments: [],
        labels,
        parameters: [],
        links: [],
      };

      if (result.duration != null) {
        const stop = Date.now();
        payload.stop = stop;
        payload.start = stop - result.duration;
      }
      if (result.message || result.trace) {
        payload.statusDetails = {
          message: result.message,
          trace: result.trace,
        };
      }

      await writeFile(`${outputDir}/${uuid}-result.json`, JSON.stringify(payload));
      count += 1;
    },
    async visitCheckResult() {},
    async visitTestFixtureResult() {},
    async visitAttachmentFile(file) {
      const original = file.getOriginalFileName();
      const buffer = await file.asBuffer();
      if (!buffer) return;
      const fileName = `${original}-attachment.txt`;
      await writeFile(`${outputDir}/${fileName}`, buffer);
      attachmentIndex.set(original, { fileName });
    },
    async visitMetadata() {},
    async visitGlobals() {},
  };

  const xmlPath = resolve(input);
  const xml = await readFile(xmlPath);
  const ok = await junitXml.read(visitor, new BufferResultFile(xml, basename(xmlPath)));
  if (!ok) {
    throw new Error(`Failed to parse JUnit XML: ${xmlPath}`);
  }
  console.log(
    `Wrote ${count} Allure results to ${outputDir} (epic=${epic}, component=${componentLabel}, layer=${layer})`,
  );
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  const args = parseArgs(process.argv.slice(2));
  await convert(args);
}
