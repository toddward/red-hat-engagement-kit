import { Router, type Response } from "express";
import { v4 as uuidv4 } from "uuid";
import {
  startAssistantQuery,
  type AssistantSession,
} from "../services/assistantService.js";

const router = Router();
const sessions = new Map<string, AssistantSession>();
const pendingQuestions = new Map<string, string>();
// Map engagement slug → active assistant session for reuse
const engagementSessions = new Map<string, string>();

function emitToClients(
  session: AssistantSession,
  event: string,
  data: Record<string, unknown>
) {
  const payload = `event: ${event}\ndata: ${JSON.stringify(data)}\n\n`;
  for (const client of session.sseClients) {
    client.write(payload);
  }
}

// Start a new assistant query or follow up on an existing session
router.post("/:slug/assistant", (req, res) => {
  const { question } = req.body as { question: string };
  if (!question?.trim()) {
    res.status(400).json({ error: "Question is required" });
    return;
  }

  const slug = req.params.slug;

  // Check if there's an existing session for this engagement we can resume
  const existingSessionId = engagementSessions.get(slug);
  const existingSession = existingSessionId
    ? sessions.get(existingSessionId)
    : null;

  if (
    existingSession &&
    existingSession.claudeSessionId &&
    existingSession.status === "completed"
  ) {
    // Resume the existing session with a follow-up
    pendingQuestions.set(existingSession.id, question);
    res.json({ sessionId: existingSession.id, resumed: true });
    return;
  }

  // Create a new session
  const sessionId = uuidv4();
  const session: AssistantSession = {
    id: sessionId,
    slug,
    claudeSessionId: null,
    status: "running",
    sseClients: new Set(),
    abort: () => {},
  };
  sessions.set(sessionId, session);
  engagementSessions.set(slug, sessionId);
  pendingQuestions.set(sessionId, question);

  res.json({ sessionId, resumed: false });
});

router.get("/:slug/assistant/:sessionId/stream", (req, res: Response) => {
  const session = sessions.get(req.params.sessionId);
  if (!session) {
    res.status(404).json({ error: "Session not found" });
    return;
  }

  // If already connected (follow-up), just trigger the query
  if (session.sseClients.has(res)) {
    return;
  }

  res.writeHead(200, {
    "Content-Type": "text/event-stream",
    "Cache-Control": "no-cache",
    Connection: "keep-alive",
  });

  session.sseClients.add(res);

  emitToClients(session, "connected", { sessionId: session.id });

  // Start/resume query
  const question = pendingQuestions.get(session.id);
  if (question) {
    pendingQuestions.delete(session.id);
    const emitSSE = (event: string, data: Record<string, unknown>) =>
      emitToClients(session, event, data);
    startAssistantQuery(session, question, emitSSE);
  }

  req.on("close", () => {
    session.sseClients.delete(res);
  });
});

// Follow-up on an existing connected session
router.post("/:slug/assistant/:sessionId/followup", (req, res) => {
  const { question } = req.body as { question: string };
  const session = sessions.get(req.params.sessionId);

  if (!session) {
    res.status(404).json({ error: "Session not found" });
    return;
  }

  if (!session.claudeSessionId) {
    res.status(400).json({ error: "No Claude session to resume" });
    return;
  }

  const emitSSE = (event: string, data: Record<string, unknown>) =>
    emitToClients(session, event, data);
  startAssistantQuery(session, question, emitSSE);

  res.json({ ok: true });
});

export default router;
