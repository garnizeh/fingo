# p2-018 — investmentdb: SQL store

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement PostgreSQL store for `investmentbus` inside `business/domain/investmentbus/stores/investmentdb/`.

---

## Context

Store implements `investmentbus.Storer`. All SQL is raw via `sqldb` helpers.

---

## Acceptance Criteria

- [ ] `business/domain/investmentbus/stores/investmentdb/investmentdb.go` exists.
- [ ] `business/domain/investmentbus/stores/investmentdb/model.go` exists with `dbInvestment` struct and converters.
- [ ] CRUD methods implementation complete.

---

## Steps

1. Implement `investmentdb.Store`.
2. Map `sql.ErrNoRows` → `investmentbus.ErrNotFound`.
3. Support pagination and filters.
