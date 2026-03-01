# p2-011 — budgetbus: model, filter, order, events, business

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the `budgetbus` package.

---

## Context

Follow the standard 19-step checklist. This domain also registers a delegate function to react to transaction events.

---

## Acceptance Criteria

- [ ] `business/domain/budgetbus/` exists with `model.go`, `filter.go`, `order.go`, `event.go`, `budgetbus.go`.
- [ ] `Budget` struct exists.
- [ ] `ExtBusiness` interface includes standard CRUD.
- [ ] `registerDelegateFunctions` registers for `transactionbus.ActionCreated`.

---

## Steps

1. Create `model.go`, `filter.go`, `order.go`, `event.go`.
2. Implement `budgetbus.go` with `Storer` and `ExtBusiness`.
3. Implement `registerDelegateFunctions` to handle updates when a transaction is added to a budget.
