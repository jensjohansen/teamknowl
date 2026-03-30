# TeamKnowl — Implementation Plan / Tracker

## 1. Purpose
This document is the core project planner for TeamKnowl implementation. It tracks delivery of milestones in the sequence defined by the PRD and Technical Design.

## 2. Tracker Conventions
### 2.1. Status Values
- Not started
- In progress
- Blocked
- Done

### 2.2. Quality Gates
- All code must pass `docs/standards.md` before pushing.
- Per-file headers are mandatory.
- Unambiguous naming and narrative style are required.

## 3. Milestones

### Milestone 0 — PRD + Architecture + Standards Baseline
**Status**: Done

- [x] Publish `docs/prd.md` and record approval
- [x] Create Technical Design doc (`docs/technical-design.md`)
- [x] Document coding standards and Definition of Done (`docs/standards.md`)
- [x] Initialize project tracker (`docs/implementation-plan-tracker.md`)

---

### Milestone 1 — K8s Operator Scaffolding
**Status**: Not started

- [ ] Initialize Kubebuilder project in `operator/`
- [ ] Define `KnowledgeBase` CRD v1alpha1
- [ ] Implement basic controller reconciliation loop (Status only)
- [ ] Add unit tests for CRD validation

---

### Milestone 2 — Git-Sync Integration
**Status**: Not started

- [ ] Integrate `git-sync` as a sidecar in the KB deployment
- [ ] Implement Secret-based authentication for private Git repos
- [ ] Verify bidirectional sync (Pull from repo, local FS write)

---

### Milestone 3 — Core API & Indexing
**Status**: Not started

- [ ] Develop Go-based API for listing/reading Markdown files
- [ ] Implement Bleve in-memory indexing for full-text search
- [ ] Add [[Wikilink]] resolution logic to the API

---

### Milestone 4 — Obsidian-like Web UI
**Status**: Not started

- [ ] Scaffold React/Next.js frontend
- [ ] Implement Markdown rendering with syntax highlighting
- [ ] Add Graph View visualization for backlinks

---

### Milestone 5 — AI Agentic API
**Status**: Not started

- [ ] Implement `/v1/context` endpoint for flat context retrieval
- [ ] Implement automatic YAML frontmatter injection for AI metadata
- [ ] Create Python-based CLI for AI agents to interact with TeamKnowl

---

### Milestone 6 — Helm Chart & Production Hardening
**Status**: Not started

- [ ] Create comprehensive Helm chart for TeamKnowl Operator
- [ ] Implement NetworkPolicies and RBAC security hardening
- [ ] Finalize documentation and Open Source release baseline
