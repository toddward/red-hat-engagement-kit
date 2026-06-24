---
name: memory
description: >
  Maintain the engagement's long-lived memory. The foundational, skill-agnostic
  engine behind the Memory Protocol: compact (distill old detail into summaries),
  promote (lift scrubbed, generalizable learnings into cross-engagement memory),
  recall (answer from engagement + institutional memory), and status (report
  memory health). Reads memory/POLICY.md to decide what to keep verbatim vs.
  summarize. Use when the architect runs /memory, says "compact memory",
  "summarize what we know", "promote these learnings", "what did we learn about
  X", "what do we know so far", or "memory status". Also the grounding layer the
  leave-behind assistant calls to answer questions from engagement artifacts.
---

# /memory — Long-Lived Memory Engine

This is the foundational skill behind the kit's memory system. It does the heavy,
judgment-driven work that the lightweight **Memory Protocol** (in `CLAUDE.md`)
deliberately leaves out: distilling detail for efficient recall, promoting
learnings across engagements, and answering questions grounded in memory.

It is **skill-agnostic** — it knows nothing about discovery, assessments, or
decks. Domain skills are interchangeable per project; this skill and the protocol
are what stay constant. Every behavior here is governed by `memory/POLICY.md`.

## Operating Principles (read before any operation)

1. **Read the policy first.** Load `memory/POLICY.md`. If it is missing, use the
   documented `balanced` defaults and state that you fell back to defaults.
2. **The audit trail is sacrosanct.** Never read-modify-write `CONTEXT.md`. You
   may *read* it; you must never edit or delete from it. If asked to "clean up"
   or "shorten" CONTEXT.md, decline and explain that distillation happens in the
   `memory/` recall layer instead.
3. **Never hard-delete a record.** Supersede it (write a replacement whose
   `supersedes:` points at the old record; drop the old line from `MEMORY.md`).
4. **Everything is local.** No external calls — the flow must work air-gapped.

## Which Engagement?

Most operations act on one engagement. Determine the target:

- If exactly one directory exists under `engagements/` (excluding `.template`),
  use it.
- If several exist, ask the architect which customer.
- `recall` may *additionally* read the repo-level `memory/` (Tier 2); `promote`
  reads one engagement and writes to Tier 2.

---

## Operation: `compact` — distill for lean recall

Triggered by: "compact memory", "summarize what we've got", "tidy up memory",
running it at a phase boundary, or `/memory compact`.

1. Load `memory/POLICY.md` (`level`, `keep_recent_phases`, `always_keep_verbatim`).
2. Read the engagement's `memory/MEMORY.md` index and `CONTEXT.md` (read-only) to
   understand phase boundaries and what detail exists.
3. Identify candidates: records in phases older than `keep_recent_phases` that are
   eligible to roll up — i.e. their `type` is **not** in `always_keep_verbatim`
   **and** they are not sensitive. A record with `sensitivity: sensitive` or any
   `[SENSITIVE]` content is **never** a candidate, regardless of `type` (the same
   predicate `promote` uses).
4. For each such phase, write one `summary` record (`type: summary`) capturing the
   durable substance at the configured `level`. Preserve: what was found, what it
   means, and pointers back to `CONTEXT.md`/artifacts. A `summary` inherits the
   highest `sensitivity` of its sources and preserves any `[SENSITIVE]` markers.
   At `aggressive`, older `summary` records may be rolled into a higher-level
   summary.
5. Mark the rolled-up detail records as superseded — set the new summary's
   `supersedes:` to the list of record names it rolls up (a single name or a
   list); leave the files on disk — and remove their lines from `MEMORY.md`. Add
   the new `summary` line(s).
6. **Do not modify `CONTEXT.md`.**
7. Report: which phases were distilled, how many records were superseded, and any
   records you flagged as promotion candidates.

**Verbatim guard:** a record is never collapsed if its `type` is in
`always_keep_verbatim` **or** it is sensitive (`sensitivity: sensitive` /
`[SENSITIVE]`). Re-evaluate against the *current* policy (the type list **and**
the sensitivity field), not the cached `verbatim` flag.

---

## Operation: `promote` — make the fork smarter

Triggered by: "promote these learnings", "save this for future engagements",
"add this to institutional memory", or `/memory promote`.

1. Confirm `promotion.enabled` in the policy. If disabled, say so and stop.
2. Scan the engagement's records for candidates. A record promotes only if ALL
   `promotion.criteria` hold: reusable across customers, customer-agnostic after
   scrub, generalizes beyond this engagement.
3. **Mandatory sensitivity scrub** (`sensitivity_scrub: required`):
   - Any record with `sensitivity: sensitive` or `[SENSITIVE]` content **never
     promotes**.
   - Strip customer name, people, locations, and any identifying detail from the
     candidate.
   - Re-read the scrubbed text and confirm it cannot identify the customer. If
     you are not sure, do not promote.
4. Write the generalized record to `memory/records/` as `type: pattern` or
   `insight`, and add its line to `memory/MEMORY.md` (Tier 2).
5. Report what was promoted (generalized titles only) and what was withheld and
   why.

> The scrub is the most important safety control in the system. It is also why
> Tier 2 is safe to keep across forks while engagement directories are not.

---

## Operation: `recall` — answer from memory

Triggered by: "what do we know about X", "what did we learn", "recall …",
questions from the leave-behind assistant, or `/memory recall <topic>`.

1. Always load the engagement's `memory/MEMORY.md` index
   (`recall.load_engagement_memory: always`).
2. Pull cross-engagement records per `recall.load_cross_engagement`
   (`by-relevance` default → match `tags`/topic, cap at `max_cross_records`).
3. Read individual record bodies (and `CONTEXT.md` sections) only when the index
   shows something relevant — keep recall cheap.
4. Answer the question, citing which records/artifacts the answer rests on
   (engagement vs. institutional). If nothing relevant exists, say so plainly.

This is the operation the GUI/TUI leave-behind assistant calls to ground its
answers in engagement and institutional memory.

---

## Operation: `status` — memory health

Triggered by: "memory status", "how big is memory", or `/memory status`.

Report, cheaply (index-level, no deep reads):

- Audit trail size (`CONTEXT.md` lines) vs. active recall size (`MEMORY.md`
  lines) — the compression you're getting.
- Record counts per `type`, and how many are superseded.
- Last compaction (most recent `summary` record date) and current `level`.
- Pending promotion candidates (records that look generalizable but haven't been
  promoted).

---

## Edge Cases

- **Missing `POLICY.md`** → use documented `balanced` defaults; state the fallback.
- **No memory yet** (new/empty engagement) → `recall`/`status` report "no memory
  yet"; `compact` is a no-op with a note.
- **Multiple engagements** → ask which customer before acting.
- **Asked to edit `CONTEXT.md`** → refuse; explain the audit-trail invariant and
  redirect to the records layer.
- **Air-gapped** → all operations are local file reads/writes; never reach out.
- **Promotion with only sensitive candidates** → promote nothing; explain why.

## Relationship to the Memory Protocol

The Memory Protocol (`CLAUDE.md`) is the cheap, in-band half that *every* skill
performs: recall first, append to the audit trail, capture durable records,
flag-don't-overwrite, and recommend `/memory compact` at phase boundaries. This
skill is the expensive, on-demand half. Together they are the "hybrid" trigger
model: automatic-feeling capture, controlled distillation and promotion.
