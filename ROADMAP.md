# RHEK Roadmap — The Compounding Delivery Engine

**A north star for the Red Hat Engagement Kit (RHEK).**
Status: vision · living document · Last updated: 2026-06-19

> **RHEK's north star:** every engagement leaves behind not just a *deliverable*,
> but *capability* and *knowledge* that make the next engagement faster.
> Automatically. Safely. Across the whole field. The 500th engagement should be
> dramatically easier than the 5th — and no architect should ever start from a
> blank repo again.

This is a dream, not a Gantt chart. The horizons below are the wheel gaining
speed, not calendar dates. We pull threads in whatever order builds momentum.

---

## Where we are now

RHEK today is a **brilliant scribe**. It's a fork-and-own, markdown-native,
dual-runtime (Claude Code + OpenCode), air-gap-friendly kit. A **default set** of
interchangeable skills — `setup → discover → assess → deck` — walks an architect
through an engagement and leaves behind clean artifacts; fork the kit and swap
them for whatever your engagement demands. And we just built **configurable,
tiered memory**: an append-only audit trail, a distilled recall layer, a
cross-engagement institutional tier, and a `promote` step that scrubs customer
detail before anything is shared.

Two honest truths about today:

1. **The scribe documents; it doesn't *do*.** The skills interview, score, and
   present. They don't execute the actual technical delivery.
2. **Every fork is an island.** Promoted knowledge exists — but it's trapped
   inside the fork. Tier 2 is built and empty. The flywheel has a hub and one
   spoke.

Everything below turns the scribe into a *doer*, and connects the islands into a
*commons* — so that knowledge **and** capability compound together.

---

## The skill model — bring your own skills

RHEK isn't a four-step wizard wearing a roadmap. It's a **platform for skills** —
and the four that ship are a *default walk*, not the product. Each skill is a
self-contained unit that reads the engagement's shared context, does its job, and
honors the lightweight Memory Protocol. Fork the kit, delete what you don't need,
and drop in the skills your engagement actually demands — an RHOAI install, a
security review, a migration plan, a customer-specific checklist.

Only two pieces are load-bearing:

- **`/setup` — the one true prerequisite.** It sets the project start: the
  engagement directory, the metadata, the `CONTEXT.md` everything else reads from.
  Nothing meaningful happens before it.
- **`/memory` — the foundational engine.** It owns the heavy memory work —
  compact, promote, recall — so the phase skills stay swappable without breaking
  the flywheel. Every skill honors the Memory Protocol; `/memory` is what makes
  that protocol pay off.

Everything else — `discover → assess → deck` and whatever you add — is the
**swappable layer**. `setup → discover → assess → deck` is just the path most
engagements happen to walk first. The point of a fork-and-own kit is that *your*
fork carries *your* skills. The roadmap below assumes exactly that: when we say
RHEK "becomes a doer" (H2) or skills "harden themselves" (H4), we mean *whatever
skills your fork runs* — not a fixed roster of four.

---

## The engine

The whole roadmap is one loop, getting faster:

```
   ┌─▶ RUN ───────▶ REMEMBER ───────▶ DISTILL & SCRUB ───────▶ FEDERATE ─┐
   │   execute      capture what      promote customer-        publish a  │
   │   the build    actually          agnostic learnings       scrubbed   │
   │   (not just    happened          (the safety valve)       knowledge  │
   │   document)    worked / broke                             pack       │
   │                                                                      │
   └──────────────────────────  RECALL  ◀─────────────────────────────────┘
            the next engagement starts pre-loaded from the commons
```

The moat is the word **together**. Most knowledge management compounds only the
docs (Remember → Distill). RHEK compounds the *doing* too — the execution skills
themselves get hardened and shared. Knowledge and capability, on the same wheel.

---

## The four horizons

```
  H1 ─────────────▶ H2 ───────────▶ H3 ─────────────▶ H4
  REMEMBER          DO              FEDERATE          COMPOUND
  scribe that       curated         islands become    every engagement
  never forgets     execution       a knowledge       starts ahead
  (SHIPPED)         skills          commons           (moonshot)
                    (RHOAI)         (pack + scrub)

           ↻ each turn of the wheel feeds the next ↻
```

*Horizon order is the order we build the capabilities — the engine itself is a cycle with no first step.*

### H1 · Remember — the scribe that never forgets · *Shipped*

Configurable, tiered memory governs *what to keep verbatim vs. distill for
recall*. `CONTEXT.md` stays the append-only audit trail; a separate recall layer
stays lean; a cross-engagement tier (`memory/`) is the seed of the commons;
`promote` + a mandatory scrub is the safety valve. `recall` is deliberately
designed as the grounding layer a leave-behind assistant calls.

- **Today:** built and merging. The wheel exists but barely turns — Tier 2 is
  empty, `promote` has never run on a live engagement.
- **The remaining job is small and symbolic — spin it once:** run `promote` on a
  real engagement, watch a learning land in Tier 2, prove the scrub holds.
- **Signal it's working:** the first scrubbed learning appears in Tier 2.

### H2 · Do — the scribe becomes a doer · *Next*

Curate **executable delivery skills** and drive them from the Red Hat side, not
just on the customer's site. The beachhead is **John Hurlocker's RHOAI install
skill** — the first RHEK skill that *executes* instead of *documents*.

Why this matters beyond automation: a discovery interview produces *opinions*; a
real install produces *evidence*. The moment RHEK runs the build, the memory it
captures becomes **"what actually worked / what actually broke."** **Do feeds
Federate** — execution is what generates knowledge worth federating.

