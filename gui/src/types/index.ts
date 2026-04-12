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

export type EngagementPhase = "pre-engagement" | "live" | "leave-behind";

export interface PhaseInfo {
  phase: EngagementPhase;
  hasContext: boolean;
  hasDeliverables: boolean;
  artifactCounts: {
    discovery: number;
    assessments: number;
    deliverables: number;
    checklists: number;
  };
}

export interface ArtifactNode {
  name: string;
  path: string;
  type: "file" | "directory";
  children?: ArtifactNode[];
}

export interface ChecklistItem {
  text: string;
  checked: boolean;
  line: number;
}

export interface ChecklistSection {
  title: string;
  items: ChecklistItem[];
}

export interface Checklist {
  name: string;
  fileName: string;
  sections: ChecklistSection[];
  completionPercent: number;
}

export type AppView = "detail" | "execution" | "artifacts" | "checklists" | "assistant";
