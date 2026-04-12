import fs from "node:fs/promises";
import path from "node:path";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");
const ENGAGEMENTS_DIR = path.join(PROJECT_ROOT, "engagements");
const TEMPLATES_DIR = path.join(PROJECT_ROOT, "knowledge", "checklists");

export interface ChecklistItem {
  text: string;
  checked: boolean;
  line: number;
}

export interface ChecklistSection {
  title: string;
  items: ChecklistItem[];
}

export interface Checklist {
  name: string;
  fileName: string;
  sections: ChecklistSection[];
  completionPercent: number;
}

function parseChecklist(content: string, fileName: string): Checklist {
  const lines = content.split("\n");
  const sections: ChecklistSection[] = [];
  let currentSection: ChecklistSection | null = null;
  let totalItems = 0;
  let checkedItems = 0;

  // Extract name from first H1
  let name = fileName.replace(/\.md$/, "");
  for (const line of lines) {
    if (line.startsWith("# ")) {
      name = line.slice(2).trim();
      break;
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    if (line.startsWith("## ")) {
      currentSection = { title: line.slice(3).trim(), items: [] };
      sections.push(currentSection);
      continue;
    }

    const checkMatch = line.match(/^- \[([ x])\] (.+)$/);
    if (checkMatch && currentSection) {
      const checked = checkMatch[1] === "x";
      currentSection.items.push({
        text: checkMatch[2],
        checked,
        line: i + 1, // 1-based line number
      });
      totalItems++;
      if (checked) checkedItems++;
    }
  }

  return {
    name,
    fileName,
    sections,
    completionPercent: totalItems > 0 ? Math.round((checkedItems / totalItems) * 100) : 0,
  };
}

export async function listChecklists(slug: string): Promise<Checklist[]> {
  const checklistDir = path.join(ENGAGEMENTS_DIR, slug, "checklists");

  let entries;
  try {
    entries = await fs.readdir(checklistDir);
  } catch {
    return [];
  }

  const checklists: Checklist[] = [];
  for (const entry of entries.sort()) {
    if (!entry.endsWith(".md")) continue;
    const content = await fs.readFile(path.join(checklistDir, entry), "utf-8");
    checklists.push(parseChecklist(content, entry));
  }
  return checklists;
}

export async function instantiateChecklists(slug: string): Promise<string[]> {
  const checklistDir = path.join(ENGAGEMENTS_DIR, slug, "checklists");
  await fs.mkdir(checklistDir, { recursive: true });

  let templates;
  try {
    templates = await fs.readdir(TEMPLATES_DIR);
  } catch {
    return [];
  }

  const copied: string[] = [];
  for (const template of templates) {
    if (!template.endsWith(".md")) continue;
    const destPath = path.join(checklistDir, template);
    try {
      await fs.access(destPath);
      // Already exists, skip
    } catch {
      await fs.copyFile(path.join(TEMPLATES_DIR, template), destPath);
      copied.push(template);
    }
  }
  return copied;
}

export async function toggleChecklistItem(
  slug: string,
  fileName: string,
  lineNumber: number,
  checked: boolean
): Promise<boolean> {
  const filePath = path.join(ENGAGEMENTS_DIR, slug, "checklists", fileName);

  // Security: validate fileName doesn't contain path traversal
  if (fileName.includes("/") || fileName.includes("..")) {
    return false;
  }

  let content;
  try {
    content = await fs.readFile(filePath, "utf-8");
  } catch {
    return false;
  }

  const lines = content.split("\n");
  const idx = lineNumber - 1;
  if (idx < 0 || idx >= lines.length) return false;

  const line = lines[idx];
  if (checked) {
    lines[idx] = line.replace("- [ ] ", "- [x] ");
  } else {
    lines[idx] = line.replace("- [x] ", "- [ ] ");
  }

  if (lines[idx] === line) return false; // No change made

  await fs.writeFile(filePath, lines.join("\n"), "utf-8");
  return true;
}
