# Red Hat Solution Patterns — Quick Reference

This document provides a curated reference of Red Hat solution patterns for use during engagement architecture recommendations. Each pattern includes selection criteria, key components, and common pairings.

> **Note:** This is a living reference. Architects should validate against current Red Hat documentation and reference architectures. Product capabilities evolve rapidly.

## Platform Patterns

### Multi-Node OpenShift (HA Production)
**When:** Production workloads requiring HA, multi-team environments, large-scale deployments
**Components:** OCP (3+ control plane, n workers), RHCOS, integrated monitoring/logging
**Sizing Baseline:** 3 control (8 vCPU, 32 GB), 3+ worker (16 vCPU, 64 GB) — scale based on workload
**Pairs with:** ACM (multi-cluster), ACS (security), ODF (storage), Service Mesh, GitOps

### Single Node OpenShift (SNO)
**When:** Edge, remote sites, lab/dev, resource-constrained, single-purpose deployments
**Components:** OCP on a single node (combined control+worker), RHCOS
**Sizing Baseline:** 8 vCPU, 32 GB RAM minimum (16 vCPU / 64 GB recommended)
**Constraints:** No HA for control plane, limited to ~50-100 pods typically, some operators may not support SNO
**Pairs with:** ACM (fleet management), GitOps (declarative management)

### ROSA / ARO (Managed Cloud)
**When:** Cloud-first strategy, managed operations preference, AWS/Azure native
**Components:** ROSA (AWS) or ARO (Azure), managed control plane
**Considerations:** Shared responsibility model, cloud-specific networking, STS/workload identity
**Pairs with:** Cloud-native storage, cloud IAM integration, ACM for hybrid

### OCP Virtualization (OCP-V)
**When:** VMware migration, VM + container coexistence, reduce virtualization licensing costs
**Components:** OCP with OpenShift Virtualization operator, KubeVirt
**Considerations:** VM migration toolkit (MTV), Windows VM support, live migration requirements
**Pairs with:** ODF (storage for VMs), bare metal (best performance), migration toolkit

### Advanced Cluster Management (ACM)
**When:** Multi-cluster governance, fleet management, policy enforcement across environments
**Components:** ACM hub cluster, managed clusters, policy engine, observability
**Considerations:** Hub sizing, network connectivity to managed clusters, policy complexity
**Pairs with:** Any OCP deployment pattern, GitOps, ACS (centralized security)

## Application Modernization Patterns

### Containerize with S2I / Buildpacks
**When:** Standard web frameworks (Java/Spring, Node.js, Python, .NET), rapid containerization
**Approach:** Source-to-Image or Cloud Native Buildpacks transform source code into container images
**Effort:** Low — minimal code changes if app is 12-factor adjacent
**Pairs with:** CI/CD pipelines, OpenShift Developer Console

### Microservices with Service Mesh
**When:** Complex distributed applications, zero-trust networking requirements, advanced traffic management
**Components:** OpenShift Service Mesh (Istio-based), Kiali, Jaeger/Tempo
**Considerations:** Sidecar overhead, mTLS complexity, team adoption curve
**Pairs with:** GitOps, ACS, distributed tracing

### Event-Driven Architecture
**When:** Asynchronous processing, data pipelines, event sourcing, decoupled microservices
**Components:** AMQ Streams (Kafka), Serverless (Knative), Camel K
**Considerations:** Event schema management, exactly-once semantics, consumer group strategy
**Pairs with:** Service Mesh, monitoring (consumer lag alerts)

### GitOps with ArgoCD
**When:** Declarative deployment, multi-environment promotion, audit trail, drift detection
**Components:** OpenShift GitOps (ArgoCD), Git repository structure, ApplicationSets
**Considerations:** Secret management (Sealed Secrets or External Secrets), RBAC model, sync policies
**Pairs with:** Every pattern — GitOps is a deployment methodology, not a workload pattern

## AI/ML Patterns

### OpenShift AI (RHOAI)
**When:** ML model training and serving on OpenShift, data science team enablement
**Components:** RHOAI operator, JupyterHub, model serving (KServe), data science pipelines
**Considerations:** GPU operator required for GPU workloads, storage for datasets, namespace isolation
**Pairs with:** GPU operator, ODF (data storage), Node Feature Discovery

### Self-Hosted LLM Inference (vLLM)
**When:** On-prem LLM serving, data sovereignty requirements, air-gapped AI
**Components:** vLLM on OCP, GPU operator, model storage, inference endpoint
**Considerations:** GPU memory requirements (model size dependent), quantization options, batching strategy
**Pairs with:** RHOAI (serving framework), persistent storage for models, Service Mesh (inference routing)

### MLOps Pipeline
**When:** Model lifecycle management, reproducible training, automated retraining
**Components:** MLflow (experiment tracking), KServe (serving), data science pipelines (Kubeflow/Elyra)
**Considerations:** Model versioning, feature store integration, monitoring for drift
**Pairs with:** RHOAI, GitOps (model deployment), S3-compatible storage

## Security Patterns

### Advanced Cluster Security (ACS / StackRox)
**When:** Container security posture management, vulnerability scanning, runtime protection
**Components:** ACS central, secured clusters, scanner, admission controller
**Considerations:** Integration with CI/CD for shift-left, policy tuning to reduce noise
**Pairs with:** Every OCP deployment, CI/CD pipelines, SIEM integration

### Compliance Operator
**When:** STIG compliance, CIS benchmarks, NIST 800-53 controls, automated compliance scanning
**Components:** Compliance Operator, OpenSCAP, tailored profiles
**Considerations:** Profile customization, remediation automation, scan scheduling
**Pairs with:** ACS, RHEL (host-level compliance), Ansible (remediation)

### Supply Chain Security
**When:** Image provenance, SBOM requirements, software supply chain integrity
**Components:** Cosign/Sigstore (signing), Syft (SBOM), SLSA provenance, Tekton Chains
**Considerations:** Key management for signing, policy enforcement at admission, SBOM storage
**Pairs with:** ACS (policy enforcement), Tekton (CI pipeline), OPA/Gatekeeper

## Automation Patterns

### Ansible Automation Platform
**When:** Cross-domain automation, infrastructure provisioning, compliance remediation, Day 2 ops
**Components:** AAP controller, execution environments, automation hub, Event-Driven Ansible
**Considerations:** Execution environment strategy, credential management, RBAC model
**Pairs with:** Every pattern — automation is horizontal

### Infrastructure as Code (Terraform + Ansible)
**When:** Reproducible infrastructure provisioning, multi-cloud, GitOps for infrastructure
**Components:** Terraform (provisioning), Ansible (configuration), Git (state/playbook management)
**Considerations:** State management, drift detection, blast radius controls
**Pairs with:** OCP installation (IPI/UPI), cloud provisioning, GitOps workflows
