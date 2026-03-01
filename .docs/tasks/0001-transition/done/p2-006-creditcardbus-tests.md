# p2-006 — creditcardbus: testutil and unit tests

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** Done
**Assignee:** unassigned
**Started:** 2026-02-28
**Completed:** 2026-02-28

---

## Goal

Add `testutil.go` and `creditcardbus_test.go` to the `creditcardbus` package, covering core business rules and user ownership isolation.

---

## Context

Business layer tests use a real PostgreSQL instance provided by `business/sdk/dbtest`. The implementation follows the existing project pattern used in `user`, `product`, and `home` domains: tests use `db.BusDomain.CreditCard` directly from `dbtest.New(...)` instead of adding a dedicated `NewUnit` helper. Use `productbus` test files as structural templates.

---

## Acceptance Criteria

- [x] `business/domain/creditcardbus/testutil.go` exists with credit card test data helpers (`TestGenerateNewCreditCards`, `TestGenerateSeedCreditCards`), and tests use `db.BusDomain.CreditCard` per existing domain pattern.
- [x] `business/domain/creditcardbus/creditcardbus_test.go` covers: `TestCreate`, `TestUpdate`, `TestDelete`, `TestQuery`, `TestQueryByID`.
- [x] `TestQuery_UserIsolation` asserts that cards created for user A are never returned when querying for user B.
- [x] `TestCreate` verifies that a zero or negative `Limit` returns `ErrCardLimit`.
- [x] `go test ./business/domain/creditcardbus/...` exits zero.

---

## Steps

1. Use existing `testutil.go` helpers and keep domain test wiring consistent with existing project pattern (`db.BusDomain.CreditCard`).

2. Create `creditcardbus_test.go`:
   - Use `dbtest.NewUnit` to get a database and logger.
   - Seed two users (user A and user B) using `userbus` testutil.
   - `TestCreate`: create a card for user A, assert fields match.
   - `TestUpdate`: update the card name; assert `DateUpdated` changed.
   - `TestDelete`: delete the card; assert `QueryByID` returns `ErrNotFound`.
   - `TestQuery`: create 3 cards for user A; assert `Query` with `UserID` filter returns all 3.
   - `TestQuery_UserIsolation`: create 1 card for user B; query with user A filter; assert user B card is absent.
   - `TestQueryByID`: assert `QueryByID` with a non-existent UUID returns `ErrNotFound`.

---

## Notes

Do not use mocks. The test must exercise real SQL. Tests requiring a database must be skipped if the `TEST_DB_URL` environment variable is not set.
