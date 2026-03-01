# p2-017 — investmentbus: model, filter, order, events, business, testutil

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the standard `investmentbus` package.

---

## Context

Follow the 19-step checklist. Note that investments use several value objects: `ticker`, `assetclass`, `currency`.

---

## Acceptance Criteria

- [ ] `business/domain/investmentbus/` exists with all core files.
- [ ] `Investment` entity uses correct types from `business/types/`.
- [ ] `ExtBusiness` includes standard CRUD.

---

## Steps

1. Create `model.go`, `filter.go`, `order.go`, `event.go`.
2. Implement `investmentbus.go` with `Storer` and `ExtBusiness`.
3. Add `testutil.go`.
4. Test Query and Count.
