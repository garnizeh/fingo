# p2-013 — budgetdb and transactiondb: SQL stores

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement PostgreSQL stores for both `budgetbus` and `transactionbus`.

---

## Context

Stores implement the `Storer` interfaces. All SQL is raw via `sqldb` helpers.

---

## Acceptance Criteria

- [ ] `business/domain/budgetbus/stores/budgetdb/budgetdb.go` exists.
- [ ] `business/domain/transactionbus/stores/transactiondb/transactiondb.go` exists.
- [ ] Each store implements the required methods: `Create`, `Update`, `Delete`, `Query`, `Count`, `QueryByID`.
- [ ] `go build ./...` exits zero.

---

## Steps

1. Implement `budgetdb.Store`.
2. Implement `transactiondb.Store`.
3. Test query filters for both.
