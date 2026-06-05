# Configurable Long-Lived Memory — Design Spec

**Date:** 2026-06-05
**Branch:** `feature/memory`
**Status:** Approved design — ready for implementation planning

---

## 1. Problem & Context

The Red Hat Engagement Kit is a fork-and-own toolkit. Today, each engagement
accumulates a single `engagements/<customer>/CONTEXT.md` — a living,
**append-only** memory that every skill reads from and appends to. This works,
but it has two limits:

1. **CONTEXT.md only ever grows.** There is no notion of *deciding* what is
   worth keeping verbatim versus distilling for efficient recall. Over a long
   engagement the file bloats and recall gets noisier.
2. **Nothing compounds across engagements.** Each fork starts cold. Learnings
   from one engagement never make the fork smarter for the next.

We want a **configurable, long-lived memory** that decides *when to remember
detail verbatim* versus *when to summarize for later recall*, governed by
settings committed in the repo.

## 2. Goals

- Add an intelligent memory layer that makes the **remember-vs-summarize**
  decision based on a repo-committed policy.
- Make memory **tiered**: per-engagement working/long-term memory *and* a
  repo-level cross-engagement institutional memory that compounds over time.
- Keep the mechanism **runtime-agnostic** (identical behavior under Claude Code
  and OpenCode) and **air-gap-friendly** (plain local markdown, no external
  calls, no runtime-specific hooks).
- Keep memory **skill-agnostic**: skills are interchangeable per project, so
  memory must not live inside any one domain skill. Current skills become
  *consumers* of memory via a documented contract.
- Preserve the existing invariant: **CONTEXT.md is an append-only audit trail**
  and is never rewritten.

## 3. Non-Goals

- No new runtime dependency, service, database, or API. Everything is markdown
  files plus skill instructions.
- No automatic OS-level hooks firing every turn (not reliably shared across both
  runtimes). Capture is in-band; heavy distillation is an explicit operation.
- Not building the `memory/ → knowledge/` human-curation graduation step now
  (noted as future work).
- Not building the GUI/TUI surface here — but `recall` is designed to be the
  grounding layer those will call.

## 4. Constraints (derived)

| Constraint | Source | Consequence |
|------------|--------|-------------|
| Runtime-agnostic | Genericization commits; OpenCode + Claude Code parity | Instructions + plain files only; no Claude-only hooks |
| Markdown-first | Repo philosophy | Stores and policy are markdown |
| Skill-interchangeable per project | User requirement | Memory = protocol + foundational skill + config, never wired into a domain skill |
| Audit trail is append-only | CLAUDE.md Context File Protocol | Summarization writes a **separate** distilled surface; never edits CONTEXT.md |
| Air-gapped support | README security section | All local, no external calls |
| Sensitive data must never leak | `[SENSITIVE]` tags, pre-push hook, fork privacy | Cross-engagement promotion requires a mandatory sensitivity scrub |

## 5. Locked Design Decisions

1. **Locus:** Tiered — per-engagement *and* cross-engagement memory.
2. **Trigger model:** Hybrid — cheap capture in-band as skills run; heavy
   distillation/promotion via an explicit, on-demand foundational skill.
3. **Per-engagement recall shape:** A directory of typed records + a `MEMORY.md`
   index — **the same shape as the cross-engagement tier.** Uniformity makes
   promotion a simple "lift a scrubbed record up a tier."
4. **Policy:** A single `memory/POLICY.md` with YAML frontmatter (hard knobs) +
   prose (judgment guidance), governing both tiers.

## 6. Architecture & Storage Layout

The key insight that keeps this safe: **summarization never rewrites
CONTEXT.md.** "Summarize for later recall" means distilling detail *into a
separate recall surface*, leaving the audit trail intact.

```
memory/                          # TIER 2 — repo-level, cross-engagement institutional memory
├── POLICY.md                    # the configurable policy (governs BOTH tiers)
├── MEMORY.md                    # index: one line per record, loaded for recall
└── records/                     # distilled, customer-AGNOSTIC, sensitivity-scrubbed learnings
    ├── pattern-airgap-registry-mirroring.md
    └── insight-fedramp-discovery-sequence.md

engagements/<customer>/
├── CONTEXT.md                   # TIER 1 RAW — unchanged: append-only audit trail
├── memory/                      # TIER 1 DISTILLED — new: lean recall layer for THIS engagement
│   ├── MEMORY.md                # index: one line per record
│   └── records/                 # typed records: decision-*, constraint-*, fact-*, summary-*, ...
└── discovery/ assessments/ deliverables/   # unchanged
```

