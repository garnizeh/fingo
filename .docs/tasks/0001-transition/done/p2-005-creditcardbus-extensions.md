# p2-005 — creditcardbus extensions: otel and audit

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status: Done**
**Assignee:** unassigned
**Started:** —
**Completed: 2026-02-28**

---

## Goal

Implement the two decorator extensions for `creditcardbus`: the OpenTelemetry tracing extension and the audit trail extension.

---

## Context

Extensions wrap `creditcardbus.ExtBusiness` and add cross-cutting concerns (tracing, audit log) without modifying the core `Business` struct. Use `useraudit` and `userotel` as direct structural templates.

---

## Acceptance Criteria

- [x] `business/domain/creditcardbus/extensions/creditcardotel/creditcardotel.go` exists and implements all `ExtBusiness` methods, adding an OTel span to each.
- [x] `business/domain/creditcardbus/extensions/creditcardaudit/creditcardaudit.go` exists and calls `auditbus.Create` after `Create`, `Update`, and `Delete`.
- [x] `creditcardotel.NewExtension()` returns a `creditcardbus.Extension`.
- [x] `creditcardaudit.NewExtension(auditBus auditbus.ExtBusiness)` returns a `creditcardbus.Extension`.
- [x] `go build ./business/domain/creditcardbus/...` exits zero.

---

## Steps

1. Create `creditcardotel/creditcardotel.go`:
   - `otelBusiness` struct wrapping `creditcardbus.ExtBusiness`.
   - Each method calls `otel.AddSpan(ctx, "business.creditcardbus.<method>")`.
   - `NewExtension()` returns `creditcardbus.Extension` (a `func(ExtBusiness) ExtBusiness`).

2. Create `creditcardaudit/creditcardaudit.go`:
   - `Extension` struct with `bus creditcardbus.ExtBusiness` and `auditBus auditbus.ExtBusiness`.
   - `Create` calls `ext.bus.Create`, then calls `ext.auditBus.Create` with `ActionCreated`.
   - `Update` calls `ext.bus.Update`, then calls `ext.auditBus.Create` with `ActionUpdated`.
   - `Delete` calls `ext.bus.Delete`, then calls `ext.auditBus.Create` with `ActionDeleted`.
   - Read-only methods (`Query`, `Count`, `QueryByID`) delegate directly with no audit call.
   - `NewExtension(auditBus)` returns `creditcardbus.Extension`.

---

## Notes

Audit entries use `domain.Domain("creditcard")` for the `ObjDomain` field.
