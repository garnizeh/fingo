# p2-001 — Migration SQL: credit_cards and invoices tables

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** Done  
**Assignee:** GitHub Copilot  
**Started:** 2026-02-28  
**Completed:** 2026-02-28  

---

## Goal

Append versioned migration blocks for the `credit_cards` and `invoices` tables to the existing migration file so that Darwin applies them on next startup.

---

## Context

Darwin tracks each `-- Version: X.XX` block by checksum. Existing blocks must never be edited. New blocks are appended at the end of `business/sdk/migrate/sql/migrate.sql`. The migration must be idempotent (`CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`).

---

## Acceptance Criteria

- [x] Version `2.01` block creates table `credit_cards` with all required columns, primary key, foreign key to `users(user_id)`, and `idx_credit_cards_user` index.
- [x] Version `2.02` block creates table `invoices` with all required columns, primary key, foreign key to `credit_cards(credit_card_id)`, and the `uq_invoice_card_month` unique constraint.
- [x] Running the migration twice against the same database produces no error (idempotent).
- [x] `go build ./...` exits zero after the change.

---

## Steps

1. Open `business/sdk/migrate/sql/migrate.sql`.
2. Append version `2.01` block:

```sql
-- Version: 2.01
-- Description: Create table credit_cards
CREATE TABLE IF NOT EXISTS credit_cards (
    credit_card_id   UUID           NOT NULL,
    user_id          UUID           NOT NULL REFERENCES users(user_id),
    name             TEXT           NOT NULL,
    card_limit       NUMERIC(12, 2) NOT NULL,
    closing_day      INT            NOT NULL,
    due_day          INT            NOT NULL,
    last_four_digits CHAR(4)        NOT NULL,
    enabled          BOOLEAN        NOT NULL DEFAULT true,
    date_created     TIMESTAMPTZ    NOT NULL,
    date_updated     TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (credit_card_id),
    CONSTRAINT chk_closing_day CHECK (closing_day BETWEEN 1 AND 31),
    CONSTRAINT chk_due_day     CHECK (due_day BETWEEN 1 AND 31)
);

CREATE INDEX IF NOT EXISTS idx_credit_cards_user ON credit_cards (user_id);
```

3. Append version `2.02` block:

```sql
-- Version: 2.02
-- Description: Create table invoices
CREATE TABLE IF NOT EXISTS invoices (
    invoice_id       UUID           NOT NULL,
    credit_card_id   UUID           NOT NULL REFERENCES credit_cards(credit_card_id),
    reference_month  DATE           NOT NULL,
    total_amount     NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status           TEXT           NOT NULL DEFAULT 'open',
    due_date         TIMESTAMPTZ    NOT NULL,
    date_created     TIMESTAMPTZ    NOT NULL,
    date_updated     TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (invoice_id),
    CONSTRAINT uq_invoice_card_month UNIQUE (credit_card_id, reference_month)
);
```

4. Spin up a local PostgreSQL instance and confirm both migrations apply cleanly.
5. Run the migration a second time to confirm idempotency.

---

## Notes

The `reference_month` column is stored as `DATE` (not `TIMESTAMPTZ`) and must always be set to the first day of the billing month in UTC before insertion.
