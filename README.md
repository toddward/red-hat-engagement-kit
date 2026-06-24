# Red Hat Engagement Kit

An AI-powered engagement delivery toolkit for Red Hat architects. Fork → run your AI coding agent → `/setup` → deliver.

Small enough to understand, built to be customized, skills over features.

Supports both **Claude Code** and **OpenCode** as the AI agent runtime.

---

## What This Is

A structured, skills-driven framework for delivering customer engagements. Each engagement phase is a skill that guides the architect through discovery, assessment, and deliverable generation. Memory accumulates as the engagement runs — an append-only `CONTEXT.md` audit trail plus a distilled, configurable recall layer that later skills (and future engagements) build on automatically.

**This is not a SaaS product or a framework with dependencies.** It's a Git repo with markdown files that your AI coding agent knows how to execute. Fork it, customize it, run engagements with it.

## Quick Start

### With Claude Code

```bash
gh repo fork <org>/rh-engagement-kit --clone
cd rh-engagement-kit
claude
```

Then type `/setup` inside Claude Code.

### With OpenCode

```bash
gh repo fork <org>/rh-engagement-kit --clone
cd rh-engagement-kit
opencode
```

Then type `/setup` inside OpenCode.

## Skills

| Skill | Phase | What It Does |
|-------|-------|-------------|
| `/setup` | Initialize | Gather engagement metadata, create workspace, write initial CONTEXT.md |
| `/discover-infrastructure` | Discovery | Structured infrastructure interview across 6 domains, maturity scoring |
| `/assess-app-portfolio` | Assessment | Run system info collection script, present findings, produce assessment report |
| `/build-deliverable-deck` | Delivery | Customer-facing executive presentation (Quick Deck or PPTX) |
| `/memory` | Foundational | Long-lived memory engine — `compact` (distill), `promote` (across engagements), `recall`, `status` |

## How It Works

Every skill follows the lightweight **Memory Protocol** (defined in `CLAUDE.md`): recall from the `memory/` index first, append narrative to the append-only `CONTEXT.md` audit trail, and capture durable facts as typed records — so context compounds without anyone hand-bridging phases.

```
/setup
  ├── Gathers engagement metadata via structured interview
  ├── Creates engagements/<customer>/ workspace (discovery/, assessments/, deliverables/)
  └── Writes CONTEXT.md (audit trail) + memory/ recall layer (MEMORY.md index + records/)

/discover-infrastructure
  ├── Recalls from memory/MEMORY.md first (then CONTEXT.md for deeper history)
  ├── Conducts structured interview (adapts to engagement type)
  ├── Writes discovery/infrastructure-discovery.md
  ├── Appends findings to CONTEXT.md + captures typed records
  └── Recommends /memory compact at the phase boundary

/assess-app-portfolio
  ├── Recalls from memory/MEMORY.md first (then CONTEXT.md)
  ├── Runs .claude/skills/assess-app-portfolio/collect-system-info.sh (demonstrates script execution)
  ├── Presents system landscape to architect for review
  ├── Writes assessments/system-info-collection.md
  ├── Appends findings to CONTEXT.md + captures typed records
  └── Recommends /memory compact

/build-deliverable-deck
  ├── Recalls from memory/MEMORY.md, then reads EVERYTHING (CONTEXT.md + all reports)
  ├── Structures executive narrative from findings
  ├── Generates presentation (HTML Quick Deck or PPTX)
  ├── Writes to deliverables/
  └── Offers /memory compact + /memory promote (lift reusable learnings cross-engagement)
```

## Memory

The kit has a **configurable, tiered memory** so engagements stay sharp and the fork gets smarter over time.

- **Per-engagement (Tier 1).** `engagements/<customer>/CONTEXT.md` stays the append-only audit trail. A new `engagements/<customer>/memory/` holds the distilled recall layer — typed records (decisions, constraints, scores, risks, summaries) plus a `MEMORY.md` index skills read first.
- **Cross-engagement (Tier 2).** Repo-level `memory/` accumulates customer-agnostic, sensitivity-scrubbed learnings promoted from real engagements. It's the auto-accumulated counterpart to the human-curated `knowledge/` base.

