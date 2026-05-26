# Strategic Research: Persistence Trade-off Analysis
**Author**: Zencoder AI
**Date**: 2026-03-30

This document evaluates the strategic and technical trade-offs between using **Git-as-Source** versus a **Document Store (JSON-based, e.g., RethinkDB)** for TeamKnowl's persistence layer.

## 1. Stakeholder Simplicity Matrix

"Simplicity" is subjective and varies across organizational levels. We evaluate each option based on the value it delivers to these stakeholders.

| **Stakeholder** | **Git-as-Source** | **Document Store (JSON)** |
| :--- | :--- | :--- |
| **C-Suite (Strategic)** | **High Simplicity**: No vendor lock-in; Data is portable and durable. | **High Simplicity**: Unstructured JSON is an open standard. Portable, schema-flexible, and natively supports rich metadata. |
| **Middle-Management (Tactical)** | **Moderate Simplicity**: Integrates with Git workflows; Audit trails are native. | **High Simplicity**: Native support for real-time collaboration, metadata indexing, and complex relationship mapping. |
| **Line Staff (Functional)** | **Moderate Simplicity**: CLI/Local-friendly, but manual "Sync" friction. | **High Simplicity**: Fast, real-time editing. No "Commit" ceremony required. |

## 2. Technical Research: Versioning Mechanics (Google Docs Style)

### 2.1. How Google Docs Handles History
Google Docs (and similar collaborative editors) use **Operational Transformation (OT)** or **Conflict-free Replicated Data Types (CRDTs)**:
- **Mechanic**: They don't just store "diffs" (like Git). They store a log of operations (e.g., "Insert 'X' at position 5") and periodically take **Full Snapshots**.
- **Business Value**: This allows for near-instant "Time Travel" to any point in the document's history by loading a snapshot and playing back a small log of operations.
- **Efficiency**: It is computationally cheaper to store and retrieve whole versions or large snapshots in cheap S3 storage than it is to reconstruct a document from thousands of low-level Git-style line-diffs.

### 2.2. Snapshots vs. Line-Diffing
- **Git (Diff-based)**: Extremely efficient for disk space (only stores what changed), but expensive for compute (must "apply" diffs to recreate an old version).
- **Document Store (Snapshot-based)**: More disk space (stores multiple whole versions), but extremely cheap for compute (O(1) retrieval of any version).
- **Strategy for TeamKnowl**: Given that S3 storage is cheap and compute/latency is expensive, a **Snapshot-first approach** in a JSON Document Store delivers better business responsiveness.

## 3. Comparative Analysis (The 5 Whys)

### Option A: Git-as-Source
1.  **Why use Git?** To provide a robust, auditable version history for every change.
2.  **Why auditable history?** To ensure accountability and "revertability" in a shared team environment.
3.  **Why not a DB-based history?** To leverage an existing industry-standard versioning tool.
4.  **Why Git specifically?** It allows documentation to mirror the code lifecycle.
5.  **Strategic Outcome**: High durability, but introduces significant implementation and operational complexity (Merge conflicts, "Sync" lag).

### Option B: JSON Document Store (Strategic Choice)
1.  **Why use a Document Store?** To store knowledge artifacts alongside their operational metadata (version history, view logs, AI insights) in a single, flexible source.
2.  **Why JSON?** To maintain an open, non-proprietary data format while allowing for partially unstructured data.
3.  **Why not Git-style diffs?** Because low-level add/delete diffing is computationally expensive compared to simple version snapshots in cheap S3 storage.
4.  **Why snapshots?** To provide near-instant "Time Travel" and high-concurrency editing for knowledge workers.
5.  **Strategic Outcome**: Prioritizes **Flexibility and High-Speed Collaboration** without the "Git friction" for non-engineering stakeholders.

## 4. Final Recommendation: The JSON-Snapshot Architecture

We should move toward a **JSON Document Store** model:
- **Primary Storage**: Artifacts stored as JSON documents (e.g., in RethinkDB or directly in S3 with a metadata index).
- **Versioning Strategy**: **Full-Version Snapshots**. Every "save" (or milestone) creates a new version of the JSON object. 
- **Business Value**: This delivers the "Google Docs simplicity" users expect, maintains data flexibility for AI, and keeps compute costs low by avoiding complex diff-reconstruction.
