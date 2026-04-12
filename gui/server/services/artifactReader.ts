import fs from "node:fs/promises";
import path from "node:path";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");
const ENGAGEMENTS_DIR = path.join(PROJECT_ROOT, "engagements");

export interface ArtifactNode {
  name: string;
  path: string;
  type: "file" | "directory";
  children?: ArtifactNode[];
}

async function buildTree(
  dirPath: string,
  relativeTo: string,
  depth: number = 0
): Promise<ArtifactNode[]> {
  if (depth > 3) return [];

  let entries;
  try {
    entries = await fs.readdir(dirPath, { withFileTypes: true });
  } catch {
    return [];
  }

  const nodes: ArtifactNode[] = [];

  for (const entry of entries.sort((a, b) => a.name.localeCompare(b.name))) {
    if (entry.name.startsWith(".") || entry.name === ".gitkeep") continue;

    const fullPath = path.join(dirPath, entry.name);
    const relPath = path.relative(relativeTo, fullPath);

    if (entry.isDirectory()) {
      const children = await buildTree(fullPath, relativeTo, depth + 1);
      nodes.push({
        name: entry.name,
        path: relPath,
        type: "directory",
        children,
      });
    } else {
      nodes.push({
        name: entry.name,
        path: relPath,
        type: "file",
      });
    }
  }

  // Sort: directories first, then files
  return nodes.sort((a, b) => {
    if (a.type !== b.type) return a.type === "directory" ? -1 : 1;
    return a.name.localeCompare(b.name);
  });
}

export async function getArtifactTree(slug: string): Promise<ArtifactNode[]> {
  const engDir = path.join(ENGAGEMENTS_DIR, slug);
  return buildTree(engDir, engDir);
}

export async function getArtifactContent(
  slug: string,
  filePath: string
): Promise<string | null> {
  const engDir = path.join(ENGAGEMENTS_DIR, slug);
  const resolved = path.resolve(engDir, filePath);

  // Security: prevent directory traversal
  if (!resolved.startsWith(engDir)) {
    return null;
  }

  // Only serve text-based files
  const ext = path.extname(resolved).toLowerCase();
  const allowed = [".md", ".txt", ".json", ".yaml", ".yml", ".csv", ".html"];
  if (!allowed.includes(ext)) {
    return null;
  }

  try {
    return await fs.readFile(resolved, "utf-8");
  } catch {
    return null;
  }
}
