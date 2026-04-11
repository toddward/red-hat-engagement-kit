export interface SkillInfo {
  name: string;
  description: string;
  path: string;
  content: string;
}

export interface EngagementInfo {
  slug: string;
  hasContext: boolean;
  contextContent?: string;
}
