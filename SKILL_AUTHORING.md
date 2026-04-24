# Skill Authoring Guide

This guide is for Red Hat practitioners writing new Claude Code skills for this
engagement kit, and for maintainers reviewing those skills before they become
part of a team fork.

The engagement kit is intentionally small: skills are the extension mechanism.
A good skill should capture a repeatable engagement motion, produce durable
customer or practitioner artifacts, and update the engagement's shared context
so later skills can build on the work.

## What A Skill Must Do

Every integrated skill must follow the engagement contract:

- Live under `.claude/skills/<skill-name>/SKILL.md`.
- Include frontmatter with `name` and `description`.
- State when the skill should trigger.
- Define prerequisites, including whether `engagements/<customer>/CONTEXT.md`
  must exist.
- Read engagement context before asking the practitioner for information.
- Avoid re-asking for facts already present in `CONTEXT.md`.
- Write durable artifacts under `engagements/<customer>/`.
- Append a structured summary back to `CONTEXT.md`.
- Recommend the next logical skill or engagement action.

Use these existing skills as references:

- `.claude/skills/setup/SKILL.md` for engagement initialization.
- `.claude/skills/discover-infrastructure/SKILL.md` for structured discovery.
- `.claude/skills/assess-app-portfolio/SKILL.md` for script-backed assessment.
- `.claude/skills/build-deliverable-deck/SKILL.md` for deliverable generation.

## Recommended Skill Shape

Use this structure for new skills unless there is a clear reason to diverge:

```markdown
---
name: assess-example
description: >
  One or two sentences describing the practitioner intent, the engagement phase,
  and the artifact this skill produces.
---

# /assess-example - Example Assessment

## When This Skill Triggers

- Practitioner runs `/assess-example`
- Practitioner asks for the same outcome in natural language

## Prerequisites

- Engagement initialized with `/setup`
- Any required access, files, scripts, or customer inputs

## Workflow

### Step 0: Load Context

Read `engagements/<customer>/CONTEXT.md`. If multiple engagements exist, ask
which customer this skill should use. Capture known facts, known gaps, and any
prior artifacts this skill should reference.

### Step 1: Gather Inputs

Ask only for missing or stale information. Group questions by topic and explain
why each group matters to the assessment or deliverable.

### Step 2: Produce Artifact

Write the durable output to the correct engagement workspace path.

### Step 3: Update Context

Append a concise summary to `CONTEXT.md` using the context protocol.

### Step 4: Recommend Next Steps

Tell the practitioner what skill or action should follow.

## Error Handling

Define what to do when required context, access, files, tools, or customer
answers are missing.
```

## Context Protocol

`CONTEXT.md` is the shared memory for an engagement. Treat it as an append-only
audit trail unless the practitioner explicitly asks to correct an error.

Every skill must:

- Read the full context before collecting new information.
- Use existing metadata such as customer name, engagement type, environment,
  clearance, constraints, and prior findings.
- Append a new `## Phase: <Skill Name>` section with the current date.
- Include generated artifact paths relative to the engagement directory.
- Flag conflicts when new information contradicts prior context.
- Mark sensitive information with `[SENSITIVE]`.

Preferred append format:

```markdown
## Phase: <Skill Name>
**Date:** YYYY-MM-DD
**Conducted by:** <practitioner or automation>

### Summary
<3-5 sentences describing what changed or was learned>

### Key Findings
- <finding>
- <finding>
- <finding>

### Artifacts Produced
- `<relative/path.md>` - <purpose>

### Impact on Engagement Direction
<How this should affect scope, recommendations, roadmap, or next steps>
```

## Artifact Guidelines

Skills should create artifacts that are useful after the Claude session ends.
Do not rely on chat history as the only record of engagement work.

Use these workspace conventions:

- `discovery/` for interview notes, inventories, and current-state findings.
- `assessments/` for scorecards, evidence-based analysis, and readiness reports.
- `deliverables/` for customer-facing documents, decks, and summaries.
- Additional subdirectories are acceptable when the skill owns a clear artifact
  family, but prefer the existing directories first.

Artifact quality expectations:

- Include customer, date, author or collection method, and source assumptions.
- Separate observed facts from recommendations.
- Capture unknowns as follow-up questions instead of inventing answers.
- Use tables for scorecards and decision matrices when they improve clarity.
- Keep customer-facing artifacts free of internal-only notes unless explicitly
  labeled.

## Designing Good Engagement Skills

A skill should represent a bounded engagement capability, not a generic prompt.
Prefer skills that advance one phase or outcome clearly.

Good skill examples:

- Discover the application portfolio and produce a disposition matrix.
- Assess OpenShift readiness and produce a gap report.
- Generate an executive readout from completed discovery and assessments.
- Build a migration wave plan from an approved application inventory.

Weak skill examples:

- Ask random architecture questions without producing an artifact.
- Generate recommendations without reading prior context.
- Mix unrelated outcomes such as kickoff planning, technical assessment, and
  final deck generation in one skill.
- Store important decisions only in chat.

When writing a skill, define:

- The practitioner role using it.
- The engagement phase it supports.
- The specific output it produces.
- The inputs it needs from context, files, scripts, or the practitioner.
- The criteria for completion.

## Knowledge And Tooling Integration

Skills should prefer local knowledge so engagements work in restricted or
air-gapped environments.

- Reference `knowledge/solution-patterns/` for Red Hat architecture patterns.
- Reference `knowledge/checklists/` for discovery and assessment criteria.
- Reference `knowledge/templates/` for reusable output structures.
- Bundle small helper scripts inside the skill directory when automation is part
  of the workflow.
- Document any required external tool, command, or access before the step that
  uses it.

Script-backed skills must explain:

- What the script collects or changes.
- Whether it reads local workstation state, customer environment state, or both.
- Where raw output is stored.
- How the practitioner can review or correct collected data.

## Security And Data Handling

Engagement data may be customer-sensitive. Skills must be conservative by
default.

- Do not instruct practitioners to push customer engagement forks to public repos.
- Mark sensitive values with `[SENSITIVE]` in `CONTEXT.md`.
- Avoid writing secrets, credentials, private keys, tokens, or classified data.
- If sensitive evidence is required, instruct the practitioner to store it in a
  `.sensitive` or `.classified` file covered by `.gitignore`.
- Make air-gapped and limited-connectivity assumptions explicit.
- Keep external references optional unless the engagement requires connected
  operation.

## Maintainer Review Checklist

Use this checklist before accepting a new skill into the kit:

- The skill has clear `name` and `description` frontmatter.
- The trigger conditions are specific and practitioner-oriented.
- The workflow starts by loading `CONTEXT.md`.
- The skill avoids re-asking facts already available in context.
- The skill produces at least one durable artifact or explicitly explains why it
  only updates context.
- Artifact paths are under `engagements/<customer>/`.
- The `CONTEXT.md` append section follows the shared protocol.
- Sensitive-data handling is documented.
- Error handling covers missing engagement context and missing required inputs.
- Next-step recommendations are defined.
- The skill scope is bounded to one engagement outcome.
- Any scripts are local, documented, and reviewed for side effects.

## Acceptance Criteria For A New Skill

A skill is ready to integrate when another practitioner can run it and understand:

- Why they would use it.
- What inputs they need.
- What it will ask them.
- What files it will produce.
- How it updates engagement memory.
- What decision or deliverable it enables next.

If those answers are not obvious from `SKILL.md`, tighten the skill before
adding more functionality.
