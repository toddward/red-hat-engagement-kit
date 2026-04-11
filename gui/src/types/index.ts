export interface SkillSummary {
  name: string;
  description: string;
  path: string;
}

export interface SkillDetail {
  name: string;
  description: string;
  content: string;
}

export interface Engagement {
  slug: string;
  hasContext: boolean;
}

export interface StreamMessage {
  id: string;
  type: "assistant" | "tool_use" | "tool_result" | "error" | "skill_start" | "skill_complete" | "result";
  text?: string;
  tool?: string;
  input?: string;
  output?: string;
  skill?: string;
  status?: string;
  cost?: number;
  timestamp: number;
}

export interface QuestionOption {
  label: string;
  description?: string;
}

export interface Question {
  question: string;
  header: string;
  options: QuestionOption[];
  multiSelect: boolean;
}

export interface QuestionEvent {
  questionId: string;
  questions: Question[];
}

export interface ApprovalEvent {
  approvalId: string;
  tool: string;
  input: Record<string, unknown>;
}

export type ExecutionStatus = "idle" | "running" | "waiting" | "done" | "error";
