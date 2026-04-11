import { useState, useRef, useCallback } from "react";
import type { StreamMessage, QuestionEvent, ApprovalEvent, ExecutionStatus } from "../types";

let messageCounter = 0;

export function useExecution() {
  const [messages, setMessages] = useState<StreamMessage[]>([]);
  const [status, setStatus] = useState<ExecutionStatus>("idle");
  const [pendingQuestion, setPendingQuestion] = useState<QuestionEvent | null>(null);
  const [pendingApproval, setPendingApproval] = useState<ApprovalEvent | null>(null);
  const [executionId, setExecutionId] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const addMessage = useCallback((msg: Omit<StreamMessage, "id" | "timestamp">) => {
    setMessages((prev) => [
      ...prev,
      { ...msg, id: String(++messageCounter), timestamp: Date.now() },
    ]);
  }, []);

  const startExecution = useCallback(
    async (skills: string[], engagementSlug?: string) => {
      // Clean up prior connection
      eventSourceRef.current?.close();
      setMessages([]);
      setPendingQuestion(null);
      setPendingApproval(null);

      const res = await fetch("/api/execute", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ skills, engagementSlug }),
      });
      const { executionId: id } = await res.json();
      setExecutionId(id);
      setStatus("running");

      const es = new EventSource(`/api/execute/${id}/stream`);

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
        addMessage({ type: "tool_result", tool: data.tool, output: data.output });
      });

      es.addEventListener("question", (e) => {
        const data = JSON.parse(e.data) as QuestionEvent;
        setPendingQuestion(data);
        setStatus("waiting");
      });

      es.addEventListener("approval", (e) => {
        const data = JSON.parse(e.data) as ApprovalEvent;
        setPendingApproval(data);
        setStatus("waiting");
      });

      es.addEventListener("skill_start", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "skill_start", skill: data.skill });
      });

      es.addEventListener("skill_complete", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "skill_complete", skill: data.skill });
      });

      es.addEventListener("result", (e) => {
        const data = JSON.parse(e.data);
        addMessage({ type: "result", status: data.status, cost: data.cost });
        setStatus("done");
        es.close();
      });

      es.addEventListener("error", () => {
        setStatus("error");
      });

      eventSourceRef.current = es;
    },
    [addMessage]
  );

  const respondToQuestion = useCallback(
    async (questionId: string, answers: Record<string, string>) => {
      if (!executionId) return;
      await fetch(`/api/execute/${executionId}/respond`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type: "answer", questionId, answers }),
      });
      setPendingQuestion(null);
      setStatus("running");
    },
    [executionId]
  );

  const respondToApproval = useCallback(
    async (approvalId: string, allow: boolean) => {
      if (!executionId) return;
      await fetch(`/api/execute/${executionId}/respond`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type: "approval", approvalId, allow }),
      });
      setPendingApproval(null);
      setStatus("running");
    },
    [executionId]
  );

  const sendFollowUp = useCallback(
    async (message: string) => {
      if (!executionId) return;
      addMessage({ type: "assistant", text: `**You:** ${message}` });
      await fetch(`/api/execute/${executionId}/respond`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type: "followup", message }),
      });
      setStatus("running");
    },
    [executionId, addMessage]
  );

  const reset = useCallback(() => {
    eventSourceRef.current?.close();
    setMessages([]);
    setStatus("idle");
    setExecutionId(null);
    setPendingQuestion(null);
    setPendingApproval(null);
  }, []);

  return {
    messages,
    status,
    executionId,
    pendingQuestion,
    pendingApproval,
    startExecution,
    respondToQuestion,
    respondToApproval,
    sendFollowUp,
    reset,
  };
}
