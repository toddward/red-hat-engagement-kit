import express from "express";
import cors from "cors";
import skillsRouter from "./routes/skills.js";
import engagementsRouter from "./routes/engagements.js";
import executionRouter from "./routes/execution.js";
import assistantRouter from "./routes/assistant.js";

const app = express();
const PORT = 3001;

app.use(cors());
app.use(express.json());

app.use("/api/skills", skillsRouter);
app.use("/api/engagements", engagementsRouter);
app.use("/api/execute", executionRouter);
app.use("/api/engagements", assistantRouter);

app.listen(PORT, () => {
  console.log(`[server] running on http://localhost:${PORT}`);
});
