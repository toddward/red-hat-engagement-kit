---
type: fact
date: 2026-06-01
verbatim: false
phase: architecture
sensitivity: none
tags: [airgap, registry, openshift, disconnected-install]
supersedes: null
---
Disconnected install path: mirror images with oc-mirror, carry them across the
one-way data diode, serve from an internal Quay registry. Reusable for any
air-gapped OpenShift build — a strong candidate to promote (after scrub) into a
cross-engagement `pattern`.