**Relationship to `knowledge/`:** `knowledge/` is human-curated, stable
source-of-truth (you author it). `memory/` is its **auto-accumulated
counterpart** — what the agent distills from real engagements. A learning that
proves broadly true could later graduate `memory/ → knowledge/` (future
human-curation step; out of scope here).

## 7. The Configurable Policy — `memory/POLICY.md`

Hard knobs in YAML frontmatter (a human tunes them, an agent reads them
deterministically) + prose guidance for the LLM's judgment calls.

```yaml
---
summarization:
  level: balanced            # conservative | balanced | aggressive
                             #   master "remember-more vs summarize-more" dial
  trigger: phase-boundary    # phase-boundary | on-demand | size
  size_threshold_lines: 400  # when trigger=size: distill once a tier's recall exceeds this
  keep_recent_phases: 2      # always keep the last N phases' detail un-summarized

always_keep_verbatim:        # record TYPES never summarized away (keys on the `type` field)
  - decision
  - commitment               # dates, SLAs, promises to the customer
  - score                    # maturity scores / quantified findings
  - stakeholder              # named people, roles, contacts
  - risk                     # open risks / conflicts

keep_sensitive_verbatim: true  # SEPARATE sensitivity floor: ANY record with sensitivity: sensitive
                               # (or [SENSITIVE] content) is never summarized or rolled up, regardless
                               # of type. Keys on the `sensitivity` field, not `type`.

promotion:                   # engagement memory → cross-engagement institutional memory
  enabled: true
  trigger: on-demand
  sensitivity_scrub: required   # MUST strip customer-identifying + [SENSITIVE] before promoting
  criteria:                     # promotes only if the agent judges ALL true
    - reusable across customers
    - customer-agnostic after scrub
    - generalizes beyond this engagement

recall:
  load_engagement_memory: always        # always load this engagement's MEMORY.md index
  load_cross_engagement: by-relevance   # by-relevance | always | off
  max_cross_records: 8

retention:
  audit_trail: never_modify             # CONTEXT.md is sacrosanct
  records: append_or_supersede          # never hard-delete; supersede with a pointer
---
# Memory Policy — guidance for judgment calls
# (what "balanced" means in practice, how to scrub for promotion,
#  good-vs-bad record examples, conflict handling...)
```

The `level` dial plus `always_keep_verbatim` together *are* the configurable
"when to remember detail vs. summarize" decision. The `promotion` block with its
mandatory `sensitivity_scrub` is how the fork compounds without ever leaking
customer data — tying directly into the existing `[SENSITIVE]` tags and pre-push
hook.

**Missing policy fallback:** if `memory/POLICY.md` is absent, the agent falls
back to the documented `balanced` defaults above and notes that it did so.

## 8. Shared Record Schema (identical in both tiers)

```markdown
---
type: decision      # decision|commitment|constraint|fact|score|stakeholder|risk|summary|pattern|insight
date: 2026-06-02
verbatim: true      # honored against the always_keep_verbatim policy
phase: discovery
sensitivity: none   # none | sensitive   (sensitive never promotes)
tags: [airgap, ocp]
supersedes: null    # name — or list of names — of record(s) this replaces (no hard deletes)
---
One durable fact / decision / distilled summary. Links with [[other-record]].
```

**`MEMORY.md` index format** (both tiers) — one line per record:

```
- [decision] target-platform — OCP 4.16 on bare metal (2026-06-02)
- [constraint] airgap — diode-only data transfer, disconnected registry (2026-06-02)
- [summary] discovery — VMware-heavy, no containers yet; F5 ingress; AD identity (2026-06-03)
```

`type` is one of: `decision`, `commitment`, `constraint`, `fact`, `score`,
`stakeholder`, `risk`, `summary` (a distilled rollup), `pattern`, `insight` (the
last two are typical Tier-2 promotion outputs).

