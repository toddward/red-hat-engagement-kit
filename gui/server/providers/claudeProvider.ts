import { spawn, type ChildProcess } from "node:child_process";
import path from "node:path";
import type { ExecutionSession } from "./types.js";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");

function getSessionMeta(session: ExecutionSession) {
  return session as unknown as Record<string, unknown>;
}

export function startClaudeSession(
  session: ExecutionSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void
): void {
  const skillList = session.skills.map((s) => `/${s}`).join(", then ");
  const engagementCtx = session.engagementSlug
    ? ` Use the engagement at engagements/${session.engagementSlug}/.`
    : "";

  const prompt = `Run ${skillList}.${engagementCtx}`;
  runClaudeProcess(session, emitSSE, [prompt]);
}

export function resumeClaudeSession(
  session: ExecutionSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void,
  message: string
): void {
  const claudeSessionId = getSessionMeta(session).claudeSessionId as string;
  if (!claudeSessionId) {
    emitSSE("error", { text: "No session to resume" });
    return;
  }
  runClaudeProcess(session, emitSSE, [
    "--resume",
    claudeSessionId,
    message,
  ]);
}

function runClaudeProcess(
  session: ExecutionSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void,
  extraArgs: string[]
): void {
  const args = [
    "--print",
    "--output-format",
    "stream-json",
    "--verbose",
    "--permission-mode",
    "acceptEdits",
    ...extraArgs,
  ];

  console.log(`[claude] spawning: claude ${args.join(" ")}`);

  let proc: ChildProcess;

  try {
    proc = spawn("claude", args, {
      cwd: PROJECT_ROOT,
      env: { ...process.env, PATH: process.env.PATH },
      stdio: ["pipe", "pipe", "pipe"],
    });
  } catch (err) {
    const msg = `Failed to spawn claude CLI: ${err instanceof Error ? err.message : String(err)}`;
    console.error(`[claude] ${msg}`);
    emitSSE("error", { text: msg });
    session.status = "error";
    return;
  }

  session.status = "running";

  session.abort = () => {
    proc.kill("SIGTERM");
    session.status = "completed";
  };

  // Track whether Claude asked a question this turn
  const meta = getSessionMeta(session);
  meta._askedQuestion = false;
  meta._lastResultText = "";
  meta._totalCost = (meta._totalCost as number) ?? 0;

  session.sendFollowUp = (message: string) => {
    // The process already exited — need to resume with a new process
    resumeClaudeSession(session, emitSSE, message);
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
        handleStreamMessage(msg, session, emitSSE);
      } catch {
        emitSSE("assistant", { text: line });
      }
    }
  });

  proc.stderr?.on("data", (chunk: Buffer) => {
    const text = chunk.toString().trim();
    if (text) {
      console.error(`[claude stderr] ${text}`);
      if (
        text.toLowerCase().includes("error") &&
        !text.toLowerCase().includes("warning")
      ) {
        emitSSE("error", { text });
      }
    }
  });

  proc.on("close", (code) => {
    console.log(`[claude] process exited with code ${code}`);
    // Flush remaining buffer
    if (buffer.trim()) {
      try {
        const msg = JSON.parse(buffer);
        handleStreamMessage(msg, session, emitSSE);
      } catch {
        if (buffer.trim()) {
          emitSSE("assistant", { text: buffer });
        }
      }
    }

    // If Claude asked a question, don't emit result — wait for user input
    if (meta._askedQuestion && code === 0) {
      session.status = "waiting_input";
      meta._askedQuestion = false; // Reset for next turn
      emitSSE("waiting", { message: "Waiting for your response..." });
      console.log("[claude] session waiting for user input (question asked)");
      return;
    }

    if (session.status === "running") {
      session.status = "completed";
      emitSSE("result", {
        status: code === 0 ? "success" : "error",
        exitCode: code,
        cost: meta._totalCost as number,
      });
    }
  });

  proc.on("error", (err) => {
    console.error(`[claude] spawn error: ${err.message}`);
    session.status = "error";
    emitSSE("error", { text: `Claude process error: ${err.message}` });
  });
}

function handleStreamMessage(
  msg: Record<string, unknown>,
  session: ExecutionSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void
): void {
  const type = msg.type as string;
  const meta = getSessionMeta(session);

  switch (type) {
    case "assistant": {
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
            meta._lastResultText = block.text;
          } else if (block.type === "tool_use") {
            // Track if AskUserQuestion was used
            if (block.name === "AskUserQuestion") {
              meta._askedQuestion = true;
            }
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
      break;
    }

    case "result": {
      if (msg.total_cost_usd != null) {
        meta._totalCost =
          ((meta._totalCost as number) ?? 0) +
          (msg.total_cost_usd as number);

        // Check if the result text ends with a question — Claude is waiting
        const resultText = (msg.result as string) ?? "";
        if (resultText.trim().endsWith("?")) {
          meta._askedQuestion = true;
        }

        // Don't emit result event here — let the close handler decide
        // based on whether a question was asked
        if (meta._askedQuestion) {
          console.log("[claude] detected question in result, will wait for input");
        } else {
          session.status = "completed";
          emitSSE("result", {
            status: (msg.subtype as string) ?? "success",
            cost: meta._totalCost as number,
            result: resultText,
          });
        }
        break;
      }

      // Tool result message
      const message = msg.message as {
        content?: Array<{ type: string; text?: string }>;
      };
      const toolName = (msg.tool_name as string) ?? "unknown";
      if (message?.content) {
        const text = message.content
          .filter((b) => b.type === "text")
          .map((b) => b.text)
          .join("\n");
        if (text) {
          emitSSE("tool_result", { tool: toolName, output: text });
        }
      }
      break;
    }

    case "tool_result": {
      const message = msg.message as {
        content?: Array<{ type: string; text?: string }>;
      };
      const toolName = (msg.tool_name as string) ?? "unknown";
      if (message?.content) {
        const text = message.content
          .filter((b) => b.type === "text")
          .map((b) => b.text)
          .join("\n");
        if (text) {
          emitSSE("tool_result", { tool: toolName, output: text });
        }
      }
      break;
    }

    case "system": {
      const sessionId = msg.session_id as string;
      if (sessionId) {
        meta.claudeSessionId = sessionId;
        console.log(`[claude] session ID: ${sessionId}`);
      }
      break;
    }

    case "rate_limit_event":
      break;

    default:
      console.log(`[claude] unhandled message type: ${type}`);
      break;
  }
}
