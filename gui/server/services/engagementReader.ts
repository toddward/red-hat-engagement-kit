import fs from "node:fs/promises";
import path from "node:path";
import type { EngagementInfo } from "../types/index.js";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");
const ENGAGEMENTS_DIR = path.join(PROJECT_ROOT, "engagements");

export async function listEngagements(): Promise<EngagementInfo[]> {
  let entries;
  try {
    entries = await fs.readdir(ENGAGEMENTS_DIR, { withFileTypes: true });
  } catch {
    return [];
  }

  const engagements: EngagementInfo[] = [];

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    if (entry.name.startsWith(".")) continue; // skip .template

    const contextPath = path.join(ENGAGEMENTS_DIR, entry.name, "CONTEXT.md");
    let hasContext = false;
    try {
      await fs.access(contextPath);
      hasContext = true;
    } catch {
      // no CONTEXT.md
    }

    engagements.push({
      slug: entry.name,
      hasContext,
    });
  }

  return engagements.sort((a, b) => a.slug.localeCompare(b.slug));
}

export async function getEngagementContext(
  slug: string
): Promise<string | null> {
  const contextPath = path.join(ENGAGEMENTS_DIR, slug, "CONTEXT.md");
  try {
    return await fs.readFile(contextPath, "utf-8");
  } catch {
    return null;
  }
}
