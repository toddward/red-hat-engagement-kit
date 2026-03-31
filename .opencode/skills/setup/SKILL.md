---
name: setup
description: >
  Initialize a new customer engagement. Creates the engagement directory structure,
  gathers core engagement metadata through a structured interview, and writes the
  initial CONTEXT.md that all subsequent skills will build upon. Run this first —
  it's the foundation everything else reads from.
---

# /setup — Initialize Engagement

This skill bootstraps a new customer engagement by gathering essential metadata and creating the engagement workspace. It produces the `CONTEXT.md` file that every downstream skill depends on.

## When This Skill Triggers

- Architect runs `/setup` in OpenCode
- Architect says "start a new engagement", "initialize engagement", or "set up for <customer>"
- The `engagements/` directory has no customer subdirectories yet

## Workflow

### Step 1: Gather Engagement Metadata

Conduct a structured interview with the architect. Ask these questions conversationally — don't dump them all at once. Group related questions and adapt based on answers.

**Required Fields (must have before proceeding):**

| Field | Question | Example |
|-------|----------|---------|
| `customer_name` | "What's the customer name?" | Acme Federal |
| `customer_slug` | Auto-generate from name (kebab-case) | `acme-federal` |
| `engagement_type` | "What type of engagement is this?" | See engagement types below |
| `architect_name` | "Who's the lead architect?" | Todd Wardzinski |
| `engagement_dates` | "What are the start and target end dates?" | 2026-04-01 → 2026-06-30 |
| `executive_sponsor` | "Who's the customer executive sponsor?" | Jane Smith, CTO |

**Engagement Types** (offer these as options):

1. **Application Modernization** — Assess and migrate legacy applications to cloud-native platforms
2. **AI/ML Enablement** — Evaluate and deploy AI/ML capabilities on Red Hat infrastructure
3. **Platform Assessment** — Evaluate readiness for OpenShift, RHEL, or hybrid cloud adoption
4. **Security & Compliance** — Assess security posture and compliance alignment
5. **Infrastructure Modernization** — Datacenter transformation, edge computing, hybrid cloud
6. **Custom** — Architect defines the engagement scope manually

**Optional Fields (ask but don't block on):**

| Field | Question |
|-------|----------|
| `clearance_level` | "Any clearance requirements? (Unclassified, CUI, Secret, TS/SCI)" |
| `environment_type` | "Connected, limited connectivity, or air-gapped?" |
| `team_members` | "Who else is on the engagement team?" |
| `customer_industry` | "What sector? (Federal Civilian, DoD, IC, SLED, Critical Infrastructure)" |
| `existing_rh_footprint` | "Any existing Red Hat products in the environment?" |
| `key_constraints` | "Any known constraints? (budget, timeline, compliance, legacy dependencies)" |
| `success_criteria` | "What does success look like for this engagement?" |

### Step 2: Create Engagement Directory

Once you have the required fields, create the directory structure:

```bash
CUSTOMER_SLUG="<customer_slug>"
mkdir -p "engagements/${CUSTOMER_SLUG}"/{discovery,assessments,deliverables}
```

### Step 3: Write CONTEXT.md

Generate the initial `CONTEXT.md` at `engagements/<slug>/CONTEXT.md` using this template:

```markdown
# Engagement Context: <Customer Name>

> This file is the living memory for this engagement. Every skill reads from and
> appends to this file. Do not manually edit unless correcting an error.

## Engagement Metadata
- **Customer:** <customer_name>
- **Type:** <engagement_type>
- **Lead Architect:** <architect_name>
- **Dates:** <start_date> → <end_date>
- **Executive Sponsor:** <executive_sponsor>
- **Clearance:** <clearance_level or "Unclassified">
- **Environment:** <environment_type or "Connected">
- **Industry:** <customer_industry or "Not specified">
- **Existing RH Footprint:** <existing_rh_footprint or "Unknown — to be discovered">

## Team
- <architect_name> (Lead Architect)
<additional team members if provided>

## Constraints
<key_constraints or "None identified yet">

## Success Criteria
<success_criteria or "To be defined during discovery">

## Engagement Log
### <today's date> — Engagement Initialized
- Engagement workspace created by <architect_name>
- Type: <engagement_type>
- Status: **Discovery phase**

---
<!-- Skills append their findings below this line -->
```

### Step 4: Write .gitignore

Create or update the repo-level `.gitignore` to protect sensitive engagement data:

```
# Engagement artifacts that should never be pushed to shared remotes
engagements/*/discovery/*.sensitive
engagements/*/discovery/*.classified

# OS artifacts
.DS_Store
Thumbs.db

# Editor artifacts
*.swp
*.swo
*~
```

### Step 5: Confirm and Recommend Next Steps

After creating the workspace, confirm what was created and recommend the logical next skill:

```
✅ Engagement initialized: <customer_name>
📁 Workspace: engagements/<slug>/
📄 Context: engagements/<slug>/CONTEXT.md

Recommended next steps based on engagement type:
```

**Next-step recommendations by engagement type:**

| Engagement Type | Recommended Next Skill |
|-----------------|----------------------|
| Application Modernization | `/discover-infrastructure` then `/assess-app-portfolio` |
| AI/ML Enablement | `/discover-infrastructure` then `/assess-app-portfolio` |
| Platform Assessment | `/discover-infrastructure` then `/assess-app-portfolio` |
| Security & Compliance | `/discover-infrastructure` with security focus |
| Infrastructure Modernization | `/discover-infrastructure` |
| Custom | `/discover-infrastructure` (always a good starting point) |

## Error Handling

- If an engagement directory already exists for the given slug, ask: "An engagement for `<slug>` already exists. Do you want to resume it, or create a new one with a different name?"
- If the architect provides minimal information, create the workspace with what you have and note gaps in CONTEXT.md under a `## Known Gaps` section.
- If running in an air-gapped environment, confirm that all knowledge base files are present locally.

## Example Session

```
Architect: /setup

OpenCode: Let's get this engagement started. What's the customer name?

Architect: Acme Federal Services

OpenCode: Got it — Acme Federal Services. What type of engagement is this?
  1. Application Modernization
  2. AI/ML Enablement
  3. Platform Assessment
  4. Security & Compliance
  5. Infrastructure Modernization
  6. Custom

Architect: 2 - AI/ML Enablement

OpenCode: Who's the lead architect for this engagement?

Architect: Me — Todd Wardzinski

OpenCode: What are the start and target end dates?

Architect: Starting next week, targeting end of Q2

OpenCode: And who's the customer executive sponsor?

Architect: Deputy Director of IT, Mark Johnson

OpenCode: A few more optional details that'll help downstream —
Any clearance requirements?

Architect: TS/SCI

OpenCode: Environment type — connected, limited connectivity, or air-gapped?

Architect: Air-gapped

OpenCode: Any existing Red Hat products in the environment?

Architect: RHEL 8 across the board, some Ansible Tower

OpenCode: ✅ Engagement initialized: Acme Federal Services
📁 Workspace: engagements/acme-federal-services/
📄 Context: engagements/acme-federal-services/CONTEXT.md

Since this is an AI/ML Enablement engagement in an air-gapped TS/SCI environment,
I'd recommend starting with `/discover-infrastructure` to map the compute and
data landscape, with special attention to GPU resources and data classification levels.
```