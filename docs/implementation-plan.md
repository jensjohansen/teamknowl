# Implementation Plan: TeamKnowl

## 1. Phase 1: Foundation (MVP)
### 1.1. Operator Scaffolding
- Initialize project with Kubebuilder (Go).
- Define `KnowledgeBase` CRD schema.
- Create basic controller logic for creating and deleting a deployment.

### 1.2. Git-Sync Engine
- Research and select a Git-Sync tool (likely `kubernetes-sigs/git-sync`).
- Set up automated testing for bidirectional sync between a local Git repo and the API service.

### 1.3. Core API
- Develop basic REST API (Go/Gin) for listing and reading notes.
- Implement file-based storage manager with locking (preventing write-conflicts between syncs).

### 1.4. Helm Chart
- Create initial Helm chart for installing the TeamKnowl Operator.

## 2. Phase 2: User Experience (The "Obsidian" Layer)
### 2.1. Basic UI
- Develop a React/Next.js frontend to list and read Markdown files.
- Implement syntax highlighting (Prism.js) and [[Wikilink]] navigation.

### 2.2. Search Indexing
- Integrate a lightweight indexing engine (like Bleve) to support full-text search.

## 3. Phase 3: AI & Extensions
### 3.1. Agentic API
- Create specific endpoints for AI agents (`/v1/context`, `/v1/suggest`).
- Implement automatic frontmatter injection.

### 3.2. DevOps & Support Sidecars
- Research and draft requirements for DevOps-specific sidecar (K8s event watcher).
- Research and draft requirements for Support-specific connector (HubSpot/Jira).

## 4. Phase 4: Production Ready
### 4.1. Security Hardening
- Implement NetworkPolicies and RBAC-controlled access.
- Finalize Secrets Management and certificate-based auth for Git.

### 4.2. Performance Optimization
- Optimize the indexing process for large repositories.
- Implement caching for the UI.
