# p2-008 — Mux wiring and DI for creditcardbus in main.go

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** Done
**Assignee:** GitHub Copilot
**Started:** 2026-02-28
**Completed:** 2026-02-28

---

## Goal

Wire `creditcardbus` into the service entry point (`api/services/fingo/main.go`) and register its routes via the mux configuration.

---

## Context

The five-step DI pattern from design doc Section 2 must be followed exactly. `CreditCardBus` must be added to `mux.BusConfig` and the route registration must be added to the appropriate build file.

---

## Acceptance Criteria

- [x] `app/sdk/mux/mux.go` has `CreditCardBus creditcardbus.ExtBusiness` field in `BusConfig`.
- [x] `creditcardapp.Routes` is called inside the mux route-wiring function with `cfg.BusConfig.CreditCardBus`.
- [x] `api/services/fingo/main.go` constructs `creditcarddb.Store`, the two extensions, and `creditcardbus.NewBusiness` in the correct dependency order (after `auditBus` and `userBus`).
- [x] `creditCardBus` is passed into `mux.BusConfig`.
- [x] `go build ./...` exits zero.
- [x] `go vet ./...` exits zero.

---

## Steps

1. Add import and field to `app/sdk/mux/mux.go`:
   ```go
   CreditCardBus creditcardbus.ExtBusiness
   ```

2. In the mux route-wiring function (`Routes` or equivalent), add:
   ```go
   creditcardapp.Routes(app, creditcardapp.Config{
       AuthClient:    cfg.AuthClient,
       CreditCardBus: cfg.BusConfig.CreditCardBus,
   })
   ```

3. In `api/services/fingo/main.go` inside `run()`:
   ```go
   creditCardOtelExt  := creditcardotel.NewExtension()
   creditCardAuditExt := creditcardaudit.NewExtension(auditBus)
   creditCardStorage  := creditcarddb.NewStore(log, db)
   creditCardBus      := creditcardbus.NewBusiness(log, userBus, delegate,
                             creditCardStorage, creditCardOtelExt, creditCardAuditExt)
   ```

4. Pass `creditCardBus` into the `BusConfig` struct.

5. Run `go build ./...` and `go vet ./...` to confirm zero errors.

---

## Notes

This task depends on p2-003 (bus), p2-004 (store), p2-005 (extensions), and p2-007 (app layer) all being complete.
