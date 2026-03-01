# p2-020 — investmentapp: DTOs, handlers, routes for investments and vinvestments

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the HTTP application layer for regular `investments` and the virtual `vinvestments`.

---

## Context

Handlers translate HTTP ↔ business types; return `web.Encoder`.

---

## Acceptance Criteria

- [ ] `app/domain/investmentapp/` exists with all app layer files.
- [ ] `app/domain/vinvestmentapp/` exists with read-only portfolio handlers.
- [ ] Both domains registered with `/v1/investments` and `/v1/portfolio`.

---

## Steps

1. Create DTOs and converters for `investmentapp` and `vinvestmentapp`.
2. Implement handlers: `create`, `update`, `delete`, `query`, `queryByID` for investments.
3. Implement `query` and `queryByID` for vinvestments.
4. Configure routes in `route.go`.
