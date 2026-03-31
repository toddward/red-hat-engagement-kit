---
name: discover-infrastructure
description: >
  Conduct a structured infrastructure discovery interview for a customer engagement.
  Walks the architect through compute, networking, storage, identity, CI/CD, and
  operational tooling topics. Produces a discovery report and updates CONTEXT.md.
  Use when the architect runs /discover-infrastructure, says "let's do discovery",
  "map the infrastructure", or "what's in their environment".
---

# /discover-infrastructure — Infrastructure Discovery

This skill conducts a structured discovery interview covering the customer's current infrastructure landscape. It produces a comprehensive discovery report and appends findings to the engagement's `CONTEXT.md`.

## Prerequisites

- Engagement must be initialized (`/setup` must have been run)
- `engagements/<customer>/CONTEXT.md` must exist

If no engagement exists, prompt the architect to run `/setup` first.

## Workflow

### Step 0: Load Context

Read `engagements/<customer>/CONTEXT.md` to understand:
- Engagement type (shapes which discovery areas to emphasize)
- Environment type (air-gapped changes the questioning)
- Existing RH footprint (skip questions we already know answers to)
- Any prior discovery data (don't re-ask)

If multiple engagements exist, ask which customer this discovery is for.

### Step 1: Discovery Interview

Conduct the interview conversationally. Don't dump all questions at once — work through domains sequentially, adapting based on answers. Skip domains that aren't relevant to the engagement type.

**Interview pacing:** Ask 2-3 questions at a time, grouped by domain. Summarize what you've learned at domain transitions. Flag gaps or inconsistencies as you go.

---

#### Domain 1: Compute & Hosting

| Question | Why It Matters |
|----------|---------------|
| Where do workloads run today? (on-prem DCs, cloud providers, colo, edge) | Deployment topology |
| What hypervisors or bare-metal provisioning? (VMware, Nutanix, bare metal, KVM) | Migration complexity |
| Any existing container platforms? (OpenShift, Kubernetes, ECS, etc.) | Starting point for platform work |
| What's the server fleet look like? (vintage, CPU arch, GPU availability) | Sizing and compatibility |
| How is capacity managed? (manual requests, self-service, quotas) | Operational maturity |

**If engagement type is AI/ML Enablement, also ask:**
- GPU inventory — what models, how many, how allocated?
- Any existing ML serving infrastructure? (MLflow, KServe, SageMaker)
- Data pipeline tooling? (Spark, Airflow, Kafka, Flink)
- Where does training data live and how is it classified?

#### Domain 2: Networking & Connectivity

| Question | Why It Matters |
|----------|---------------|
| Network topology — flat, segmented, zero-trust? | Architecture constraints |
| How are environments connected? (VPNs, direct connect, air-gap boundaries) | Deployment strategy |
| Load balancing approach? (F5, HAProxy, cloud ALB/NLB) | Ingress design |
| DNS and certificate management? | Operational integration |
| Any network-level security appliances? (WAF, IDS/IPS, DLP) | Compliance and traffic flow |

**If environment is air-gapped, also ask:**
- Data transfer mechanisms (CDTS, sneakernet, diode)
- Disconnected registry strategy
- Patch/update delivery cadence

#### Domain 3: Storage & Data

| Question | Why It Matters |
|----------|---------------|
| Primary storage platforms? (SAN, NAS, object, cloud-native) | Persistent volume strategy |
| Backup and DR approach? | Resilience architecture |
| Data classification levels in play? | Security boundaries |
| Database landscape? (RDBMS, NoSQL, data warehouses) | Application dependencies |
| Any shared file systems? (NFS, CIFS, Lustre) | Workload portability |

#### Domain 4: Identity & Access

| Question | Why It Matters |
|----------|---------------|
| Identity provider? (Active Directory, Okta, Keycloak, PIV/CAC) | Auth integration |
| How is RBAC managed today? | Governance model |
| Any federation or SSO across environments? | Multi-cluster identity |
| Service account / machine identity management? | Automation readiness |
| Compliance framework? (FedRAMP, FISMA, CMMC, NIST 800-53) | Compliance scope |

#### Domain 5: CI/CD & Developer Tooling

| Question | Why It Matters |
|----------|---------------|
| Source control? (GitHub, GitLab, Bitbucket, on-prem) | Pipeline integration |
| CI/CD platform? (Jenkins, GitLab CI, GitHub Actions, Tekton) | Build/deploy strategy |
| Artifact management? (Nexus, Artifactory, registry) | Supply chain |
| Developer environments? (local, VDI, cloud IDEs, DevSpaces) | Developer experience |
| Any existing GitOps practices? (ArgoCD, Flux) | Deployment methodology |

#### Domain 6: Operations & Observability

| Question | Why It Matters |
|----------|---------------|
| Monitoring stack? (Prometheus, Datadog, Splunk, Dynatrace) | Observability integration |
| Log aggregation? | Troubleshooting and compliance |
| Alerting and incident management? (PagerDuty, ServiceNow) | Operational integration |
| Configuration management? (Ansible, Puppet, Chef, Salt) | Automation maturity |
| Change management process? (CAB, automated, ad-hoc) | Deployment velocity |

### Step 2: Synthesize and Score

After completing the interview, produce a maturity summary across domains:

```
Infrastructure Maturity Summary
═══════════════════════════════
Compute & Hosting:        [1-5] — <one-line rationale>
Networking & Connectivity: [1-5] — <one-line rationale>
Storage & Data:           [1-5] — <one-line rationale>
Identity & Access:        [1-5] — <one-line rationale>
CI/CD & Developer Tooling: [1-5] — <one-line rationale>
Operations & Observability:[1-5] — <one-line rationale>
```

**Maturity Scale:**
1. **Ad-hoc** — No standardized processes, manual everything
2. **Emerging** — Some tooling in place, inconsistent practices
3. **Defined** — Standardized tools and processes, documented
4. **Managed** — Measured, automated, governance in place
5. **Optimized** — Continuous improvement, self-service, policy-driven

### Step 3: Produce Discovery Report

Write the discovery report to `engagements/<customer>/discovery/infrastructure-discovery.md`:

```markdown
# Infrastructure Discovery Report
**Customer:** <name>
**Date:** <today>
**Conducted by:** <architect>

## Executive Summary
<2-3 paragraph summary of key findings, strengths, and gaps>

## Compute & Hosting
<findings>

## Networking & Connectivity
<findings>

## Storage & Data
<findings>

## Identity & Access
<findings>

## CI/CD & Developer Tooling
<findings>

## Operations & Observability
<findings>

## Maturity Summary
<maturity scores table>

## Key Findings
### Strengths
- <bulleted list>

### Gaps & Risks
- <bulleted list>

### Immediate Recommendations
- <bulleted list of quick wins or blockers to address>

## Questions Requiring Follow-Up
- <anything the architect couldn't answer on-site>
```

### Step 4: Update CONTEXT.md

Append to the engagement's `CONTEXT.md`:

```markdown
## Phase: Infrastructure Discovery
**Date:** <today>
**Conducted by:** <architect>

### Environment Summary
<3-5 sentence summary of the infrastructure landscape>

### Maturity Scores
<scores table>

### Key Findings
- <top 3-5 findings>

### Artifacts Produced
- `discovery/infrastructure-discovery.md` — Full discovery report

### Impact on Engagement Direction
<1-2 sentences on how these findings affect the engagement plan>
```

### Step 5: Recommend Next Steps

Based on findings and engagement type, recommend the logical next skill:

- **High app count discovered** → `/assess-app-portfolio`
- **Security gaps found** → Security & compliance deep-dive
- **AI/ML engagement** → GPU/data landscape warrants specialized assessment
- **Always applicable** → `/assess-app-portfolio` then `/build-deliverable-deck`

## Adaptive Behavior

- **If the architect says "I don't know"** — Note it as a follow-up item. Don't stall.
- **If the architect provides a document** (architecture diagram, CMDB export, network diagram) — Ingest it, extract what you can, and ask only clarifying questions.
- **If the interview is interrupted** — Save progress to a partial discovery file and note where you left off in CONTEXT.md. The skill can be re-run to resume.
- **If running in an air-gapped context** — Emphasize disconnected operations, data transfer mechanisms, and registry mirroring in your questions and findings.

## Checklist Reference

For detailed per-domain checklists, see `knowledge/checklists/infrastructure-discovery-checklist.md`.
