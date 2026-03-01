# p2-003 — creditcardbus: Storer, ExtBusiness, and Business struct

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status: Done**
**Assignee:** unassigned
**Started:** —
**Completed: 2026-02-28**

---

## Goal

Implement `creditcardbus.go` — the file that declares the `Storer` interface, the `ExtBusiness` interface, the `Extension` function type, the `Business` struct, and `NewBusiness`.

---

## Context

`Storer` is the persistence contract known only to this package and its tests. `ExtBusiness` is the public contract consumed by the app layer and extension decorators. `Business` implements `ExtBusiness` by delegating persistence to `Storer`.

---

## Acceptance Criteria

- [x] `business/domain/creditcardbus/creditcardbus.go` exists.
- [x] `Storer` interface includes `NewWithTx`, `Create`, `Update`, `Delete`, `Query`, `Count`, `QueryByID`.
- [x] `ExtBusiness` interface includes `NewWithTx`, `Create`, `Update`, `Delete`, `Query`, `Count`, `QueryByID`.
- [x] `Extension` type is `func(ExtBusiness) ExtBusiness`.
- [x] `Business` struct has fields: `log *logger.Logger`, `userBus userbus.ExtBusiness`, `delegate *delegate.Delegate`, `storer Storer`.
- [x] `NewBusiness(log, userBus, delegate, storer, extensions...)` applies extensions in correct order and returns `ExtBusiness`.
- [x] Sentinel errors `ErrNotFound` and `ErrCardLimit` are declared.
- [x] `Create` validates that `Limit > 0`, returning `ErrCardLimit` if not.
- [x] `go build ./business/domain/creditcardbus/...` exits zero.

---

## Steps

1. Declare sentinel errors:
   ```go
   var (
       ErrNotFound  = errors.New("credit card not found")
       ErrCardLimit = errors.New("credit card limit must be positive")
   )
   ```
2. Declare `Storer` interface.
3. Declare `ExtBusiness` interface.
4. Declare `Extension` type.
5. Implement `Business` struct.
6. Implement `NewBusiness` — apply extensions in reverse slice order, wrapping the concrete `Business` with each decorator.
7. Implement `Create` (assign `uuid.New()` for `ID`, set `DateCreated`/`DateUpdated`, validate limit).
8. Implement `Update` (apply non-nil fields from `UpdateCreditCard`, set `DateUpdated`).
9. Implement `Delete`, `Query`, `Count`, `QueryByID`.
10. Implement `NewWithTx` — builds a new `Business` with the storer obtained from `storer.NewWithTx(tx)`.

---

## Notes

`actorID` is passed through to extensions (OTel, audit) but is not persisted directly in the `credit_cards` table (ownership is via `UserID` on the entity).
