---
name: assess-app-portfolio
description: >
  Collect system information by running a local discovery script. Demonstrates
  how a skill can execute a binary or script behind the scenes to gather data
  that feeds into the engagement assessment. Use when the architect runs
  /assess-app-portfolio, mentions "app assessment", "system info", or
  "collect system data".
---

# /assess-app-portfolio — System Information Collection

This skill runs a local bash script to collect system information from the machine, then writes the results into the engagement's assessment directory. It demonstrates how skills can execute binaries or scripts behind the scenes to gather data automatically.

## Prerequisites

- Engagement initialized (`/setup`) — need to know the customer directory name

## Workflow

### Step 0: Load Context

Read `engagements/<customer>/CONTEXT.md` to determine the customer name and engagement directory.

If multiple engagement directories exist under `engagements/`, ask the architect which one to use.

### Step 1: Run the Collection Script

Execute the system information collection script bundled with this skill:

```bash
bash .claude/skills/assess-app-portfolio/collect-system-info.sh
```

Capture the full output — this is the raw discovery data.

### Step 2: Present Results

Show the architect the collected system information and highlight anything notable:
- OS and architecture details
- Available container runtimes (Podman, Docker)
- Kubernetes/OpenShift CLI availability
- Detected runtime environments (Java, Python, Node, Go, etc.)
- Disk and memory capacity

Ask if there's anything they'd like to add or correct.

### Step 3: Write Assessment Report

Write the script output and any architect annotations to:
`engagements/<customer>/assessments/system-info-collection.md`

Format:

```markdown
# System Information Collection
**Customer:** <name>
**Date:** <today>
**Collected by:** automated script + architect review

## Raw Collection Output
<paste full script output here>

## Architect Notes
<any corrections or additions from the architect>

## Summary
<brief summary of the system landscape based on collected data>
```

### Step 4: Update CONTEXT.md

Append to the engagement's `CONTEXT.md`:

```markdown
## Phase: System Information Collection
**Date:** <today>

### System Landscape
- OS: <detected>
- Architecture: <detected>
- Container Runtime: <detected or none>
- Kubernetes/OpenShift CLI: <detected or none>
- Runtime Environments: <list detected>
- Memory: <detected>
- Disk: <detected>

### Artifacts Produced
- `assessments/system-info-collection.md` — Full system info report
```

### Step 5: Recommend Next Steps

- If customer wants a presentation → `/build-deliverable-deck`
