import fs from "node:fs/promises";
import path from "node:path";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");
const ENGAGEMENTS_DIR = path.join(PROJECT_ROOT, "engagements");

export type EngagementPhase = "pre-engagement" | "live" | "leave-behind";

export interface PhaseInfo {
  phase: EngagementPhase;
  hasContext: boolean;
  hasDeliverables: boolean;
  artifactCounts: {
    discovery: number;
    assessments: number;
    deliverables: number;
    checklists: number;
  };
}

async function countFiles(dirPath: string): Promise<number> {
  try {
    const entries = await fs.readdir(dirPath);
    return entries.filter((e) => !e.startsWith(".") && e !== ".gitkeep").length;
  } catch {
    return 0;
  }
}

export async function detectPhase(slug: string): Promise<PhaseInfo> {
  const engDir = path.join(ENGAGEMENTS_DIR, slug);
  const contextPath = path.join(engDir, "CONTEXT.md");

  let hasContext = false;
  try {
    await fs.access(contextPath);
    hasContext = true;
  } catch {
    // no CONTEXT.md
  }

  const [discovery, assessments, deliverables, checklists] = await Promise.all([
    countFiles(path.join(engDir, "discovery")),
    countFiles(path.join(engDir, "assessments")),
    countFiles(path.join(engDir, "deliverables")),
    countFiles(path.join(engDir, "checklists")),
  ]);

  const hasDeliverables = deliverables > 0;

  let phase: EngagementPhase;
  if (!hasContext) {
    phase = "pre-engagement";
  } else if (hasDeliverables) {
    phase = "leave-behind";
  } else {
    phase = "live";
  }

  return {
    phase,
    hasContext,
    hasDeliverables,
    artifactCounts: { discovery, assessments, deliverables, checklists },
  };
}
