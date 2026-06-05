---
# ─────────────────────────────────────────────────────────────────────────────
# Memory Policy — the configurable knobs that govern remember-vs-summarize.
# A human tunes these; the agent reads them deterministically.
# Governs BOTH tiers: per-engagement (engagements/<customer>/memory/) and the
# repo-level cross-engagement store (memory/).
# If this file is missing, the agent falls back to these exact defaults and says so.
# ─────────────────────────────────────────────────────────────────────────────

summarization:
  level: balanced            # conservative | balanced | aggressive
                             #   the master "remember-more vs summarize-more" dial
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

# Memory Policy — Guidance

This file is the source of truth for how memory behaves. The YAML frontmatter
above is the set of hard knobs; the prose below tells the agent how to apply the
judgment calls the knobs can't fully specify. The `/memory` skill and every skill
that follows the **Memory Protocol** (see `CLAUDE.md`) read this file.

> **If this file is absent**, behave exactly as the frontmatter defaults above
> (`level: balanced`, etc.) and note in your output that you fell back to
> defaults because no `POLICY.md` was found.

---

## The two tiers this governs

- **Tier 1 — Engagement memory** (`engagements/<customer>/memory/`): the lean,
  distilled recall layer for one engagement. `CONTEXT.md` in the same directory
  remains the **append-only audit trail** and is *never* modified by memory
  operations.
- **Tier 2 — Cross-engagement memory** (`memory/`, this directory): customer-
  agnostic, sensitivity-scrubbed learnings that persist across every engagement
  in the fork and make it smarter over time.

Both tiers use the identical record schema and `MEMORY.md` index format below, so
promoting a learning is just "lift a scrubbed record up a tier."

---

## The `level` dial — what each setting means

`level` is the single most important knob. It sets how eagerly detail is distilled.

| Level | Distill when | Keeps verbatim | Use when |
|-------|--------------|----------------|----------|
| `conservative` | Only well past `keep_recent_phases`, and only clearly redundant detail | More — when in doubt, keep it | High-stakes / audited engagements where losing nuance is costly |
| `balanced` (default) | At each phase boundary, for phases older than `keep_recent_phases` | The `always_keep_verbatim` types + anything still actively referenced | Most engagements |
| `aggressive` | As soon as a phase completes; rolls older summaries into higher-level ones | Only the `always_keep_verbatim` types | Long engagements where recall is getting noisy and speed matters |

`level` never overrides `always_keep_verbatim`. Those record types stay verbatim
at every level — that is the floor.

## `keep_recent_phases` and `trigger`

- `keep_recent_phases` is a protected window: the most recent N phases' detail is
  never summarized, regardless of `level`. Default `2`.
- `trigger` decides *when* compaction is considered:
  - `phase-boundary` (default) — when the architect moves to a new phase / runs
    the next skill.
  - `on-demand` — only when `/memory compact` is invoked explicitly.
  - `size` — when a tier's active recall (its un-superseded records) exceeds
    `size_threshold_lines`.
- Capture (writing new records as skills run) happens **in-band regardless of
  trigger**. The trigger only gates the heavier *distillation* step.

---

## What "remember in full" means — `always_keep_verbatim`

Two independent rules decide what is kept verbatim (never collapsed into a
summary). A record is verbatim if **either** rule matches:

**1. Verbatim by type** — a record whose `type` is in `always_keep_verbatim`:

- `decision` — choices made and the reason for them.
- `commitment` — dates, SLAs, and promises made to the customer.
- `score` — maturity scores and other quantified findings.
- `stakeholder` — named people, their roles, and contact context.
- `risk` — open risks and unresolved conflicts.

**2. Verbatim by sensitivity** (`keep_sensitive_verbatim: true`) — ANY record
with `sensitivity: sensitive` or `[SENSITIVE]` content, **regardless of its
`type`**. This keys on the `sensitivity` field, not `type`, so a
`sensitivity: sensitive` record stored under an ordinary type (e.g.
`type: constraint`) is still never summarized, never rolled up, and never dropped
from the index. This is the **same predicate `promote` uses** to block sensitive
records — `compact` and `promote` share it, so nothing can be laundered from one
path to the other. This guard is only as strong as capture-time tagging: skills
MUST set `sensitivity: sensitive` (or include a `[SENSITIVE]` marker) on any
clearance, network, or PII content when they create a record — see the Memory
Protocol in `CLAUDE.md`.

**Precedence:** the per-record `verbatim` field is a denormalized cache of these
rules, set at capture time. The **policy is the source of truth** — `/memory
compact` re-evaluates each record against the *current* rules (its `type` vs.
`always_keep_verbatim`, and its `sensitivity`), so editing this policy takes
effect retroactively.

---

## Record schema (identical in both tiers)

Every record is one markdown file with frontmatter, holding exactly one durable
fact, decision, or distilled summary:

```markdown
---
type: decision      # decision|commitment|constraint|fact|score|stakeholder|risk|summary|pattern|insight
date: 2026-06-02
verbatim: true      # cache of always_keep_verbatim; policy is source of truth
phase: discovery    # which phase produced it (optional)
sensitivity: none   # none | sensitive   (sensitive never promotes)
tags: [airgap, ocp] # for relevance-based recall and promotion matching
supersedes: null    # name — or list of names — of record(s) this replaces (no hard deletes)
---
One durable fact / decision / distilled summary. Link related records with
[[other-record-name]].
```

Record `type` vocabulary:

- `decision`, `commitment`, `constraint`, `fact`, `score`, `stakeholder`, `risk`
  — captured in-band by skills during an engagement.
- `summary` — a distilled rollup produced by `/memory compact`.
- `pattern`, `insight` — generalized learnings, typically the *output* of
  `/memory promote` into Tier 2.

**File naming:** `<type>-<kebab-slug>.md` (e.g. `decision-target-platform.md`,
`summary-discovery.md`, `pattern-airgap-registry-mirroring.md`).

## `MEMORY.md` index format (both tiers)

`MEMORY.md` is the recall surface — one line per *active* (un-superseded) record:

```
- [decision] target-platform — OCP 4.16 on bare metal (2026-06-02)
- [constraint] airgap — diode-only data transfer, disconnected registry (2026-06-02)
- [summary] discovery — VMware-heavy, no containers yet; F5 ingress; AD identity (2026-06-03)
```

Format: `- [<type>] <slug> — <one-line hook> (<date>)`. When a record is
superseded, its line drops out of the index (the file stays on disk).

---

## Summarization — how to distill (`/memory compact`)

1. Identify phases older than `keep_recent_phases`, then within them the detail
   records eligible to roll up: those whose `type` is **not** in
   `always_keep_verbatim` **and** that are not sensitive (`sensitivity: sensitive`
   or any `[SENSITIVE]` content is never eligible, regardless of `type`).
2. Write one `summary` record per phase **that has eligible records** (the default
   granularity) capturing the durable substance at the chosen `level`. A phase
   whose records are all verbatim/sensitive (e.g. a `setup` phase) produces no
   summary. At `aggressive`, older `summary` records may themselves be rolled into
   a higher-level summary.
3. Mark each rolled-up detail record as superseded — set the new summary's
   `supersedes:` to the list of record names it rolls up (a single name or a list;
   the detail files remain on disk) — and remove their lines from `MEMORY.md`.
4. **Never touch `CONTEXT.md`.** The audit trail already holds the full history;
   distillation only changes the *recall* layer.
5. A `summary` inherits the highest `sensitivity` of the records it rolls up and
   preserves any `[SENSITIVE]` markers — distillation can never strip a
   sensitivity flag. (Sensitive records are never eligible to roll up in the first
   place, so a clean run never even produces a `sensitivity: sensitive` summary.)

A good summary preserves: what was found, what it means for the engagement, and
pointers (to `CONTEXT.md` sections or artifacts) for anyone who needs the detail.
A bad summary drops numbers, decisions, or named risks — those should have been
`always_keep_verbatim` and must survive.

---

## Promotion — making the fork smarter (`/memory promote`)

Promotion lifts a learning from Tier 1 to Tier 2 so future engagements benefit.

1. **Select candidates.** A record promotes only if the agent judges ALL
   `promotion.criteria` true: reusable across customers, customer-agnostic after
   scrub, and generalizes beyond this engagement. Day-to-day engagement facts do
   not promote — patterns and insights do.
2. **Scrub (mandatory — `sensitivity_scrub: required`).** Before writing
   anything to `memory/records/`:
   - Remove customer name, people, locations, and any other identifying detail.
   - A record with `sensitivity: sensitive` (or any `[SENSITIVE]` content)
     **never promotes** — full stop.
   - Re-read the scrubbed text and confirm it could not identify the customer.
3. **Write** the generalized record (`type: pattern` or `insight`) to
   `memory/records/` and add its line to `memory/MEMORY.md`.

The scrub is the single most important safety control in this system. When in
doubt, do not promote.

---

## Recall — loading memory (`/memory recall` and the Memory Protocol step 1)

- Always load the engagement's `memory/MEMORY.md` index
  (`recall.load_engagement_memory: always`).
- Pull cross-engagement records per `recall.load_cross_engagement`:
  - `by-relevance` (default) — match `tags` / topic to the question, load at most
    `recall.max_cross_records`.
  - `always` — load all Tier-2 records (only sensible for small forks).
  - `off` — engagement memory only.
- Read individual record bodies (and `CONTEXT.md`) only when the index shows
  something relevant. Recall should stay cheap.

---

## Conflict handling

Carried over from the original protocol: when new information contradicts an
existing record, **do not overwrite**. Write a new record and set its
`supersedes:` to the old record's name. Both stay on disk; only the new one
appears in the index. This preserves the chain of how understanding evolved.

---

## Invariants (never violated, regardless of knob settings)

- `CONTEXT.md` is append-only and is never modified by any memory operation.
- Records are never hard-deleted; they are superseded.
- `sensitivity: sensitive` records are never summarized, rolled up, or moved out
  of their engagement directory.
- Everything is local markdown — no external calls, so the whole flow works
  air-gapped.
