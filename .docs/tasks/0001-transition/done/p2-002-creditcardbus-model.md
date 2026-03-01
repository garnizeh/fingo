# p2-002 — creditcardbus: model, filter, order, events

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** Done  
**Assignee:** GitHub Copilot  
**Started:** 2026-02-28  
**Completed:** 2026-02-28  

---

## Goal

Create the four foundational files of the `creditcardbus` package: `model.go`, `filter.go`, `order.go`, and `event.go`.

---

## Context

These files define the pure-Go business types used throughout the domain. No `net/http`, no JSON tags, and no SQL tags are allowed in any of these files. Use `productbus` as a structural reference. All update-type fields must be pointer types.

---

## Acceptance Criteria

- [x] `business/domain/creditcardbus/model.go` exists and exports `CreditCard`, `NewCreditCard`, `UpdateCreditCard`, `Invoice`, `NewInvoice`, `UpdateInvoice`, `InvoiceStatus` constants (`Open`, `Closed`, `Paid`).
- [x] All fields in `UpdateCreditCard` and `UpdateInvoice` are pointer types.
- [x] `business/domain/creditcardbus/filter.go` exists and exports `QueryFilter` with all filter fields as pointer types.
- [x] `business/domain/creditcardbus/order.go` exists and exports `DefaultOrderBy` plus at least `OrderByID`, `OrderByName`, `OrderByUserID`, `OrderByDateCreated`.
- [x] `business/domain/creditcardbus/event.go` exists and exports `DomainName`, `ActionCreated`, `ActionUpdated`, `ActionDeleted`.
- [x] `go build ./business/domain/creditcardbus/...` exits zero.

---

## Steps

1. Create `business/domain/creditcardbus/model.go`:
   - `CreditCard` struct (value semantics): `ID`, `UserID`, `Name name.Name`, `Limit money.Money`, `ClosingDay int`, `DueDay int`, `LastFourDigits string`, `Enabled bool`, `DateCreated time.Time`, `DateUpdated time.Time`.
   - `NewCreditCard` struct: same fields except `ID`, `DateCreated`, `DateUpdated`.
   - `UpdateCreditCard` struct: `Name *name.Name`, `Limit *money.Money`, `ClosingDay *int`, `DueDay *int`, `Enabled *bool`.
   - `InvoiceStatus` type + constants `Open`, `Closed`, `Paid`.
   - `Invoice` struct: `ID`, `CreditCardID`, `ReferenceMonth time.Time`, `TotalAmount money.Money`, `Status InvoiceStatus`, `DueDate time.Time`, `DateCreated`, `DateUpdated`.
   - `NewInvoice` and `UpdateInvoice` structs.

2. Create `business/domain/creditcardbus/filter.go`:
   - `QueryFilter` struct with pointer fields: `ID *uuid.UUID`, `UserID *uuid.UUID`, `Name *name.Name`, `Enabled *bool`, `StartCreatedDate *time.Time`, `EndCreatedDate *time.Time`.
   - Add `Validate() error` method.

3. Create `business/domain/creditcardbus/order.go`:
   - Declare `orderByFields` map mapping constants to SQL column names.
   - Export `DefaultOrderBy`, `OrderByID`, `OrderByName`, `OrderByUserID`, `OrderByEnabled`, `OrderByDateCreated`.

4. Create `business/domain/creditcardbus/event.go`:
   - `const DomainName = "creditcard"`.
   - `const ActionCreated = "created"`, `ActionUpdated = "updated"`, `ActionDeleted = "deleted"`.

---

## Notes

`InvoiceStatus` follows the same pattern as other enum-like types in the project. Use a `Parse` constructor if other domains do so.
