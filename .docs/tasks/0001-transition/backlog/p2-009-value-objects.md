# p2-009 — Value Objects: currency, ticker, assetclass, direction

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Create the shared value objects required for the Investment and Transaction domains.

---

## Context

Value objects live in `business/types/`. They provide type safety and validation beyond raw primitives. Use `business/types/money` or `name` as structural templates. Each type must have a `Parse(s string) (Type, error)` constructor and a `String() string` method.

---

## Acceptance Criteria

- [ ] `business/types/currency/currency.go` exists with `BRL`, `USD`, `EUR`.
- [ ] `business/types/ticker/ticker.go` exists with uppercase validation (1-10 chars).
- [ ] `business/types/assetclass/assetclass.go` exists with `Stock`, `Crypto`, `FixedIncome`.
- [ ] `business/types/direction/direction.go` exists with `Debit`, `Credit`.
- [ ] Each type implements `Parse` and `String`.
- [ ] `go build ./business/types/...` exits zero.

---

## Steps

1. Create `business/types/currency/`:
   - `type Currency string`
   - Constants: `BRL`, `USD`, `EUR`
   ```go
   func Parse(s string) (Currency, error)
   func (c Currency) String() string
   ```
2. Create `business/types/ticker/`:
   - `type Ticker string`
   - `Parse` validates length and uppercase.
3. Create `business/types/assetclass/`:
   - `type AssetClass string`
   - Constants: `Stock`, `Crypto`, `FixedIncome`.
4. Create `business/types/direction/`:
   - `type Direction string`
   - Constants: `Debit`, `Credit`.
