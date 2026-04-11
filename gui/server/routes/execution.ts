import { Router } from "express";
import type { Response } from "express";
import { v4 as uuidv4 } from "uuid";
import type { ExecutionSession } from "../providers/types.js";
import { startClaudeSession } from "../providers/claudeProvider.js";

const router = Router();
const sessions = new Map<string, ExecutionSession>();

function emitToClients(
  session: ExecutionSession,
  event: string,
  data: Record<string, unknown>
) {
  const payload = `event: ${event}\ndata: ${JSON.stringify(data)}\n\n`;
  for (const client of session.sseClients) {
    client.write(payload);
  }
}

// Start a new execution
router.post("/", (req, res) => {
  const { skills, engagementSlug, provider } = req.body as {
    skills: string[];
    engagementSlug?: string;
    provider?: string;
  };

  if (!skills || skills.length === 0) {
    res.status(400).json({ error: "At least one skill is required" });
    return;
  }

  const executionId = uuidv4();
  const session: ExecutionSession = {
    id: executionId,
    skills,
    engagementSlug: engagementSlug ?? null,
    provider: provider ?? "claude",
    status: "running",
    sseClients: new Set(),
    pendingResponses: new Map(),
    abort: () => {},
  };
  sessions.set(executionId, session);

  // Defer execution start until SSE client connects
  res.json({ executionId });
});

// SSE stream for execution output
router.get("/:id/stream", (req, res: Response) => {
  const session = sessions.get(req.params.id);
  if (!session) {
    res.status(404).json({ error: "Execution not found" });
    return;
  }

  res.writeHead(200, {
    "Content-Type": "text/event-stream",
    "Cache-Control": "no-cache",
    Connection: "keep-alive",
  });

  session.sseClients.add(res);

  // Send initial connection event
  emitToClients(session, "connected", {
    executionId: session.id,
    skills: session.skills,
    provider: session.provider,
  });

  // Start execution when first client connects
  if (session.sseClients.size === 1) {
    const emitSSE = (event: string, data: Record<string, unknown>) =>
      emitToClients(session, event, data);

    if (session.provider === "claude") {
      startClaudeSession(session, emitSSE);
    } else {
      emitSSE("error", { text: `Unknown provider: ${session.provider}` });
    }
  }

  req.on("close", () => {
    session.sseClients.delete(res);
  });
});

// Respond to a question or approval during execution
router.post("/:id/respond", (req, res) => {
  const session = sessions.get(req.params.id);
  if (!session) {
    res.status(404).json({ error: "Execution not found" });
    return;
  }

  const { type, questionId, approvalId, answers, allow, message } = req.body;

  if (type === "answer" && questionId) {
    const pending = session.pendingResponses.get(questionId);
    if (pending) {
      pending.resolve(answers);
      session.pendingResponses.delete(questionId);
    }
  } else if (type === "approval" && approvalId) {
    const pending = session.pendingResponses.get(approvalId);
    if (pending) {
      pending.resolve(allow);
      session.pendingResponses.delete(approvalId);
    }
  } else if (type === "followup" && message) {
    if (session.sendFollowUp) {
      session.sendFollowUp(message);
    }
  }

  res.json({ ok: true });
});

export default router;
