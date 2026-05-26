# Implementation Plan: TeamKnowl (REDO)

## 1. Phase 1: High-Performance Engine (The Foundation)
**Objective**: Deliver a Rust-native, S3-backed API capable of serving and indexing notes for both humans and AI agents.

### 1.1. Core API Scaffolding (Rust/Axum/Tokio)
- **Task**: Initialize the Rust project with a narrative style and mandatory headers.
- **Quality Gate (G1.1)**: Verify `cargo clippy` and `cargo fmt` pass with zero warnings; verify per-file headers in `main.rs`.

### 1.2. S3-Compatible Storage Client (CEPH)
- **Task**: Implement retrieval logic from S3 (CEPH) instead of local files.
- **Quality Gate (G1.2)**: Provide unit test logs showing successful object retrieval from a mock S3 endpoint.

### 1.3. Sharded Indexing Engine (Tantivy)
- **Task**: Implement a sharded, high-performance search index.
- **Quality Gate (G1.3)**: Verify `cargo test` for index creation, update, and search retrieval against real markdown artifacts.

### 1.4. Wikilink & Context Resolution Logic
- **Task**: Implement the core logic for resolving `[[Wikilinks]]` and generating "flat" context for LLMs.
- **Quality Gate (G1.4)**: Provide `curl` output showing the `/v1/context` endpoint returning markdown with resolved links and YAML frontmatter metadata.

## 2. Phase 2: Obsidian-like Interface (The Human Layer)
**Objective**: Provide a dynamic, high-performance web UI for documentation navigation.

### 2.1. Next.js/TypeScript Scaffolding
- **Task**: Scaffold the UI with strict typing and mandatory headers.
- **Quality Gate (G2.1)**: Verify `npm run lint` and `tsc` pass; verify per-file headers in `page.tsx` and components.

### 2.2. Dynamic Markdown & Mermaid Rendering
- **Task**: Implement real-time rendering of Markdown, Wikilinks, and Mermaid diagrams.
- **Quality Gate (G2.2)**: Provide a screenshot showing a rendered note with active Wikilinks and a functional Mermaid diagram.

### 2.3. Graph View Visualization
- **Task**: Implement the "Obsidian-like" graph view of document relationships.
- **Quality Gate (G2.3)**: Demonstrate the graph view dynamically updating as new notes are indexed.

## 3. Phase 3: AI Orchestration & Context
**Objective**: Enable autonomous agents to interact with the knowledge base.

### 3.1. Agentic Context API
- **Task**: Optimize endpoints for high-volume context retrieval by LLM agents.
- **Quality Gate (G3.1)**: Provide performance benchmarks showing context retrieval latency < 50ms for notes < 10KB.

## 4. Phase 4: Production-Ready & Helm (The Operator Integration)
**Objective**: Reconnect the new Rust/Next.js stack into the K8s Operator and Helm charts.

### 4.1. Helm Chart REDO
- **Task**: Update Helm templates to use the new Rust API and Next.js UI images.
- **Quality Gate (G4.1)**: Verify `helm lint` passes; verify `kubectl get pods` shows all services running and healthy (liveness/readiness).

### 4.2. Security Hardening
- **Task**: Implement NetworkPolicies and RBAC-controlled access.
- **Quality Gate (G4.2)**: Verify that the API is unreachable except via authorized Ingress or Service channels.