- **Today:** designed direction; Hurlocker's skill is the first to port.
- **Signal it's working:** an architect completes an RHOAI install *through*
  RHEK, and the run emits worked/broke telemetry into memory.

### H3 · Federate — the islands become a commons · *Designed*

How does distilled knowledge escape a fork that is *deliberately sealed*? The
fork is a private box — today's `pre-push` hook blocks pushes outright, precisely
so customer data can't leak. The challenge isn't "sync"; it's **cutting one safe,
one-way hole in a box we welded shut.**

The answer — **the knowledge pack as the unit of exchange:**

- `promote` bundles scrubbed, customer-agnostic records into a portable,
  **signable knowledge pack** that can cross an air-gap diode. *Air-gap-first.*
- When connected, the *same pack* simply rides git as an upstream pull request.
  **One artifact, two transports.**
- A second fork **imports** the pack (or pulls the upstream commons) and is now
  smarter than it was.

- **Today:** model locked; not yet built. The first concrete task is the
  **`pre-push` hook carve-out** — keep blocking the engagement directories, but
  permit *only* the curated commons-contribution path. The box gets one door.
- **Signal it's working:** a knowledge pack crosses from one fork into the
  commons and lands in a *different* fork.

### H4 · Compound — the flywheel at full speed · *Horizon (moonshot)*

Where network effects take over and RHEK stops being a *tool* and becomes
*infrastructure*:

- **Cold-start dies.** A new "air-gapped RHOAI for federal" engagement opens
  pre-loaded — the commons recognizes the *shape* and seeds the fork with proven
  patterns, hardened skills, checklists, and known gotchas. `recall` becomes
  **pre-load.** Engagement #500 starts where #499 left off.
- **Skills harden themselves.** Every run of a shared execution skill emits
  *worked / broke / edge-case.* Aggregated, the skill *improves* — the commons
  learns "step 4 fails on disconnected registries ~30% of the time; here's the
  fix." **Capability compounds, not just knowledge.**
- **Trust is earned by evidence.** A pattern corroborated by nine engagements
  outranks a one-off. The commons is *ranked by provenance* — it doesn't rot like
  a wiki.
- **The moat (for leadership).** When an architect moves on, their hard-won
  delivery knowledge doesn't walk out the door — it's in the commons. Red Hat's
  delivery capability **compounds instead of resetting with turnover.**
- **The kicker — it feeds the product.** Aggregated "what actually breaks in the
  field" is signal gold for the RHOAI / OCP product teams. The *delivery*
  flywheel starts feeding the *product* flywheel.
- **Signal it's working:** cold-start time drops measurably; a shared skill
  self-corrects from field telemetry; top patterns carry real corroboration
  counts.

---

## The framing: open source, applied to delivery knowledge

The commons is not a new idea — it's Red Hat's *oldest* idea, pointed at a new
target. The entire model is **open source applied to how we deliver:**

| Open source (code)        | RHEK (delivery knowledge)                  |
|---------------------------|--------------------------------------------|
| Fork                      | The engagement fork                        |
| Downstream                | One customer engagement                    |
| Upstream                  | The shared commons                         |
| Prepare a patch           | `promote` + scrub into a knowledge pack    |
| Submit upstream (PR)      | Contribute the pack (over git, when connected) |
| Maintainer review + merge | The human gate on the commons              |
| Rebase on upstream        | Import the commons / re-fork               |

Red Hat already lives this ritual for *software*. RHEK extends it to *how we
deliver* — a flag only Red Hat can credibly plant.

---

## The safety model: one safe hole in a sealed box

Federation must never weaken the thing that makes forks safe. The fork→commons
edge is a **one-way valve with two locks**:

1. **The automated scrub** (already built) — strips customer-identifying detail;
   anything tagged `sensitive` *never* leaves the engagement directory.
2. **The human gate** (new) — a maintainer reviews a pack before it joins the
   commons.

Invariants, never violated:

- **Never an automatic sync** from a sensitive fork. Contribution is always a
  deliberate, reviewable, outbound-only act.
- The `pre-push` hook keeps blocking engagement data; only the curated pack path
  gets a door.
- **Air-gap-first:** the pack is a file that can cross an approved transfer
  boundary; nothing in the flow requires a live service.

---

## The connective tissue (already in flight)

The dream isn't vapor — the pieces are being built in parallel:

- **The face — the OpenTUI workbench + GUI lifecycle.** Where the architect
  *operates* the flywheel: run skills, browse artifacts, host the leave-behind
  assistant — and, newly, **the human gate lives here**: review what gets
  promoted, sign a knowledge pack, import the commons. The cockpit.
- **Commons gardening.** A maintainer model, dedupe/merge of similar patterns,
  staleness review, and provenance counts — the OSS governance that keeps the
  commons healthy.
- **The Idea Lab as the R&D engine.** A research loop that already dreamed this
  (its first idea was the memory flywheel → shared commons). The continuous
  idea-generator that feeds this roadmap.
- **Role-based checklists.** The engagement's readiness fabric (PM / architect /
  consultant / kickoff), first-class in the workbench.

---

## What this asks of you

**Architects — the flywheel only spins if you push knowledge back.**
Fork it. Run a real engagement. When you learn something that would have saved you
a week, `promote` it and ship the pack. Port the execution skill you keep redoing
by hand. The commons is only as smart as what the field contributes.

**Leadership — back the commons as shared infrastructure.**
The ask is small: a git repo, a maintainer model, a review gate. The return
compounds: faster delivery, a knowledge asset that survives turnover, and field
signal that feeds the products. This is how Red Hat's delivery capability stops
resetting and starts compounding.

---

*RHEK: every engagement makes the next one faster.*
