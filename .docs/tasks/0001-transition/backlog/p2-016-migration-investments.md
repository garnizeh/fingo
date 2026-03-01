# p2-016 — Migration SQL: investments table

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Append versioned migration block for the `investments` table to the existing migration file.

---

## Context

Darwin tracks each block by checksum. Existing blocks must never be edited. New blocks are appended at the end of `business/sdk/migrate/sql/migrate.sql`.

---

## Acceptance Criteria

- [ ] Version `2.05` block creates table `investments`.
- [ ] Running the migration twice produces no error.
- [ ] `go build ./...` exits zero.

---

## Steps

1. Open `business/sdk/migrate/sql/migrate.sql`.
2. Append version `2.05`:
```sql
-- Version: 2.05
-- Description: Create table investments
CREATE TABLE IF NOT EXISTS investments (
    investment_id  UUID           NOT NULL,
    user_id        UUID           NOT NULL REFERENCES users(user_id),
    asset_class    TEXT           NOT NULL,
    ticker         TEXT           NOT NULL,
    quantity       NUMERIC(12, 4) NOT NULL,
    avg_cost       NUMERIC(12, 2) NOT NULL,
    currency       TEXT           NOT NULL,
    date_created   TIMESTAMPTZ    NOT NULL,
    date_updated   TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (investment_id)
);
```
