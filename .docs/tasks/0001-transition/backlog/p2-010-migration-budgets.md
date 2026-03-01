# p2-010 — Migration SQL: budgets and transactions tables

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Append versioned migration blocks for the `budgets` and `transactions` tables to the existing migration file.

---

## Context

Darwin tracks each block by checksum. Existing blocks must never be edited. New blocks are appended at the end of `business/sdk/migrate/sql/migrate.sql`.

---

## Acceptance Criteria

- [ ] Version `2.03` block creates table `budgets`.
- [ ] Version `2.04` block creates table `transactions`.
- [ ] Running the migration twice produces no error.
- [ ] `go build ./...` exits zero.

---

## Steps

1. Open `business/sdk/migrate/sql/migrate.sql`.
2. Append version `2.03`:
```sql
-- Version: 2.03
-- Description: Create table budgets
CREATE TABLE IF NOT EXISTS budgets (
    budget_id      UUID           NOT NULL,
    user_id        UUID           NOT NULL REFERENCES users(user_id),
    name           TEXT           NOT NULL,
    total_amount   NUMERIC(12, 2) NOT NULL,
    period         TEXT           NOT NULL,
    start_date     DATE           NOT NULL,
    date_created   TIMESTAMPTZ    NOT NULL,
    date_updated   TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (budget_id)
);
```
3. Append version `2.04`:
```sql
-- Version: 2.04
-- Description: Create table transactions
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID           NOT NULL,
    user_id        UUID           NOT NULL REFERENCES users(user_id),
    budget_id      UUID           NULL REFERENCES budgets(budget_id),
    direction      TEXT           NOT NULL,
    amount         NUMERIC(12, 2) NOT NULL,
    description    TEXT           NOT NULL,
    category       TEXT           NOT NULL,
    transaction_at TIMESTAMPTZ    NOT NULL,
    date_created   TIMESTAMPTZ    NOT NULL,
    date_updated   TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (transaction_id)
);
```