**Configurable.** `memory/POLICY.md` decides *when to keep detail verbatim vs. summarize for later recall* — a master `level` dial (`conservative` / `balanced` / `aggressive`), a list of record types always kept verbatim, and a mandatory sensitivity scrub before anything is promoted across engagements.

**The `/memory` skill** is the engine: `compact` distills old detail into summaries (never touching the audit trail), `promote` lifts scrubbed learnings to Tier 2, `recall` answers questions grounded in memory (what a leave-behind assistant calls), and `status` reports memory health. Every other skill stays interchangeable by following the lightweight **Memory Protocol** in `CLAUDE.md`. Invoke it by intent — "compact memory", "recall what we know about X", "memory status" — which works identically across runtimes (on Claude Code, `/memory` is a reserved built-in, so the intent trigger is the canonical entry point).

Everything is local markdown — no service, no database, air-gap friendly. Nothing customer-identifying ever reaches the shared cross-engagement tier.

## Writing New Skills

The kit is **skills over features** — you extend it by adding skills, not code. A skill is just a directory the AI agent knows how to execute, and both Claude Code and OpenCode read from `.claude/skills/`.

### Anatomy of a skill

```
.claude/skills/<skill-name>/
├── SKILL.md          # Required — frontmatter (name, description) + instructions
├── *.sh / *.py       # Optional — scripts the skill runs (e.g. collect-system-info.sh)
└── hooks/            # Optional — lifecycle hooks (e.g. setup/hooks/pre-push)
```

The `description` is what *triggers* the skill — write it for the intent you want the agent to fire on, not just what it does.

### Every skill is a memory citizen

A new phase skill is only interchangeable if it honors the **Memory Protocol** (defined in `CLAUDE.md`). That protocol is the contract — five steps, every skill, every run:

1. **Recall first** — read `engagements/<customer>/memory/MEMORY.md` before asking the architect anything.
2. **Append to the audit trail** — write narrative to `CONTEXT.md` under `## Phase: <Skill>` headers (append-only, never rewrite).
3. **Capture durable facts** — write typed records to `memory/records/`, one fact each, indexed in `MEMORY.md`. Tag sensitivity at capture.
4. **Flag, don't overwrite** — contradictions get a new record with `supersedes:`, never a silent edit.
5. **Recommend `/memory compact`** at the phase boundary.

### Skill scaffold

Copy this into a new `.claude/skills/<skill-name>/SKILL.md` and fill it in — the Memory Protocol steps are pre-wired so you can't forget them:

```markdown
---
name: <skill-name>
description: <intent that should trigger this skill — be specific>
---

# /<skill-name> — <one-line purpose>

## 1. Recall
Read engagements/<customer>/memory/MEMORY.md first. Read CONTEXT.md only for
deeper history. Don't re-ask for anything memory already knows.

## 2. Do the work
<interview / run a script / generate an artifact — your skill's actual job>.
Write outputs to engagements/<customer>/{discovery,assessments,deliverables}/.

## 3. Append to the audit trail
Append findings to engagements/<customer>/CONTEXT.md under
`## Phase: <Skill Name>`. Never rewrite prior content.

## 4. Capture durable facts
Write typed records (decision-*, constraint-*, risk-*, summary-*) to
engagements/<customer>/memory/records/, add their lines to MEMORY.md, and tag
`sensitivity: sensitive` on any clearance / network / PII content.

