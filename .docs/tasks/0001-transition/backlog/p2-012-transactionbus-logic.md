# p2-012 — transactionbus: model, filter, order, events, business

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the `transactionbus` package.

---

## Context

Follow the standard 19-step checklist. Note that a transaction can have a nullable `BudgetID`.

---

## Acceptance Criteria

- [ ] `business/domain/transactionbus/` exists with `model.go`, `filter.go`, `order.go`, `event.go`, `transactionbus.go`.
- [ ] `Transaction` struct exists.
- [ ] `ExtBusiness` interface includes standard CRUD.

---

## Steps

1. Create `model.go`, `filter.go`, `order.go`, `event.go`.
2. Implement `transactionbus.go` with `Storer` and `ExtBusiness`.
3. In `Create`, if `BudgetID` is not nil, notify the delegate of `ActionCreated`.
