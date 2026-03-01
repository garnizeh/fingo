# p2-007 — creditcardapp: model, filter, order, handlers, routes

**Design Doc:** `.docs/design-docs/0001-transition.md`
**Phase:** 2 — Domain Implementation
**Status:** To Do
**Assignee:** unassigned
**Started:** —
**Completed:** —

---

## Goal

Implement the full HTTP application layer for the Credit Card domain under `app/domain/creditcardapp/`.

---

## Context

The app layer translates HTTP ↔ business types only. No business logic lives here. Handlers must return `web.Encoder`. DTOs have JSON tags. Use `productapp` as a structural template.

---

## Acceptance Criteria

- [x] `app/domain/creditcardapp/model.go` exists with `AppCreditCard`, `NewCreditCard` (input), and `UpdateCreditCard` (input) DTOs carrying `json` tags; includes `toAppCreditCard`, `toBusNewCreditCard`, `toBusUpdateCreditCard` converters.
- [x] `app/domain/creditcardapp/filter.go` exists and parses `user_id`, `name`, `enabled`, `start_date`, `end_date` query params into `creditcardbus.QueryFilter`.
- [x] `app/domain/creditcardapp/order.go` exists and parses `orderBy` query param into `order.By`.
- [x] `app/domain/creditcardapp/creditcardapp.go` exists with handler methods: `create`, `update`, `delete`, `query`, `queryByID`.
- [x] `app/domain/creditcardapp/route.go` exists; registers routes under `/v1/credit-cards` with appropriate auth middleware.
- [x] `go build ./app/domain/creditcardapp/...` exits zero.

---

## Steps

1. Create `model.go`:
   - `AppCreditCard` with all fields serialised to snake_case JSON.
   - `NewCreditCard` input DTO (fields: `name`, `limit`, `closing_day`, `due_day`, `last_four_digits`).
   - `UpdateCreditCard` input DTO (all pointer fields).
   - Converter functions.

2. Create `filter.go`: parse URL query params; return `errs.New(errs.InvalidArgument, ...)` on parse errors.

3. Create `order.go`: map string `fieldName` to `order.By` using `creditcardbus.OrderByXxx` constants.

4. Create `creditcardapp.go`:
   - `app` struct with `creditCardBus creditcardbus.ExtBusiness`.
   - `create`: decode body → `NewCreditCard` → `toBusNewCreditCard` → `bus.Create` → `toAppCreditCard`.
   - `update`: parse `credit_card_id` path param → `bus.QueryByID` → `toBusUpdateCreditCard` → `bus.Update` → `toAppCreditCard`.
   - `delete`: parse `credit_card_id` → `bus.QueryByID` → `bus.Delete`.
   - `query`: parse filter + order + page → `bus.Query` → map to `[]AppCreditCard`.
   - `queryByID`: parse `credit_card_id` → `bus.QueryByID` → `toAppCreditCard`.

5. Create `route.go`:
   ```go
   func Routes(app *web.App, cfg Config) {
       const version = "v1"
       authen := mid.Authenticate(cfg.AuthClient)
       a := newApp(cfg.CreditCardBus)
       app.HandlerFunc(http.MethodPost,   version, "/credit-cards",          a.create,      authen)
       app.HandlerFunc(http.MethodPut,    version, "/credit-cards/{id}",     a.update,      authen)
       app.HandlerFunc(http.MethodDelete, version, "/credit-cards/{id}",     a.delete,      authen)
       app.HandlerFunc(http.MethodGet,    version, "/credit-cards",          a.query,       authen)
       app.HandlerFunc(http.MethodGet,    version, "/credit-cards/{id}",     a.queryByID,   authen)
   }
   ```

---

## Notes

`mid.GetSubjectID(ctx)` provides the authenticated user's UUID for ownership checks. The `delete` handler returns HTTP 204 with no body.