## 5. Close the loop
If new info contradicts a record, supersede it (don't overwrite). Recommend
`/memory compact` at the phase boundary.
```

### Enforcing memory utilization

The protocol only compounds if every skill actually follows it. Layered ways to enforce it, cheapest first:

- **Scaffold by default.** Start every new skill from the template above so the five steps ship with the skill. The cheapest enforcement is making the right thing the default.
- **CLAUDE.md is the always-on backstop.** The Memory Protocol lives in `CLAUDE.md` and is marked as overriding default agent behavior, so the agent honors it even when a skill's own instructions are thin. Keep it there — don't let skills restate (and drift from) it.
- **Definition of Done.** End each `SKILL.md` with a self-check the agent must pass before reporting completion: *did I recall first, append to `CONTEXT.md`, write ≥1 record and index it, and tag sensitivity?* A skill that can't tick these isn't done.
- **Lint new skills in a hook.** The repo already ships hooks (`setup/hooks/pre-push`). Add a pre-commit check that greps each `SKILL.md` for the protocol touchpoints and fails when one is missing:

  ```bash
  for s in .claude/skills/*/SKILL.md; do
    grep -q "MEMORY.md"  "$s" || { echo "FAIL: $s never recalls memory";            exit 1; }
    grep -q "CONTEXT.md" "$s" || { echo "FAIL: $s never appends to the audit trail"; exit 1; }
    grep -q "records/"   "$s" || { echo "FAIL: $s never captures records";           exit 1; }
  done
  ```

- **Audit with `/memory status`.** At each phase boundary, "memory status" reports record counts and last compaction. If a skill ran but record counts didn't grow, it skipped step 3 — that's your signal it isn't a full memory citizen yet.

## Customization

**Add a new assessment type:**
Create a skill at `.claude/skills/assess-<topic>/SKILL.md` — see [Writing New Skills](#writing-new-skills) for the scaffold and the Memory Protocol contract every skill must honor.

**Modify an assessment:**
Edit the SKILL.md directly. The skills are just markdown instructions — change the questions, scoring criteria, or output format to match your methodology.

**Add knowledge base content:**
Drop reference material into `knowledge/solution-patterns/`, `knowledge/checklists/`, or `knowledge/templates/`. Skills will reference it.

**Fork for your team:**
Each team or practice area can maintain their own fork with customized skills, checklists, and solution patterns. The base repo provides the framework; your fork encodes your team's methodology.

## Repository Structure

```
rh-engagement-kit/
├── .claude/skills/              # Skills (shared by Claude Code & OpenCode)
│   ├── setup/                   # Includes hooks/pre-push
│   ├── discover-infrastructure/
│   ├── assess-app-portfolio/    # Includes collect-system-info.sh
│   ├── build-deliverable-deck/
│   └── memory/                  # Foundational long-lived memory engine
├── engagements/                 # Customer engagement workspaces
│   └── .template/               # CONTEXT.md (audit trail) + memory/ (recall layer)
├── memory/                      # Cross-engagement institutional memory (Tier 2)
│   ├── POLICY.md                # Configurable remember-vs-summarize policy
│   ├── MEMORY.md                # Recall index
│   └── records/                 # Scrubbed, customer-agnostic learnings
├── knowledge/                   # Human-curated knowledge base
│   ├── solution-patterns/
│   ├── checklists/
│   └── templates/
├── CLAUDE.md                    # Agent instructions (shared by Claude Code & OpenCode)
├── opencode.json                # OpenCode project config (points to CLAUDE.md)
└── README.md
```

## Security & Data Handling

- **Never push customer forks to public repos.** Engagement data is customer-sensitive.
- **Air-gapped support.** All knowledge base content is local. No external API calls required during engagement execution (beyond the AI agent itself).
- **Sensitive data tagging.** Skills mark sensitive information with `[SENSITIVE]` tags in CONTEXT.md.
- **Classification boundaries.** The `.gitignore` excludes `.sensitive` and `.classified` files.

## Contributing

**Don't add features, add skills.**

Want to add a new assessment type? Create a skill. Want to support a different deliverable format? Create a skill. The base repo stays minimal — your fork encodes your specific needs.

Add new skills to `.claude/skills/` — both Claude Code and OpenCode read from this directory.

## Requirements

- [Claude Code](https://claude.ai/download) **or** [OpenCode](https://github.com/opencode-ai/opencode)
- Git
- That's it.

## License

Internal Red Hat use. See your team's guidelines for external distribution.
