import { useState, useRef, useCallback } from "react";
import type { StreamMessage, ExecutionStatus } from "../types";

let messageCounter = 0;

export function useAssistant(slug: string | null) {
  const [messages, setMessages] = useState<StreamMessage[]>([]);
  const [status, setStatus] = useState<ExecutionStatus>("idle");
  const [sessionId, setSessionId] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const addMessage = useCallback(
    (msg: Omit<StreamMessage, "id" | "timestamp">) => {
      setMessages((prev) => [
        ...prev,
        { ...msg, id: String(++messageCounter), timestamp: Date.now() },
      ]);
    },
    []
  );

  const connectSSE = useCallback(
    (slug: string, sid: string) => {
      const es = new EventSource(
        `/api/engagements/${slug}/assistant/${sid}/stream`
      );

      es.addEventListener("assistant", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "assistant", text: data.text });
      });

      es.addEventListener("tool_use", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "tool_use", tool: data.tool, input: data.input });
      });

      es.addEventListener("tool_result", (e) => {
        const data = JSON.parse(e.data);
        addMessage({
          type: "tool_result",
          tool: data.tool,
          output: data.output,
        });
      });

      es.addEventListener("result", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "result", status: data.status, cost: data.cost });
        setStatus("done");
        // Don't close the EventSource — keep it for follow-ups
      });

      es.addEventListener("error", () => {
        setStatus("error");
      });

      eventSourceRef.current = es;
    },
    [addMessage]
  );

  const askQuestion = useCallback(
    async (question: string) => {
      if (!slug) return;

      addMessage({ type: "assistant", text: `**You:** ${question}` });
      setStatus("running");

      if (sessionId) {
        // Follow-up on existing session — use the followup endpoint
        await fetch(
          `/api/engagements/${slug}/assistant/${sessionId}/followup`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ question }),
          }
        );
        // SSE connection is still open, new events will stream in
      } else {
        // First question — create a new session
        const res = await fetch(`/api/engagements/${slug}/assistant`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ question }),
        });
        const data = await res.json();
        setSessionId(data.sessionId);
        connectSSE(slug, data.sessionId);
      }
    },
    [slug, sessionId, addMessage, connectSSE]
  );

  const reset = useCallback(() => {
    eventSourceRef.current?.close();
    eventSourceRef.current = null;
    setMessages([]);
    setStatus("idle");
    setSessionId(null);
  }, []);

  return { messages, status, askQuestion, reset };
}
