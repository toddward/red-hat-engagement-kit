# Memory — Scenario Evals

The memory system is instructions + markdown, so its "tests" are scenario
walkthroughs an architect (or the `skill-creator` eval harness) can run against a
fixture and check by eye or by simple diff.

`sample-engagement/` is a small, **fixture** engagement (not a real customer)
seeded with three phases of `CONTEXT.md` plus a pre-compaction `memory/` recall
layer. Each eval below states a starting condition, an action, and the assertion.

> The fixture lives under `docs/` on purpose — it must never sit in the real
> `engagements/` tree, which is customer-sensitive space.

## How to run

1. Copy `sample-engagement/` to a scratch location (so you can re-run from a
   clean state).
2. Point the `/memory` skill at the copy and perform the action.
3. Check the assertion.

The repo-level `memory/POLICY.md` supplies the knobs; some evals ask you to tweak
a knob on the copy first.

---

## E1 — Compaction never touches the audit trail

- **Given** the fixture (work phases discovery → assessment → architecture;
  `discovery` is older than the protected `keep_recent_phases: 2` window).
- **When** you run `/memory compact`.
- **Then** `sample-engagement/CONTEXT.md` is **byte-for-byte unchanged**
  (`git diff --stat` shows no change to it). Only `memory/` files change.

## E2 — Old detail becomes a lean summary

- **Given** the fixture, where `fact-vmware-estate` and `fact-airgap-ops-untested`
  are non-verbatim details in the older `discovery` phase.
- **When** you run `/memory compact`.
- **Then** a single `summary-discovery` record appears, both discovery facts are
  marked superseded and **drop out of `memory/MEMORY.md`**, and the index is
  shorter than before (8 → 7 lines). The superseded files still exist on disk (no
  hard delete).

## E3 — Verbatim types AND sensitive records survive compaction

- **Given** verbatim-by-type records `score-discovery-maturity`, `risk-skills-gap`,
  `stakeholder-sponsor`, `decision-target-platform`, and the
  verbatim-by-sensitivity record `constraint-clearance` (`type: constraint`,
  `sensitivity: sensitive`, in the oldest `setup` phase).
- **When** you run `/memory compact`.
- **Then** none are summarized away — all five still appear in `MEMORY.md` with
  their original content intact. In particular `constraint-clearance` survives
  even though its `type` (`constraint`) is **not** in `always_keep_verbatim`,
  because the sensitivity floor (`keep_sensitive_verbatim: true`) protects it.
  This is the regression guard for the type-vs-sensitivity gap.

## E4 — Promotion is scrubbed and generalized

- **Given** `fact-airgap-registry-approach` (a customer-agnostic, generalizable
  learning — the oc-mirror → diode → Quay disconnected-install path) and
  `constraint-clearance` (`sensitivity: sensitive`).
- **When** you run `/memory promote`.
- **Then** a customer-agnostic `pattern`/`insight` record (e.g.
  `pattern-airgap-registry-mirroring`) is written to the repo-level
  `memory/records/` with **no customer name** ("Acme Federal") and no
  `[SENSITIVE]` content; `constraint-clearance` is **not** promoted; the skill
  reports what it withheld and why.

## E5 — The `level` dial measurably changes behavior

- **Given** two copies of the fixture (the `discovery` phase holds two
  summarizable details: `fact-vmware-estate` and `fact-airgap-ops-untested`).
- **When** you set `summarization.level: conservative` on one and `aggressive` on
  the other, then run `/memory compact` on each.
- **Then** the two indexes differ in **composition** (the robust signal): the
  `aggressive` index has a `[summary] discovery` line and neither `[fact]
  vmware-estate` nor `[fact] airgap-ops-untested`, while `conservative` keeps at
  least one of those detail facts and produces fewer (or no) `summary` lines. In
  line count `aggressive` is shorter than or equal to `conservative` (it rolls
  both discovery facts into one summary; `conservative` keeps more verbatim).

## E6 — Missing policy falls back to balanced

- **Given** a copy of the fixture with the repo-level `memory/POLICY.md`
  temporarily renamed/removed.
- **When** you run any `/memory` operation.
- **Then** it proceeds using the documented `balanced` defaults **and states in
  its output that it fell back to defaults** because no `POLICY.md` was found.

---

## Invariants checked across all evals

- `CONTEXT.md` is never modified by a memory operation (E1).
- No record is ever hard-deleted — only superseded (E2).
- `sensitivity: sensitive` records never leave the engagement directory (E4).
- All actions are local; nothing reaches the network (air-gap safe).
