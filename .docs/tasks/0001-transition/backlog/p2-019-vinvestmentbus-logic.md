# p2-019 — vinvestmentbus: mark-to-market virtual view

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the `vinvestmentbus` read-only virtual domain for portfolio valuation.

---

## Context

Use `vproductbus` as the structural template. This bus joins `investments` with current price data (mocked or from a simple adapter). The current value is never persisted.

---

## Acceptance Criteria

- [ ] `business/domain/vinvestmentbus/vinvestmentbus.go` exists.
- [ ] `Storer` interface includes `Query` and `QueryByID` for the enriched view.
- [ ] Portfolio valuation is computed correctly in the store.
- [ ] `go build ./...` exits zero.

---

## Steps

1. Create `model.go` with `VInvestment` struct (adds `CurrentValue`, `GainLoss` fields).
2. Create `vinvestmentdb` store.
3. Use a CTE or join in the SQL to compute valuations.
4. Implement `vinvestmentbus.Business` delegates.
