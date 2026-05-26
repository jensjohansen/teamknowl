# TeamKnowl — Implementation Plan / Tracker (REDO)

## 1. Tracker Conventions
### 1.1. Status Values
- **Not started**: Task initialized but not touched.
- **In progress**: Active work on task.
- **Blocked**: Waiting on external resolution.
- **Done**: Proof attached and verified.

### 1.2. Quality Gate Protocol
No task may be marked "Done" without attaching the following **Proof** to the session:
1. **Standards Compliance**: Verification that `docs/standards.md` (Headers, Naming) is followed.
2. **Static Analysis**: Successful `cargo clippy` (Rust) or `npm run lint` (UI) logs.
3. **Logic Verification**: `cargo test` or `curl` responses showing real (non-mocked) output.
4. **Contract Adherence**: Explicit confirmation that the implementation matches the PRD/Technical Design.

---

## 2. Milestones

### Milestone 1 — High-Performance Engine (REDO)
**Status**: In progress

- [ ] **1.1. Core API Scaffolding (Rust/Axum/Tokio)**
  - [ ] **Proof**: `cargo clippy` and `ls -l` showing per-file headers.
  - [ ] **Status**: Not started
- [ ] **1.2. S3-Compatible Storage Client (CEPH)**
  - [ ] **Proof**: `cargo test` log showing successful mock-S3 retrieval.
  - [ ] **Status**: Not started
- [ ] **1.3. Sharded Indexing Engine (Tantivy)**
  - [ ] **Proof**: Search benchmark output and logic verification logs.
  - [ ] **Status**: Not started
- [ ] **1.4. Wikilink & Context Resolution Logic**
  - [ ] **Proof**: `curl /v1/context` output with resolved links and metadata.
  - [ ] **Status**: Not started

---

### Milestone 2 — Obsidian-like Interface (REDO)
**Status**: Not started

- [ ] **2.1. Next.js/TypeScript Scaffolding**
  - [ ] **Proof**: `npm run lint` and header verification.
- [ ] **2.2. Dynamic Markdown & Mermaid Rendering**
  - [ ] **Proof**: Screenshot or browser logs demonstrating rendering logic.
- [ ] **2.3. Graph View Visualization**
  - [ ] **Proof**: Demonstration of note node/edge generation.

---

### Milestone 3 — AI Orchestration & Context
**Status**: Not started

- [ ] **3.1. Agentic Context API**
  - [ ] **Proof**: Latency benchmark results.

---

### Milestone 4 — Production-Ready & Helm
**Status**: Not started

- [ ] **4.1. Helm Chart REDO**
  - [ ] **Proof**: `helm lint` and `kubectl get pods` logs.
- [ ] **4.2. Security Hardening**
  - [ ] **Proof**: NetworkPolicy test logs.
