# p2-015 — budget-transaction wiring and tests

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Wire the Budget and Transaction domains into the service entry point and register routes.

---

## Context

Follow the five-step DI pattern for both domains.

---

## Acceptance Criteria

- [ ] `app/sdk/mux/mux.go` supports `BudgetBus` and `TransactionBus`.
- [ ] `api/services/fingo/main.go` constructs both stores, extensions, and buses.
- [ ] `go build ./...` exits zero.

---

## Steps

1. Update `BusConfig` in `app/sdk/mux/mux.go`.
2. Add DI logic in `api/services/fingo/main.go`.
3. Add `budgetapp.Routes` and `transactionapp.Routes` to the mux.
4. Confirm everything compiles.
