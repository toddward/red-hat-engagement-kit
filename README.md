# Red Hat Engagement Kit

An AI-powered engagement delivery toolkit for Red Hat architects. Fork → run your AI coding agent → `/setup` → deliver.

Small enough to understand, built to be customized, skills over features.

Supports both **Claude Code** and **OpenCode** as the AI agent runtime.

---

## What This Is

A structured, skills-driven framework for delivering customer engagements. Each engagement phase is a skill that guides the architect through discovery, assessment, and deliverable generation. Context accumulates in a living `CONTEXT.md` file — later skills build on earlier findings automatically.

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

## How It Works

```
/setup
  └── Creates engagements/<customer>/CONTEXT.md (living memory)

/discover-infrastructure
  ├── Reads CONTEXT.md (knows what's already known)
  ├── Conducts structured interview (adapts to engagement type)
  ├── Writes discovery/infrastructure-discovery.md
  └── Appends findings to CONTEXT.md

/assess-app-portfolio
  ├── Reads CONTEXT.md (builds on infrastructure findings)
  ├── Runs scripts/collect-system-info.sh (demonstrates script execution)
  ├── Presents system landscape to architect for review
  ├── Writes assessments/system-info-collection.md
  └── Appends findings to CONTEXT.md

/build-deliverable-deck
  ├── Reads EVERYTHING (CONTEXT.md + all reports)
  ├── Structures executive narrative from findings
  ├── Generates presentation (HTML Quick Deck or PPTX)
  └── Writes to deliverables/
```

## Customization

**Add a new assessment type:**
Create a skill at `.claude/skills/assess-<topic>/SKILL.md` following the existing patterns. Both Claude Code and OpenCode read from this directory.

It should read from `CONTEXT.md`, conduct an interview or intake, produce a report, and append findings back to `CONTEXT.md`.

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
│   ├── setup/
│   ├── discover-infrastructure/
│   ├── assess-app-portfolio/
│   └── build-deliverable-deck/
├── scripts/                     # Shared scripts (used by both runtimes)
│   └── collect-system-info.sh
├── engagements/                 # Customer engagement workspaces
│   └── .template/
├── knowledge/                   # Institutional knowledge base
│   ├── solution-patterns/
│   ├── checklists/
│   └── templates/
├── CLAUDE.md                    # Claude Code agent instructions
├── AGENTS.md                    # OpenCode agent instructions
├── opencode.json                # OpenCode project config
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
