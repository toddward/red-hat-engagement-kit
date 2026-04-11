import fs from "node:fs/promises";
import path from "node:path";
import type { SkillInfo } from "../types/index.js";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");
const SKILLS_DIR = path.join(PROJECT_ROOT, ".claude", "skills");

function extractDescription(content: string): string {
  const lines = content.split("\n");
  let inFrontmatter = false;
  let passedFrontmatter = false;

  for (const line of lines) {
    const trimmed = line.trim();
    if (trimmed === "---") {
      if (!passedFrontmatter) {
        inFrontmatter = !inFrontmatter;
        if (!inFrontmatter) passedFrontmatter = true;
        continue;
      }
      continue;
    }
    if (inFrontmatter) continue;
    if (
      trimmed &&
      !trimmed.startsWith("#") &&
      !trimmed.startsWith("```")
    ) {
      return trimmed.length > 120 ? trimmed.slice(0, 120) + "..." : trimmed;
    }
  }
  return "";
}

export async function listSkills(): Promise<SkillInfo[]> {
  const entries = await fs.readdir(SKILLS_DIR, { withFileTypes: true });
  const skills: SkillInfo[] = [];

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;

    const skillMdPath = path.join(SKILLS_DIR, entry.name, "SKILL.md");
    try {
      const content = await fs.readFile(skillMdPath, "utf-8");
      skills.push({
        name: entry.name,
        description: extractDescription(content),
        path: path.relative(PROJECT_ROOT, skillMdPath),
        content,
      });
    } catch {
      // Skip directories without SKILL.md
    }
  }

  return skills.sort((a, b) => a.name.localeCompare(b.name));
}

export async function getSkill(name: string): Promise<SkillInfo | null> {
  const skills = await listSkills();
  return skills.find((s) => s.name === name) ?? null;
}
