# Infrastructure Discovery Checklist

Reference checklist for the `/discover-infrastructure` skill. Use this to ensure comprehensive coverage during customer interviews.

## Pre-Interview Preparation

- [ ] Review any existing customer documentation (architecture diagrams, CMDB exports)
- [ ] Check CRM for prior Red Hat engagement history
- [ ] Identify known Red Hat products already in the environment
- [ ] Understand the engagement type to prioritize relevant domains
- [ ] Prepare for environment constraints (air-gap, classification level)

## Compute & Hosting

- [ ] Datacenter locations and strategy (consolidation, hybrid, multi-cloud)
- [ ] Hypervisor platform and version (VMware vSphere version, Nutanix AHV, KVM, Hyper-V)
- [ ] Bare metal provisioning method (PXE, Satellite, Foreman, manual)
- [ ] Server inventory summary (count, vintage, CPU architecture — x86_64, ARM, Power)
- [ ] Existing container platforms (OpenShift, Kubernetes, ECS, AKS, EKS — versions)
- [ ] GPU inventory (model, count, allocation method) — critical for AI/ML engagements
- [ ] Cloud provider accounts and services in use
- [ ] Capacity request and fulfillment process
- [ ] Hardware refresh cycle and upcoming procurements

## Networking & Connectivity

- [ ] Network topology type (flat, segmented, microsegmented, zero-trust)
- [ ] VLAN/subnet strategy for workloads
- [ ] Inter-environment connectivity (VPN, Direct Connect, MPLS, air-gap boundaries)
- [ ] Load balancer platform and configuration style (L4/L7, shared vs dedicated)
- [ ] DNS infrastructure (internal/external, delegation capabilities)
- [ ] Certificate authority and management (internal CA, ACME, manual)
- [ ] Proxy requirements for outbound traffic
- [ ] Network security appliances (WAF, IDS/IPS, DLP, CASB)
- [ ] Air-gap data transfer mechanisms (CDTS, diode, sneakernet) — if applicable
- [ ] IPv4/IPv6 dual-stack requirements
- [ ] MTU considerations (jumbo frames, overlay networks)

## Storage & Data

- [ ] Primary storage platforms (NetApp, Pure, Dell, Ceph, cloud-native)
- [ ] Block, file, and object storage availability
- [ ] Storage performance tiers (SSD, NVMe, spinning disk)
- [ ] Backup solution and RPO/RTO targets
- [ ] Disaster recovery strategy and tested recovery procedures
- [ ] Data classification levels (Unclassified, CUI, Secret, TS/SCI)
- [ ] Database platforms in use (PostgreSQL, Oracle, SQL Server, MongoDB, Redis)
- [ ] Shared filesystem requirements (NFS, CIFS, Lustre, GPFS)
- [ ] Data sovereignty or residency requirements
- [ ] Storage encryption at rest and in transit

## Identity & Access

- [ ] Primary identity provider (Active Directory, Okta, Keycloak, Azure AD)
- [ ] Authentication mechanisms (Kerberos, SAML, OIDC, PIV/CAC)
- [ ] RBAC model and governance process
- [ ] Federation across environments
- [ ] Privileged access management (CyberArk, Thycotic, native)
- [ ] Service account lifecycle management
- [ ] Multi-factor authentication requirements
- [ ] Compliance framework (FedRAMP, FISMA, CMMC, NIST 800-53, PCI-DSS, HIPAA)
- [ ] Audit and access logging requirements

## CI/CD & Developer Tooling

- [ ] Source control platform and branching strategy
- [ ] CI/CD platform and pipeline maturity
- [ ] Artifact repository (Nexus, Artifactory, container registry)
- [ ] Container image build process (Dockerfile, S2I, Buildpacks, Kaniko)
- [ ] Image scanning and vulnerability management in pipeline
- [ ] Deployment methodology (imperative, declarative, GitOps)
- [ ] Environment promotion strategy (dev → stage → prod)
- [ ] Developer environment (local, VDI, DevSpaces, Codespaces)
- [ ] Inner loop developer experience (build/test/debug cycle time)
- [ ] Feature flag or progressive delivery tooling

## Operations & Observability

- [ ] Monitoring platform (Prometheus, Datadog, Splunk, Dynatrace, CloudWatch)
- [ ] Log aggregation and retention (ELK, Splunk, CloudWatch Logs, Loki)
- [ ] Alerting and on-call management (PagerDuty, OpsGenie, ServiceNow)
- [ ] Configuration management (Ansible, Puppet, Chef, Salt)
- [ ] Infrastructure as Code practices and tooling
- [ ] Change management process (CAB, automated approval, ad-hoc)
- [ ] Incident management process and MTTR targets
- [ ] Runbook and documentation practices
- [ ] Capacity monitoring and forecasting
- [ ] SLA/SLO definitions and tracking

## AI/ML Specific (if applicable)

- [ ] ML frameworks in use (PyTorch, TensorFlow, scikit-learn, JAX)
- [ ] Model serving infrastructure (KServe, TF Serving, Triton, SageMaker)
- [ ] Training infrastructure (GPU clusters, cloud training, distributed training)
- [ ] Data pipeline tooling (Spark, Airflow, Kafka, dbt, Flink)
- [ ] Feature store (Feast, Tecton, internal)
- [ ] Experiment tracking (MLflow, W&B, Neptune)
- [ ] Model registry and versioning
- [ ] Data labeling and annotation tooling
- [ ] Training data storage and access patterns
- [ ] Model monitoring and drift detection
