import { spawn, type ChildProcess } from "node:child_process";
import path from "node:path";
import type { ExecutionSession } from "./types.js";

const PROJECT_ROOT = path.resolve(import.meta.dirname, "../../..");

export function startClaudeSession(
  session: ExecutionSession,
  emitSSE: (event: string, data: Record<string, unknown>) => void
): void {
  const skillList = session.skills.map((s) => `/${s}`).join(", then ");
  const engagementCtx = session.engagementSlug
    ? ` Use the engagement at engagements/${session.engagementSlug}/.`
    : "";

  const prompt = `Run ${skillList}.${engagementCtx}`;

  const args = [
    "--print",
    "--output-format",
    "stream-json",
    "--verbose",
    "--permission-mode",
    "acceptEdits",
    prompt,
  ];

  console.log(`[claude] spawning: claude ${args.join(" ")}`);
  console.log(`[claude] cwd: ${PROJECT_ROOT}`);

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

  session.abort = () => {
    proc.kill("SIGTERM");
    session.status = "completed";
  };

  session.sendFollowUp = (message: string) => {
    if (proc.stdin && !proc.stdin.destroyed) {
      const msg = JSON.stringify({
        type: "user_message",
        content: message,
      });
      proc.stdin.write(msg + "\n");
    }
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
        // Non-JSON output, treat as raw text
        emitSSE("assistant", { text: line });
      }
    }
  });

  proc.stderr?.on("data", (chunk: Buffer) => {
    const text = chunk.toString().trim();
    if (text) {
      console.error(`[claude stderr] ${text}`);
      // Surface errors to the frontend
      if (
        text.toLowerCase().includes("error") ||
        text.toLowerCase().includes("fatal")
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

    if (session.status === "running") {
      session.status = "completed";
      emitSSE("result", {
        status: code === 0 ? "success" : "error",
        exitCode: code,
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
      break;
    }

    case "tool_result":
    case "result": {
      // Check if this is a final result message (has total_cost_usd)
      if (msg.total_cost_usd != null) {
        session.status = "completed";
        emitSSE("result", {
          status: (msg.subtype as string) ?? "success",
          cost: msg.total_cost_usd as number,
          result: msg.result as string,
        });
        break;
      }

      // Otherwise it's a tool result
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
        (session as unknown as Record<string, unknown>).claudeSessionId =
          sessionId;
      }
      break;
    }

    case "rate_limit_event":
      // Ignore rate limit events
      break;

    default:
      console.log(`[claude] unhandled message type: ${type}`);
      break;
  }
}
