import { Router } from "express";
import {
  listEngagements,
  getEngagementContext,
} from "../services/engagementReader.js";
import { detectPhase } from "../services/phaseDetector.js";
import {
  getArtifactTree,
  getArtifactContent,
} from "../services/artifactReader.js";
import {
  listChecklists,
  instantiateChecklists,
  toggleChecklistItem,
} from "../services/checklistService.js";

const router = Router();

router.get("/", async (_req, res) => {
  const engagements = await listEngagements();
  res.json({ engagements });
});

router.get("/:slug/context", async (req, res) => {
  const content = await getEngagementContext(req.params.slug);
  if (content === null) {
    res.status(404).json({ error: "Engagement or context not found" });
    return;
  }
  res.json({ slug: req.params.slug, content });
});

router.get("/:slug/phase", async (req, res) => {
  const phaseInfo = await detectPhase(req.params.slug);
  res.json({ slug: req.params.slug, ...phaseInfo });
});

router.get("/:slug/artifacts", async (req, res) => {
  const tree = await getArtifactTree(req.params.slug);
  res.json({ tree });
});

router.get("/:slug/artifacts/*filePath", async (req, res) => {
  const raw = req.params.filePath;
  const filePath = Array.isArray(raw) ? raw.join("/") : String(raw);
  const content = await getArtifactContent(req.params.slug, filePath);
  if (content === null) {
    res.status(404).json({ error: "Artifact not found" });
    return;
  }
  res.json({ path: filePath, content });
});

router.get("/:slug/checklists", async (req, res) => {
  const checklists = await listChecklists(req.params.slug);
  res.json({ checklists });
});

router.post("/:slug/checklists/instantiate", async (req, res) => {
  const copied = await instantiateChecklists(req.params.slug);
  res.json({ ok: true, copied });
});

router.patch("/:slug/checklists/toggle", async (req, res) => {
  const { fileName, lineNumber, checked } = req.body as {
    fileName: string;
    lineNumber: number;
    checked: boolean;
  };
  const success = await toggleChecklistItem(
    req.params.slug,
    fileName,
    lineNumber,
    checked
  );
  if (!success) {
    res.status(400).json({ error: "Failed to toggle item" });
    return;
  }
  res.json({ ok: true });
});

export default router;
