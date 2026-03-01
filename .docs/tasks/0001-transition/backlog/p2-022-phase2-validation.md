# p2-022 — Phase 2 validation: migration, cross-domain, seed-data

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Validate that all Phase 2 domains are correctly integrated, the migration is idempotent, and seed data is available for testing.

---

## Context

Use `api/tooling/admin/main.go` for seeding and Darwin migrations for verification.

---

## Acceptance Criteria

- [ ] All 10 table create migrations (up to 2.05) are idempotent.
- [ ] `go build ./...` and `go test ./...` exit zero.
- [ ] Cross-domain logic works (e.g. transaction updates budget balance).

---

## Steps

1. Run all migrations on a fresh PostgreSQL.
2. Verify table count.
3. Add seed data for `credit_cards`, `budgets`, `transactions`, `investments` in `api/tooling/admin/commands/seed.go`.
4. Create a test user and confirm CRUD across all new domains using a simple manual check or Integration Test.