**Precedence of `verbatim`:** the per-record `verbatim` field is a denormalized
cache of the policy decision — set at capture time so simple consumers don't have
to re-read the policy. The **policy is the source of truth**: `/memory compact`
re-evaluates each record against the *current* rules — its `type` vs.
`always_keep_verbatim` **and** its `sensitivity` field (any `sensitivity:
sensitive` / `[SENSITIVE]` record is kept verbatim regardless of `type`, the same
predicate `promote` uses) — so editing the policy takes effect retroactively even
on records captured earlier.

## 9. The Memory Protocol — the contract every skill honors

Lives in `CLAUDE.md`, replacing the current `## Context File Protocol` section.
Every skill — current or future, in any fork — honors these five steps
generically. Swapping a domain skill changes nothing about memory.

1. **Recall first** — read `engagements/<customer>/memory/MEMORY.md` (cheap
   index) before asking anything; read `CONTEXT.md` only when deeper history is
   needed.
2. **Append to the audit trail** — write your narrative to `CONTEXT.md` exactly
   as today. Unchanged. Sacrosanct.
3. **Capture durable facts in-band** — write new typed records into
   `engagements/<customer>/memory/records/` + add index lines to `MEMORY.md`,
   honoring `always_keep_verbatim` and `[SENSITIVE]`. (The cheap half of
   "hybrid.")
4. **Flag, don't overwrite** — on conflict, write a record that `supersedes:`
   the old one with a pointer. Never hard-delete.
5. **Recommend `/memory compact`** at phase boundaries. (Hands off the expensive
   half.)

A skill becomes a full memory citizen with **one line** — "follow the Memory
Protocol." That is what keeps skills interchangeable.

## 10. The `/memory` Skill — the engine

A new foundational skill (the one a fork keeps even as it swaps the others). It
reads `POLICY.md` and does the heavy, judgment-driven work the protocol leaves
out. It triggers on **natural-language intent**, not just literal slash-args,
because Claude Code headless and OpenCode do not both run literal
`/memory compact`.

| Operation | What it does |
|-----------|--------------|
| `compact` | Distill phases beyond `keep_recent_phases` into `summary` records per `level`; keep `always_keep_verbatim` types untouched. The verbose detail records it rolls up are **superseded, not deleted** — they stay on disk marked `supersedes`/superseded and drop out of the active `MEMORY.md` index so recall stays lean. **Never edits CONTEXT.md.** |
| `promote` | Scan engagement records for promotion candidates → apply mandatory sensitivity scrub → write customer-agnostic records to `memory/records/` + update `memory/MEMORY.md`. |
| `recall`  | Given a topic/question, load engagement memory + relevance-matched cross-engagement records and answer. **This is the grounding layer the leave-behind GUI/TUI assistant calls.** |
| `status`  | Report memory health: audit size vs. recall size, record counts per type, last compaction, pending promotion candidates. |

Single skill with verbs (recommended over splitting). `/recall` may later be a
thin alias for the assistant use case.

## 11. Existing Skills — Light Touch

- Swap `## Context File Protocol` → `## Memory Protocol` in `CLAUDE.md`.
- Add the one-line "follow the Memory Protocol" pointer and a "recommend
  `/memory compact`" line to each current skill's next-steps section
  (`setup`, `discover-infrastructure`, `assess-app-portfolio`,
  `build-deliverable-deck`).
- No domain logic moves into memory; no memory logic moves into domain skills.

The `.template` engagement gains an empty `memory/` scaffold (`MEMORY.md` stub +
`records/.gitkeep`) so new engagements start memory-ready. `setup` creates the
`memory/` directory alongside `discovery/ assessments/ deliverables/`.

## 12. Data Flow

```
Run a skill
  → recall (read engagements/<customer>/memory/MEMORY.md)
  → do the work
  → append narrative to CONTEXT.md            (audit trail)
  → write typed records + index lines         (in-band capture)

At a phase boundary
  → /memory compact   → distill old phases into summary records, flag candidates
  → /memory promote   → scrub + lift generalizable learnings to memory/records/ (Tier 2)

Next engagement
  → recall pulls relevance-matched Tier-2 records → the fork is now smarter
```

