# NNNN — <Short Title>

<!--
INSTRUCTIONS (delete this block before committing):
- Replace NNNN with the next four-digit sequence number (e.g. 0002).
- Replace <Short Title> with a concise description of the feature or change.
- Fill in Status, Author, and Last Updated.
- Remove sections that are not relevant, but keep the fixed structure.
- Each phase or section must produce at least one task file under
  .docs/tasks/NNNN-<slug>/backlog/.
- Move tasks to .docs/tasks/NNNN-<slug>/done/ when complete.
-->

**Status:** Draft | In Progress | Completed  
**Author:** <name>  
**Last Updated:** YYYY-MM-DD  
**Reference:** `.docs/copilot-instructions.md`

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State](#2-current-state)
3. [Target State](#3-target-state)
4. [Implementation Phases](#4-implementation-phases)
5. [Dependency Map](#5-dependency-map)
6. [Risk Register](#6-risk-register)
7. [Definition of Done](#7-definition-of-done)
8. [Backlog](#8-backlog)

---

## 1. Executive Summary

<!--
One or two paragraphs answering:
- What problem does this solve?
- What is the high-level approach?
- What are the non-negotiable constraints?
-->

---

## 2. Current State

<!--
Describe what exists today.
Use a table where possible.

| Layer | Package | Status |
|---|---|---|
| ... | ... | ✅ Keep / ⚠️ Extend / 🔁 Replace / 🆕 New |

Include any relevant patterns or conventions the current code already establishes
that this work must respect.
-->

---

## 3. Target State

<!--
Describe the desired end state: directory layout, domain map, service map.
Use code blocks for tree views.

```
business/domain/
├── existing/      ✅
└── new-domain/    🆕
```
-->

---

## 4. Implementation Phases

<!--
Break the work into sequential, independently-validatable phases.
Each phase must have:
  - A clearly stated Goal
  - A "Why this order matters" rationale
  - Sub-sections per deliverable (migration, code change, wiring, test)
  - A Validation section at the end

Phase numbers must match the task folder entries in .docs/tasks/NNNN-<slug>/backlog/.
-->

### Phase 1 — <Phase Name>

**Goal:** <!-- one sentence -->

**Why this order matters:** <!-- one sentence -->

---

#### 1.1 — <Deliverable>

<!-- Description, code snippets, SQL blocks, shell commands as needed. -->

---

#### 1.N — Validation

```bash
go build ./...
go test ./...
```

Expected outcome: <!-- what must be true before proceeding -->

---

### Phase 2 — <Phase Name>

<!-- repeat structure above -->

---

## 5. Dependency Map

<!--
Show the construction order for new buses/packages in main.go.
Use an ASCII tree or a numbered list.

Example:
  delegate        (no dependencies)
      │
  auditBus        (delegate, auditdb)
      │
  newBus          (delegate, userBus, newdb, otel ext, audit ext)
-->

---

## 6. Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| <!-- describe risk --> | Low / Medium / High | Low / Medium / High | <!-- mitigation --> |

---

## 7. Definition of Done

A phase is complete only when **all** of the following are true:

- [ ] `go build ./...` exits zero.
- [ ] `go vet ./...` exits zero.
- [ ] `go test ./...` exits zero.
- [ ] The migration runs against a fresh PostgreSQL instance without errors.
- [ ] The migration is **idempotent** (running it twice produces no error).
- [ ] Every new domain handler is covered by at least one integration test in `app/sdk/apitest`.
- [ ] <!-- add phase-specific acceptance criteria here -->

---

## 8. Backlog

<!--
List every task to be created under .docs/tasks/NNNN-<slug>/backlog/.
Each line here corresponds to one task file.

Naming convention: pN-NNN-<slug>.md
  - N   = phase number (matches the Phase N section in this doc)
  - NNN = three-digit sequence that RESETS to 001 at the start of each phase
  - slug = short kebab-case description

This prefix makes it immediately clear which phase each task belongs to,
which is critical for multi-phase design docs like 0001-transition.

Format: pN-NNN-<slug> — <one-line description>

Example (two phases, multiple tasks each):
p1-001-create-migration     — Append SQL migration versions for new tables
p1-002-implement-model      — Write model.go with business types
p1-003-implement-store      — Write store SQL implementation
p2-001-implement-bus        — Write bus interface, Business struct, NewBusiness
p2-002-implement-otel-ext   — Write OTel extension decorator
p2-003-implement-audit-ext  — Write audit extension decorator
p2-004-implement-app-layer  — Write app DTOs and handler methods
p2-005-wire-mux             — Add bus to BusConfig and main.go DI
p2-006-write-tests          — Unit + integration tests
-->
