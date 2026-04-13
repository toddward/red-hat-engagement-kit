# Engagement Team Agents

This directory defines specialized sub-agents that form the engagement delivery team. Each agent has a distinct role, expertise, and model configuration.

## Team Composition

| Agent | Model | Role |
|-------|-------|------|
| **Architect** | Opus | Team lead — understands the project holistically, coordinates work, makes architectural decisions |
| **Senior Developer** | Sonnet | Implementation — writes production code, debugs issues, executes technical tasks |
| **QA Specialist** | Opus | Quality — tests deliverables, validates requirements, ensures standards |
| **Documentation Specialist** | Opus | Writing — creates customer-facing docs, maintains CONTEXT.md, produces deliverables |

## How to Invoke

Use the Agent tool with `subagent_type` or spawn them via the CLI:

### From Claude Code

```typescript
// Invoke the Architect for project planning
Agent({
  description: "Architect: plan engagement phases",
  prompt: "Review CONTEXT.md and break down the next phase into tasks for the team",
  model: "opus"
})

// Invoke the Senior Developer for implementation
Agent({
  description: "Developer: implement discovery script",
  prompt: "Create a shell script that collects system inventory information",
  model: "sonnet"
})

// Invoke QA for validation
Agent({
  description: "QA: validate assessment report",
  prompt: "Review the assessment report against requirements and document any gaps",
  model: "opus"
})

// Invoke Documentation Specialist for writing
Agent({
  description: "Docs: write executive summary",
  prompt: "Create an executive summary from the assessment findings in CONTEXT.md",
  model: "opus"
})
```

### Parallel Execution

For independent tasks, invoke multiple agents in parallel:

```typescript
// Kick off implementation and documentation simultaneously
Agent({
  description: "Developer: implement feature",
  prompt: "...",
  model: "sonnet",
  run_in_background: true
})

Agent({
  description: "Docs: write user guide",
  prompt: "...",
  model: "opus",
  run_in_background: true
})
```

## Agent Coordination Patterns

### Architect-Led Delegation

```
User Request
    │
    ▼
Architect (analyze, plan, delegate)
    │
    ├──► Senior Developer (implement)
    │         │
    │         ▼
    │    QA Specialist (test)
    │         │
    └──► Documentation Specialist (document)
              │
              ▼
         Architect (review, deliver)
```

### Parallel Workstreams

```
Architect (define tasks)
    │
    ├──► Senior Developer ─► QA ─┐
    │                            │
    └──► Documentation ──────────┴──► Architect (integrate)
```

## When to Use Each Agent

| Situation | Agent |
|-----------|-------|
| "Plan the next phase" | Architect |
| "Break this down into tasks" | Architect |
| "Write a script to..." | Senior Developer |
| "Fix this bug" | Senior Developer |
| "Validate the deliverables" | QA Specialist |
| "Test this implementation" | QA Specialist |
| "Write the executive summary" | Documentation Specialist |
| "Update CONTEXT.md" | Documentation Specialist |

## Context Sharing

All agents should read `engagements/<customer>/CONTEXT.md` as their first action. This ensures:

- Consistent understanding of engagement scope
- Awareness of prior findings and decisions
- Alignment with customer constraints and goals

## Model Selection Rationale

- **Opus** for Architect, QA, and Documentation: These roles require nuanced judgment, holistic understanding, and careful reasoning about complex tradeoffs.
- **Sonnet** for Senior Developer: Implementation tasks benefit from Sonnet's speed while maintaining high code quality. Escalate to Opus for complex debugging or architectural decisions.
