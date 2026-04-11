import { Router } from "express";
import {
  listEngagements,
  getEngagementContext,
} from "../services/engagementReader.js";

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

export default router;
