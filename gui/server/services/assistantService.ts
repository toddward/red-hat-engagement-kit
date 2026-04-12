import { spawn } from "node:child_process";
import path from "node:path";
import type { Response } from "express";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");

export interface AssistantSession {
  id: string;
  slug: string;
  claudeSessionId: string | null;
  status: "running" | "completed" | "error";
  sseClients: Set<Response>;
  abort: () => void;
}

export function startAssistantQuery(
  session: AssistantSession,
  question: string,
  emitSSE: (event: string, data: Record<string, unknown>) => void
): void {
  const isFollowUp = session.claudeSessionId !== null;

  const prompt = isFollowUp
    ? question
    : `Read all files in engagements/${session.slug}/ including CONTEXT.md, everything in discovery/, assessments/, and deliverables/. Then answer this question based only on those engagement artifacts:\n\n${question}`;

  const args = [
    "--print",
    "--output-format",
    "stream-json",
    "--verbose",
    "--model",
    "sonnet",
    "--permission-mode",
    "acceptEdits",
    "--allowedTools",
    "Read",
    "Glob",
    "Grep",
  ];

  if (isFollowUp) {
    args.push("--resume", session.claudeSessionId!);
  }

  args.push("--", prompt);

  console.log(
    `[assistant] ${isFollowUp ? "resuming" : "starting"} query for ${session.slug}`
  );

  const proc = spawn("claude", args, {
    cwd: PROJECT_ROOT,
    env: { ...process.env },
    stdio: ["pipe", "pipe", "pipe"],
  });

  session.status = "running";

  session.abort = () => {
    proc.kill("SIGTERM");
    session.status = "completed";
  };

  let buffer = "";

  proc.stdout?.on("data", (chunk: Buffer) => {
    buffer += chunk.toString();
    const lines = buffer.split("\n");
    buffer = lines.pop() ?? "";

    for (const line of lines) {
      if (!line.trim()) continue;
      try {
        const msg = JSON.parse(line);
        handleMessage(msg, session, emitSSE);
      } catch {
        emitSSE("assistant", { text: line });
      }
    }
  });

  proc.stderr?.on("data", (chunk: Buffer) => {
    const text = chunk.toString().trim();
    if (text) console.error(`[assistant stderr] ${text}`);
  });

  proc.on("close", (code) => {
    if (buffer.trim()) {
      try {
        const msg = JSON.parse(buffer);
        handleMessage(msg, session, emitSSE);
      } catch {
        if (buffer.trim()) emitSSE("assistant", { text: buffer });
      }
    }

    if (session.status === "running") {
      session.status = "completed";
      emitSSE("result", { status: code === 0 ? "success" : "error" });
    }
  });

  proc.on("error", (err) => {
    session.status = "error";
    emitSSE("error", { text: `Assistant error: ${err.message}` });
  });
}

function handleMessage(
  msg: Record<string, unknown>,
  session: AssistantSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void
): void {
  const type = msg.type as string;

  if (type === "system") {
    // Capture the Claude session ID for resume
    const sessionId = msg.session_id as string;
    if (sessionId && !session.claudeSessionId) {
      session.claudeSessionId = sessionId;
      console.log(`[assistant] claude session ID: ${sessionId}`);
    }
  } else if (type === "assistant") {
    const message = msg.message as {
      content?: Array<{
        type: string;
        text?: string;
        name?: string;
        input?: unknown;
      }>;
    };
    if (message?.content) {
      for (const block of message.content) {
        if (block.type === "text" && block.text) {
          emitSSE("assistant", { text: block.text });
        } else if (block.type === "tool_use") {
          emitSSE("tool_use", {
            tool: block.name,
            input:
              typeof block.input === "string"
                ? block.input
                : JSON.stringify(block.input, null, 2),
          });
        }
      }
    }
  } else if (type === "result" && msg.total_cost_usd != null) {
    session.status = "completed";
    emitSSE("result", {
      status: (msg.subtype as string) ?? "success",
      cost: msg.total_cost_usd as number,
    });
  }
}
