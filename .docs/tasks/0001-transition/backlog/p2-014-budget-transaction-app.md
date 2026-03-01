# p2-014 — budgetapp and transactionapp: DTOs, handlers, routes

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the HTTP application layer for regular `budgets` and `transactions`.

---

## Context

Handlers translate HTTP ↔ business types; return `web.Encoder`.

---

## Acceptance Criteria

- [ ] `app/domain/budgetapp/` exists with `model.go`, `filter.go`, `order.go`, `budgetapp.go`, `route.go`.
- [ ] `app/domain/transactionapp/` exists with `model.go`, `filter.go`, `order.go`, `transactionapp.go`, `route.go`.
- [ ] Both domains registered with `/v1/budgets` and `/v1/transactions`.

---

## Steps

1. Create DTOs and converters for `budgetapp` and `transactionapp`.
2. Implement handlers: `create`, `update`, `delete`, `query`, `queryByID`.
3. Configure routes in `route.go`.
