# Engagement Context: Acme Federal  (FIXTURE — not a real customer)

> This file is the append-only audit trail for this engagement — every skill
> appends and never rewrites it. The distilled recall layer lives in `memory/`
> (see the Memory Protocol). Do not manually edit unless correcting an error.

## Engagement Metadata
- **Customer:** Acme Federal
- **Type:** Platform Assessment
- **Lead Architect:** Pat Architect
- **Dates:** 2026-05-04 → 2026-07-31
- **Executive Sponsor:** Dana Sponsor, Deputy CIO
- **Clearance:** [SENSITIVE] Secret
- **Environment:** Air-gapped
- **Industry:** Federal Civilian
- **Existing RH Footprint:** RHEL 8 fleet, some Ansible

## Engagement Log
### 2026-05-04 — Engagement Initialized
- Engagement workspace created by Pat Architect
- Type: Platform Assessment
- Status: **Discovery phase**

---
<!-- Skills append their findings below this line -->

## Phase: Infrastructure Discovery
**Date:** 2026-05-11
**Conducted by:** Pat Architect

### Environment Summary
Workloads run on an on-prem VMware estate (~480 VMs, vSphere 7) with no container
platform yet. Ingress via F5; identity via Active Directory. The environment is
air-gapped, with updates delivered by a one-way data diode and a partial
disconnected registry mirror. Capacity is request-based with no self-service.

### Maturity Scores
- Compute & Hosting: 3
- Networking & Connectivity: 2
- Storage & Data: 3
- Identity & Access: 3
- CI/CD & Developer Tooling: 2
- Operations & Observability: 2

### Key Findings
- No container platform; greenfield for OpenShift.
- Air-gap operations are mature for RHEL but untested for a registry-heavy
  platform.
- Container/Kubernetes skills are thin across the ops team.

### Artifacts Produced
- `discovery/infrastructure-discovery.md` — Full discovery report

## Phase: Platform Readiness Assessment
**Date:** 2026-05-26
**Conducted by:** Pat Architect

### Summary
Confirmed OpenShift 4.16 on bare metal as the target platform. Disconnected
install is the critical path; the diode + mirror approach needs a registry
mirroring runbook. Skills gap is the top delivery risk.

### Decisions
- Target platform: OpenShift 4.16 on bare metal (over VMware-hosted) — better
  fit for the eventual edge sites and avoids a hypervisor tax.

### Artifacts Produced
- `assessments/platform-readiness.md` — Readiness scorecard

## Phase: Target Architecture
**Date:** 2026-06-01
**Conducted by:** Pat Architect

### Summary
Proposed a hub-and-spoke OpenShift topology with a disconnected registry mirror
fed over the existing data diode. The disconnected install runbook (mirror →
diode → internal registry) is the critical path; it is a reusable approach for
any air-gapped OpenShift build.

### Decisions
- Adopt a disconnected registry mirroring runbook (oc-mirror → diode → Quay) as
  the standard install path for this and future air-gapped sites.

### Artifacts Produced
- `deliverables/target-architecture.md` — Topology + disconnected install runbook
