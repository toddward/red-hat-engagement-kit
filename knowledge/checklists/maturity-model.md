# Engagement Maturity Model

Standard maturity scale used across all assessment skills for consistent scoring.

## 5-Level Maturity Scale

| Level | Label | Characteristics |
|-------|-------|----------------|
| **1** | **Ad-hoc** | No standardized processes. Manual, reactive, tribal knowledge. Individual heroics. No documentation. |
| **2** | **Emerging** | Some tooling in place. Inconsistent practices across teams. Partial documentation. Aware of gaps. |
| **3** | **Defined** | Standardized tools and processes. Documented procedures. Consistent across most teams. Repeatable. |
| **4** | **Managed** | Measured and automated. Governance in place. KPIs tracked. Proactive monitoring. Self-service where appropriate. |
| **5** | **Optimized** | Continuous improvement culture. Policy-driven automation. Data-informed decisions. Self-healing systems. Industry-leading practices. |

## Scoring Guidelines

- Score based on **current state**, not aspirations or plans
- Consider the **worst-performing aspect** within a domain — a chain is only as strong as its weakest link
- Use half-point scores (e.g., 2.5) sparingly — only when the domain genuinely straddles two levels
- Document your rationale for each score — it's the "why" that matters for recommendations
- Scores below 3 generally indicate a gap requiring remediation before production deployment
- Scores of 4-5 indicate strengths to leverage and build upon

## Domain-Specific Scoring Examples

### Compute & Hosting
| Score | Example |
|-------|---------|
| 1 | Manual server provisioning, no capacity tracking, hardware age unknown |
| 3 | Automated VM provisioning (Terraform/Ansible), capacity dashboards, refresh cycle defined |
| 5 | Self-service infrastructure, auto-scaling, predictive capacity management, IaC for everything |

### CI/CD & Developer Tooling
| Score | Example |
|-------|---------|
| 1 | Manual builds, FTP deployments, no source control governance |
| 3 | CI/CD pipelines for most apps, artifact management, environment promotion defined |
| 5 | GitOps everywhere, progressive delivery, automated security scanning, <15 min build-to-deploy |

### Security & Compliance
| Score | Example |
|-------|---------|
| 1 | No scanning, manual patching, compliance is a spreadsheet |
| 3 | Automated vulnerability scanning, RBAC defined, compliance framework mapped |
| 5 | Shift-left security, automated compliance, supply chain signing, zero-trust enforced |

## Aggregation

When computing an overall score from multiple domains:
- Use **weighted average** if engagement type makes some domains more critical
- Default to **simple average** if all domains are equally important
- Always report **both** the overall score and individual domain scores
- Highlight any domain scoring ≤ 2 as a **critical gap** regardless of overall score