## 13. Edge Cases & Error Handling

- **Missing `POLICY.md`** → fall back to documented `balanced` defaults; note it.
- **Multiple engagements** → ask which customer (consistent with other skills).
- **Air-gapped** → all local; no external calls anywhere in the flow.
- **Promotion safety** → hard-blocked on `sensitivity: sensitive` records; scrub
  is mandatory; respects `.gitignore` and the pre-push hook.
- **Audit invariant** → the `/memory` skill must refuse to modify `CONTEXT.md`;
  if asked, it explains and redirects to the records layer.
- **Conflicts** → preserved from the existing protocol: write a `supersedes:`
  record, never silently overwrite.
- **Empty/new engagement** → `recall` and `status` degrade gracefully (report
  "no memory yet").

## 14. Security & Sensitivity

- The mandatory promotion scrub is the single most important safety control:
  nothing reaches the shared, cross-engagement tier without customer-identifying
  detail and `[SENSITIVE]` content removed.
- `sensitivity: sensitive` records are confined to their engagement directory
  and never promoted.
- All existing protections (gitignore patterns for `.sensitive`/`.classified`,
  the pre-push hook) remain in force and unchanged.

## 15. Testing

This is an instructions+files system, so testing is scenario-based:

- A tiny **fixture engagement** with a few phases of seeded CONTEXT.md + records.
- Evals asserting:
  1. `compact` leaves `CONTEXT.md` byte-for-byte unchanged.
  2. Recall index stays lean after compaction (old detail becomes `summary`
     records).
  3. `always_keep_verbatim` record types survive compaction untouched.
  4. Promoted records are scrubbed (no customer name / `[SENSITIVE]`) and
     generalized.
  5. `conservative` vs. `aggressive` `level` produces measurably different
     output (more vs. less verbatim retained).
  6. Missing `POLICY.md` falls back to `balanced` and says so.
- The `skill-creator` eval harness can drive these scenarios.

## 16. Forward Compatibility

- `recall` is, by design, the grounding layer for the GUI/TUI leave-behind
  assistant on the roadmap (`project_gui_lifecycle`, `project_tui_workbench`).
- The typed-record + `MEMORY.md` index shape is exactly what an artifact browser
  renders, and is consistent with how role-based checklists are envisioned to
  live as markdown.

## 17. File Manifest (for the implementation plan)

**New files:**
- `memory/POLICY.md` — policy (knobs + guidance) and documented defaults.
- `memory/MEMORY.md` — Tier-2 index (starts empty with a header).
- `memory/records/.gitkeep`
- `.claude/skills/memory/SKILL.md` — the `/memory` engine skill.
- `engagements/.template/memory/MEMORY.md` — per-engagement index stub.
- `engagements/.template/memory/records/.gitkeep`
- A fixture engagement + eval notes under `docs/` or a `tests/`-style location
  (final location decided in the plan).

**Modified files:**
- `CLAUDE.md` — replace `## Context File Protocol` with `## Memory Protocol`
  (the five-step contract); add `memory/` to the directory structure; document
  the tiered model and the foundational `/memory` skill in "Available Skills."
- `.claude/skills/setup/SKILL.md` — create `memory/` in the engagement
  scaffold; mention the Memory Protocol.
- `.claude/skills/discover-infrastructure/SKILL.md`,
  `.claude/skills/assess-app-portfolio/SKILL.md`,
  `.claude/skills/build-deliverable-deck/SKILL.md` — one-line protocol pointer +
  "recommend `/memory compact`" in next-steps.
- `README.md` — document the memory model and `/memory` skill.
- `engagements/.template/CONTEXT.md` — note that distilled recall lives under
  `memory/` while CONTEXT.md remains the audit trail.

## 18. Open Questions (resolve during planning)

- Exact home for the fixture engagement + evals (`docs/` vs a `tests/` dir).
- Whether `summary` rollups are one-per-phase or one-per-domain at
  `aggressive` level (default: one-per-phase; revisit if recall still bloats).
- Whether `status` should be a separate cheap path or folded into `recall`.
