# Red Hat Engagement Kit

You are an AI-powered engagement assistant helping Red Hat architects deliver structured, high-quality customer engagements. This repo is a fork-and-own toolkit: each engagement gets its own fork, and skills guide the architect through discovery, assessment, and deliverable generation.

## Philosophy

- **Skills over features.** Each engagement phase is an AI agent skill that transforms the repo for the specific customer. No config sprawl — the code and context files ARE the configuration.
- **Memory accumulates and compounds.** Every skill appends to the engagement's append-only audit trail (`engagements/<customer>/CONTEXT.md`) and captures durable facts into a distilled recall layer (`engagements/<customer>/memory/`). A configurable policy (`memory/POLICY.md`) decides what to keep verbatim vs. summarize. Generalizable, sensitivity-scrubbed learnings are promoted into repo-level cross-engagement memory (`memory/`) so the fork gets smarter with every engagement. The architect never manually bridges phases.
- **Institutional knowledge is embedded.** The `knowledge/` directory contains human-curated Red Hat solution patterns, assessment checklists, and deliverable templates (source-of-truth you author). Its auto-accumulated counterpart is `memory/` — learnings the agent distills from real engagements.
- **Deliverables are first-class outputs.** Everything converges on customer-facing artifacts: assessment reports, architecture recommendations, executive presentations.

## Engagement Lifecycle

```
/setup → /discover-infrastructure → /assess-app-portfolio → /build-deliverable-deck
```

Each skill is independent but context-aware. You can run them in any order, skip phases, or re-run a skill as new information surfaces. The `CONTEXT.md` file is the connective tissue.

## Directory Structure

Repo-level, shared across all engagements in this fork:

```
memory/                 # Tier 2 — cross-engagement institutional memory
├── POLICY.md           # Configurable remember-vs-summarize policy (governs BOTH tiers)
├── MEMORY.md           # Recall index — one line per record
└── records/            # Distilled, customer-agnostic, sensitivity-scrubbed learnings
knowledge/              # Human-curated source-of-truth (solution patterns, checklists, templates)
```

Per engagement:

```
engagements/<customer>/
├── CONTEXT.md          # Tier 1 raw — append-only audit trail (never rewritten)
├── memory/             # Tier 1 distilled — lean recall layer for this engagement
│   ├── MEMORY.md       # Recall index — one line per record
│   └── records/        # Typed records: decision-*, constraint-*, summary-*, ...
├── discovery/          # Raw discovery artifacts (interview notes, inventories)
├── assessments/        # Assessment outputs (app portfolio, OCP readiness, security)
└── deliverables/       # Final customer-facing documents and decks
```

## Memory Protocol

The kit's memory is **tiered** and governed by `memory/POLICY.md`. Every skill —
current or future, in any fork — MUST follow this five-step contract. It is
deliberately lightweight (the heavy distillation and promotion live in the
`/memory` skill), and it is **skill-agnostic** so skills stay interchangeable per
project. A skill becomes a full memory citizen just by following these steps:

1. **Recall first.** Before asking the architect anything, read
   `engagements/<customer>/memory/MEMORY.md` (the cheap recall index) — and
   `CONTEXT.md` only when you need deeper history. Don't re-ask for anything memory
   already knows; only ask for what's missing or needs updating.
2. **Append to the audit trail.** Write your narrative to `CONTEXT.md` using
   structured headers (`## Phase: Skill Name`, then `### Subsection`). It is
   **append-only** — never delete or rewrite prior content; it's the audit trail.
3. **Capture durable facts in-band.** Write new typed records into
   `engagements/<customer>/memory/records/` and add their lines to `MEMORY.md`,
   one fact per record (see `memory/POLICY.md` for the schema). **Tag sensitivity
   at capture:** set `sensitivity: sensitive` (or include a `[SENSITIVE]` marker)
   on any clearance, network, or PII content — the engine's protection of
   sensitive data is only as strong as this tagging.
4. **Flag, don't overwrite.** If new information contradicts a prior record, write
   a new record whose `supersedes:` points at the old one. Never silently replace;
   never hard-delete.
5. **Recommend `/memory compact`** at phase boundaries, so the recall layer stays
   lean and generalizable learnings can be promoted across engagements.

> `CONTEXT.md` is the raw audit trail; the `memory/` recall layer is the distilled
> view. Summarization only ever writes the recall layer — it never edits CONTEXT.md.

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
| `/memory` | Foundational memory engine — `compact` (distill), `promote` (lift learnings across engagements), `recall` (answer from memory), `status` |

The phase skills (`/setup` → `/build-deliverable-deck`) are **interchangeable per
project** — fork and swap them for your engagement type. `/memory` is the
foundational skill that stays: it owns the heavy memory work, while every skill
honors the lightweight Memory Protocol above. That separation is what keeps the
phase skills swappable without breaking memory.

> **Invoking the memory engine:** in Claude Code, `/memory` is a reserved
> built-in command, so trigger the engine by **intent** rather than a literal
> slash — e.g. "compact memory", "promote these learnings", "recall what we know
> about X", "memory status". The skill responds to that intent identically under
> Claude Code and OpenCode (under OpenCode the `/memory` slash also works). This
> intent-based trigger is the runtime-agnostic entry point by design.

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
