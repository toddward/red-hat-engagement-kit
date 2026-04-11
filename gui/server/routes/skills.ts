import { Router } from "express";
import { listSkills, getSkill } from "../services/skillReader.js";

const router = Router();

router.get("/", async (_req, res) => {
  const skills = await listSkills();
  res.json({
    skills: skills.map(({ name, description, path }) => ({
      name,
      description,
      path,
    })),
  });
});

router.get("/:name", async (req, res) => {
  const skill = await getSkill(req.params.name);
  if (!skill) {
    res.status(404).json({ error: "Skill not found" });
    return;
  }
  res.json({
    name: skill.name,
    description: skill.description,
    content: skill.content,
  });
});

export default router;
