# Product Requirements Document: TeamKnowl

## 1. Executive Summary
**TeamKnowl** is a Kubernetes-native, Git-driven knowledge base platform designed for the modern era of "AI + Human" engineering teams. It provides an Obsidian-like experience (Markdown, [[Wikilinks]], Graph View) but with the lifecycle management, scalability, and security of a Kubernetes-native application.

## 2. Target Audience
- **Software Engineers & Architects**: For design docs and internal knowledge sharing.
- **DevOps/SRE Teams**: For active incident response and SOPs.
- **Support Teams**: For customer-facing and internal documentation.
- **AI Agents**: For "active memory" and context retrieval.

## 3. Core Features (MVP)
### 3.1. Kubernetes-Native Management
- Managed via a `KnowledgeBase` Custom Resource Definition (CRD).
- Deploy multiple independent instances (DevOps, Support, etc.) using the same Operator.
- Installed and managed via **Helm**.

### 3.2. Git-as-Source
- Persistent storage in a Git repository.
- Automated bidirectional sync (Pull from repo, Push from UI).
- Markdown-first content.

### 3.3. Obsidian-like UI
- High-performance web interface.
- Support for `[[Wikilinks]]` and backlinking.
- Mermaid diagram support and basic callouts.

### 3.4. AI-Ready Architecture
- **Headless-First API**: Dedicated endpoints for AI agents to query and update documentation.
- **Structured Metadata**: Support for YAML frontmatter to give AI agents more context.
- **Vector Indexing (Future)**: Integration with a vector DB for semantic search.

## 4. Non-Functional Requirements
- **Simplicity**: No complex directory requirements; content is as flat as possible.
- **Security**: RBAC-controlled access to different knowledge bases.
- **Performance**: Instant loading of notes even with thousands of files.
- **Open Source**: Permissive license (MIT/Apache 2.0).

## 5. Success Metrics
- Average time for a new developer/AI to find a documented solution.
- Ease of deployment via Helm (target < 5 mins).
- High adoption by both human and non-human (agentic) team members.
