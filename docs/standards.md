# Coding Standards & Definition of Done: TeamKnowl

## 1. Coding Best Practices
### 1.1. Unambiguous Naming & Narrative Style
- **No abbreviations** in variable names. The code should read like a narrative; non-programmers should be able to follow the general logic.
- Use names that make sense to all stakeholders (CEO, CTO, architects), not just software engineers (e.g., `gitRepositoryUrl` instead of `gitRepoUrl`).
- Function names should be descriptive and clearly indicate their purpose (e.g., `FetchLatestDocumentationFromGit` instead of `pullData`).

### 1.2. Meaningful Commenting & Documentation
- Code should explain **what** it does by using clear and expressive naming.
- Every method that can be imported requires doc comments (e.g., docstrings) that explain:
  - **what** it is,
  - **why** an importing module would use it,
  - and links to supporting documentation.
- Avoid repeating the code's logic in comments (e.g., `// Increment counter` is redundant).

### 1.3. Per-File Headers (Required)
Every file that supports comments must include a header at the top containing:
- File name and repository-relative path.
- Short description of what it provides.
- Why it matters to the product/business.
- Copyright and license (MIT).

**Go Template:**
```go
/*
File: teamknowl/pkg/example.go
Purpose: <what this file provides>
Product/business importance: <why this matters to TeamKnowl>

Copyright (c) 2026 John K Johansen
License: MIT (see LICENSE)
*/
```

**YAML/Helm Template:**
```yaml
# File: charts/teamknowl/values.yaml
# Purpose: <what this file provides>
# Product/business importance: <why this matters to TeamKnowl>
#
# Copyright (c) 2026 John K Johansen
# License: MIT (see LICENSE)
```

### 1.4. Language-Specific Guidelines
- **Go**: Primary language for the **Operator, API Gateway, and AI Orchestration**. Focus on high-concurrency via goroutines and idiomatic error handling. Use `golangci-lint`.
- **Rust**: Primary language for the **Indexing Engine and performance-critical data processing**. Follow `clippy` and `rustfmt` standards.
- **TypeScript**: Used for the **Web UI and IDE Extensions**. Enforce `eslint` and `prettier`.

## 2. Security Best Practices
- **No Hardcoded Secrets**: Never commit secrets or credentials. Use Kubernetes Secrets.
- **Least Privilege**: Grant only necessary permissions to the Operator and services.
- **Dependency Hygiene**: Pin versions and review dependency changes for non-permissive licenses (MIT/Apache preferred).

## 3. Definition of Done (DoD)
A task is considered **Done** only when:
1. **Product**: Relevant PRD milestone acceptance criteria are satisfied.
2. **Design**: Consistent with Technical Design doc or deviation is documented.
3. **Planning**: Implementation tracker updated.
4. **Quality**: Tests updated and passing; 80% coverage target for core logic.
5. **Security**: No secrets introduced; passes security linting.
6. **Documentation**: Per-file headers present; API docs updated; narrative style followed.
7. **OSS Posture**: Packaging stance preserved (Helm charts only, no bundled proprietary binaries).
