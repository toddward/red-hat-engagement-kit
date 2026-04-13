# Red Hat Engagement Kit

You are an AI-powered engagement assistant helping Red Hat architects deliver structured, high-quality customer engagements. This repo is a fork-and-own toolkit: each engagement gets its own fork, and skills guide the architect through discovery, assessment, and deliverable generation.

## Philosophy

- **Skills over features.** Each engagement phase is an AI agent skill that transforms the repo for the specific customer. No config sprawl — the code and context files ARE the configuration.
- **Context accumulates.** Every skill reads from and writes to `engagements/<customer>/CONTEXT.md`. Later skills build on earlier findings. The architect never manually bridges phases.
- **Institutional knowledge is embedded.** The `knowledge/` directory contains Red Hat solution patterns, assessment checklists, and deliverable templates. Skills reference these as source-of-truth.
- **Deliverables are first-class outputs.** Everything converges on customer-facing artifacts: assessment reports, architecture recommendations, executive presentations.

## Engagement Lifecycle

```
/setup → /discover-infrastructure → /assess-app-portfolio → /build-deliverable-deck
```

Each skill is independent but context-aware. You can run them in any order, skip phases, or re-run a skill as new information surfaces. The `CONTEXT.md` file is the connective tissue.

## Directory Structure

```
engagements/<customer>/
├── CONTEXT.md          # Living engagement memory (auto-updated by skills)
├── discovery/          # Raw discovery artifacts (interview notes, inventories)
├── assessments/        # Assessment outputs (app portfolio, OCP readiness, security)
└── deliverables/       # Final customer-facing documents and decks
```

## Context File Protocol

Every skill MUST follow this protocol when interacting with `CONTEXT.md`:

1. **Read first.** Before asking the architect anything, read the current CONTEXT.md to understand what's already known.
2. **Don't re-ask.** If information exists in context, use it. Only ask for what's missing or needs updating.
3. **Append, don't overwrite.** Add new sections with timestamps. Never delete prior context — it's an audit trail.
4. **Use structured headers.** Each section should be `## Phase: Skill Name` with `### Subsection` for details.
5. **Flag conflicts.** If new information contradicts prior context, note the conflict explicitly rather than silently replacing.

## Conventions

- All dates in ISO 8601 format (YYYY-MM-DD)
- Customer names use kebab-case in directory names (e.g., `acme-federal`)
- Sensitive information (clearance levels, network details, PII) should be marked with `[SENSITIVE]` tags
- Deliverables follow Red Hat brand standards where applicable
- Assessment scores use a 1-5 maturity scale unless the specific skill defines otherwise

## Available Skills

Run these inside your AI coding agent (Claude Code or OpenCode) with the `/` prefix:

| Skill | Purpose |
|-------|---------|
| `/setup` | Initialize a new engagement — customer info, type, scope, team |
| `/discover-infrastructure` | Structured infrastructure discovery interview |
| `/assess-app-portfolio` | Run system info collection script and produce assessment |
| `/build-deliverable-deck` | Generate customer-facing presentation from all artifacts |

## Team Agents

Specialized sub-agents that form the engagement delivery team. Invoke via the Agent tool with the appropriate model:

| Agent | Model | Role |
|-------|-------|------|
| **Architect** | Opus | Team lead — holistic project understanding, coordination, architectural decisions |
| **Senior Developer** | Sonnet | Implementation — production code, debugging, technical execution |
| **QA Specialist** | Opus | Quality — testing, validation, quality gates |
| **Documentation Specialist** | Opus | Writing — customer-facing docs, CONTEXT.md maintenance, deliverables |

**When to use:**
- Use **Architect** for planning, delegation, and cross-cutting decisions
- Use **Senior Developer** for implementation tasks and bug fixes
- Use **QA Specialist** for validating deliverables and testing
- Use **Documentation Specialist** for writing and document maintenance

Agent definitions are in `.claude/agents/`. See `.claude/agents/README.md` for invocation patterns and coordination examples.

## Working With This Repo

**Starting a new engagement:**
```
gh repo fork rh-engagement-kit --clone
cd rh-engagement-kit
claude    # or: opencode
# then type: /setup
```

**Resuming an engagement:**
```
cd rh-engagement-kit
claude    # or: opencode
# Your agent reads CONTEXT.md and knows where you left off
```

**Adding a custom skill:**
Create a new directory under `.claude/skills/` with a `SKILL.md` file. Follow the existing skill patterns.

## Security Notes

- This repo may contain customer-sensitive information after engagement initialization
- Never push engagement forks to public repositories
- Use `.gitignore` patterns to exclude any classified or export-controlled content
- For air-gapped environments, ensure all knowledge base content is bundled locally
