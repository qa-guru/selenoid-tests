#!/usr/bin/env node
/**
 * Convert JUnit XML (gotestsum --junitfile) to Allure 2 JSON results.
 * Uses @allurereport/reader junitXml reader (Allure 3 JUnit plugin).
 */
import { createHash, randomUUID } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { basename, resolve } from "node:path";
import { junitXml } from "@allurereport/reader";
import { BufferResultFile } from "@allurereport/reader-api";

function parseArgs(argv) {
  const args = { input: "", output: "", epic: "", layer: "unit" };
  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--input" || arg === "-i") args.input = argv[++i];
    else if (arg === "--output" || arg === "-o") args.output = argv[++i];
    else if (arg === "--epic") args.epic = argv[++i];
    else if (arg === "--layer") args.layer = argv[++i];
  }
  if (!args.input || !args.output || !args.epic) {
    console.error("Usage: junit-to-allure.mjs --input FILE.xml --output DIR --epic NAME [--layer unit]");
    process.exit(2);
  }
  return args;
}

function buildHistoryId(fullName, name) {
  return createHash("md5").update(fullName || name || randomUUID()).digest("hex");
}

async function convert({ input, output, epic, layer }) {
  const outputDir = resolve(output);
  await mkdir(outputDir, { recursive: true });

  const attachmentIndex = new Map();
  let count = 0;

  const visitor = {
    async visitTestResult(result) {
      const uuid = randomUUID();
      const name = result.name ?? "unnamed test";
      const fullName = result.fullName ?? name;
      const labels = [
        ...(result.labels ?? []),
        { name: "epic", value: epic },
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
  console.log(`Wrote ${count} Allure results to ${outputDir} (epic=${epic}, layer=${layer})`);
}

const args = parseArgs(process.argv.slice(2));
await convert(args);
