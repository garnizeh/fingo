# p2-004 — creditcarddb: SQL store and DB model

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status: Done**
**Assignee:** unassigned
**Started:** —
**Completed: 2026-02-28**

---

## Goal

Implement the PostgreSQL persistence layer for `creditcardbus` inside `business/domain/creditcardbus/stores/creditcarddb/`.

---

## Context

The store implements `creditcardbus.Storer`. All SQL is raw — no ORM. Use named parameters via `sqldb.NamedExecContext` and `sqldb.NamedQuerySlice`. The DB model (`dbCreditCard`) has `sqlx` struct tags; the business model has none.

---

## Acceptance Criteria

- [x] `business/domain/creditcardbus/stores/creditcarddb/creditcarddb.go` exists and exports `Store` implementing `creditcardbus.Storer`.
- [x] `business/domain/creditcardbus/stores/creditcarddb/model.go` exists with `dbCreditCard` struct (snake_case `db` tags), `toDBCreditCard`, and `toCreditCard` converters.
- [x] `Create` uses an explicit column list INSERT (no `SELECT *`).
- [x] `Query` applies `QueryFilter` conditions and respects `order.By` and `page.Page`.
- [x] `Count` returns the total matching records for a given `QueryFilter`.
- [x] `QueryByID` returns `creditcardbus.ErrNotFound` (wrapped) when no row is found.
- [x] `NewWithTx` returns a new `Store` backed by the provided transaction.
- [x] `go build ./business/domain/creditcardbus/...` exits zero.

---

## Steps

1. Create `model.go`:
   - `dbCreditCard` struct with all columns and `db:"..."` tags.
   - `toDBCreditCard(cc creditcardbus.CreditCard) dbCreditCard`.
   - `toCreditCard(dbc dbCreditCard) creditcardbus.CreditCard`.

2. Create `creditcarddb.go`:
   - `Store` struct with `log *logger.Logger` and `db sqldb.SqlDB` fields.
   - `NewStore(log, db) *Store`.
   - `NewWithTx(tx sqldb.CommitRollbacker) (creditcardbus.Storer, error)`.
   - `Create`: `INSERT INTO credit_cards (...) VALUES (...)` with named params.
   - `Update`: `UPDATE credit_cards SET ... WHERE credit_card_id = :credit_card_id`.
   - `Delete`: `DELETE FROM credit_cards WHERE credit_card_id = :credit_card_id`.
   - `Query`: build dynamic WHERE clause from `QueryFilter`, apply ORDER BY and LIMIT/OFFSET.
   - `Count`: `SELECT COUNT(*) FROM credit_cards WHERE ...`.
   - `QueryByID`: single-row select; map `sql.ErrNoRows` → `creditcardbus.ErrNotFound`.

---

## Notes

Do not add DB store for `invoices` in this task — invoice persistence is a stretch goal; the initial schema just needs the table. Invoice CRUD can be a follow-up task if needed before Phase 3.
