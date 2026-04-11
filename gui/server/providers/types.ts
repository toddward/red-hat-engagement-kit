import type { Response } from "express";

export interface StreamEvent {
  type:
    | "assistant"
    | "tool_use"
    | "tool_result"
    | "question"
    | "approval"
    | "skill_start"
    | "skill_complete"
    | "result"
    | "error";
  data: Record<string, unknown>;
}

export interface ExecutionSession {
  id: string;
  skills: string[];
  engagementSlug: string | null;
  provider: string;
  status: "running" | "waiting_input" | "completed" | "error";
  sseClients: Set<Response>;
  pendingResponses: Map<string, { resolve: (value: unknown) => void }>;
  abort: () => void;
  sendFollowUp?: (message: string) => void;
}
