---
name: build-deliverable-deck
description: >
  Generate a customer-facing executive presentation from all engagement artifacts.
  Reads CONTEXT.md and all assessment/discovery reports to produce a polished
  deliverable deck. Supports output as Red Hat Quick Deck (HTML), PPTX, or both.
  Use when the architect runs /build-deliverable-deck, says "build the deck",
  "create the presentation", "executive readout", "customer deliverable", or
  "put together the final slides".
---

# /build-deliverable-deck — Deliverable Presentation Generator

This is the capstone skill. It reads all engagement artifacts and produces a polished, customer-facing presentation summarizing findings, recommendations, and next steps.

## Prerequisites

- Engagement initialized (`/setup`)
- At least one discovery or assessment phase completed
- More phases = richer deck. Ideally all relevant phases are done before running this.

## Workflow

### Step 0: Load All Engagement Context

Read everything available:
1. `engagements/<customer>/CONTEXT.md` — full engagement memory
2. `engagements/<customer>/discovery/*.md` — all discovery reports
3. `engagements/<customer>/assessments/*.md` — all assessment reports
4. Any architecture recommendation if it exists

Catalog what's available and what's missing. The deck should only cover phases that have been completed — don't invent content for phases that weren't run.

### Step 1: Choose Output Format

Ask the architect:

> "What format do you want for the deliverable?"
> 1. **Red Hat Quick Deck** (HTML) — Shareable, cinematic, Red Hat branded. Best for screen sharing and modern presentations.
> 2. **PowerPoint (.pptx)** — Traditional format. Best for customers who need to edit or redistribute.
> 3. **Both** — Generate both formats.

**If Red Hat Quick Deck is chosen:**
- Invoke the `red-hat-quick-deck` skill (check for SKILL.md in sibling skill directories or user skills)
- Pass the structured content as the deck narrative

**If PPTX is chosen:**
- Use the `template-slide-deck` skill if a customer/RH template is available
- Otherwise use the `pptx` skill for a clean, professional deck

### Step 2: Structure the Deck Narrative

Build the deck using this engagement-tested story arc:

```
1. TITLE SLIDE
   Customer name, engagement title, date, Red Hat + architect info

2. AGENDA / OVERVIEW
   What we did, what we found, what we recommend

3. ENGAGEMENT SCOPE
   Engagement type, timeline, team, methodology used

4. CURRENT STATE SUMMARY
   Key findings from discovery — infrastructure landscape at a glance
   (Maturity scores visual if infrastructure discovery was run)

5. ASSESSMENT FINDINGS
   Per assessment that was run:
   - Application Portfolio: disposition breakdown, key themes
   - OCP Readiness: scorecard, critical gaps
   - Other assessments: key findings

6. KEY FINDINGS & THEMES
   Cross-cutting themes that emerged across all phases
   (consolidate, don't repeat — this is the synthesis slide)

7. TARGET-STATE ARCHITECTURE
   Architecture recommendation overview (if available from discovery/assessment)
   Deployment topology visual
   Key design decisions

8. IMPLEMENTATION ROADMAP
   Phased plan visual (if available)
   Migration waves (if app assessment was run)
   Timeline and milestones

9. INVESTMENT SUMMARY
   Red Hat product BOM (if available)
   Resource requirements
   Effort estimates

10. RISKS & MITIGATIONS
    Top 3-5 risks with mitigations

11. RECOMMENDED NEXT STEPS
    Immediate actions (next 30 days)
    Near-term plan (60-90 days)
    Success criteria and how to measure them

12. APPENDIX (optional)
    Detailed data tables
    Methodology notes
    Full assessment references
```

**Adaptive slide selection:** Only include slides for phases that were completed. If only discovery was done, the deck focuses on findings and recommendations for next steps — don't leave blank sections.

### Step 3: Draft Content Per Slide

For each slide, draft the content:

- **Headlines should tell the story.** Not "Application Portfolio" but "73% of Applications Are Modernization Candidates"
- **Limit text per slide.** 3-5 bullet points max, each under 15 words
- **Use data where available.** Maturity scores, disposition percentages, readiness verdicts
- **Executive-level language.** No deep technical jargon on main slides — save that for appendix
- **Every slide earns its place.** If a slide doesn't advance the narrative or influence a decision, cut it

### Step 4: Review With Architect

Present the deck outline to the architect before generating:

```
Here's the deck structure I'm proposing:

1. Title: "<Engagement Title>"
2. Agenda
3. Scope: <engagement type>, <dates>
4. Current State: <1-line summary of infrastructure>
5. Findings: <key finding headlines>
6. Architecture: <recommendation summary>
7. Roadmap: <n> phases, <timeline>
8. Investment: <product summary>
9. Risks: <top 3>
10. Next Steps: <key actions>

Want me to add, remove, or change anything before I generate the deck?
```

Incorporate feedback before generating.

### Step 5: Generate the Deck

Invoke the appropriate skill:

**For Red Hat Quick Deck:**
- Pass the structured content to the quick deck skill
- Ensure cinematic narrative flow, dark-mode aesthetics, Red Hat brand
- Include speaker notes with talking points for each slide

**For PPTX:**
- Use template if available, otherwise clean professional layout
- Red Hat color palette (red #EE0000, dark gray #333, white)
- Include speaker notes

### Step 6: Save and Report

Save the deliverable to `engagements/<customer>/deliverables/`:
- `engagement-readout.html` (Quick Deck)
- `engagement-readout.pptx` (PowerPoint)

Update `CONTEXT.md`:

```markdown
## Phase: Deliverable Generation
**Date:** <today>
**Format:** <format chosen>

### Artifacts Produced
- `deliverables/engagement-readout.<ext>` — Customer-facing executive presentation

### Slides Covered
- <list of slide topics included>

### Engagement Status
- Discovery: ✅ Complete
- Assessment(s): ✅ Complete / ⬜ Not run
- Architecture: ✅ Complete / ⬜ Not run
- Deliverable: ✅ Generated
```

### Step 7: Suggest Final Actions

- "Want me to generate a summary email to send to the customer with the deck attached?"
- "Should I create an ADR for the key architecture decisions in this engagement?"
- "Need a one-page executive summary document alongside the deck?"
