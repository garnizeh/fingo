# 0001 — FinGo: Migration Design Document

**Status:** Draft  
**Author:** Engineering  
**Last Updated:** 2026-02-28  
**Reference:** `.docs/copilot-instructions.md`

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State Audit](#2-current-state-audit)
3. [Target Architecture](#3-target-architecture)
4. [Migration Phases](#4-migration-phases)
   - [Phase 1 — Identity and Rebranding](#phase-1--identity-and-rebranding)
   - [Phase 2 — Domain Implementation](#phase-2--domain-implementation)
   - [Phase 3 — Real-time Dashboard](#phase-3--real-time-dashboard)
   - [Phase 4 — Legacy Cleanup](#phase-4--legacy-cleanup)
5. [Dependency Map](#5-dependency-map)
6. [Risk Register](#6-risk-register)
7. [Definition of Done](#7-definition-of-done)

---

## 1. Executive Summary

This document describes the complete engineering plan to migrate the cloned `ardanlabs/service` boilerplate into **FinGo**, a personal finance and investment manager.

The boilerplate provides a production-ready skeleton: layered architecture, PostgreSQL persistence via `sqlx`, OpenTelemetry tracing, JWT authentication, structured logging, and Docker/Kubernetes infrastructure. The migration strategy is **additive first, destructive last**: we build all FinGo domains alongside the existing `productbus`/`homebus` code, validate functionality, then remove the boilerplate artifacts. This avoids a "big bang" rewrite and allows the CI pipeline to remain green throughout.

**Non-negotiable constraints throughout all phases:**
- Every output (code, comments, logs) must be in **English**.
- If any step is ambiguous, stop and ask rather than guess.

---

## 2. Current State Audit

### What Exists Today

| Layer | Package | Status |
|---|---|---|
| Foundation | `foundation/web`, `foundation/logger`, `foundation/otel` | ✅ Keep — stable, do not modify |
| Business SDK | `business/sdk/sqldb`, `delegate`, `migrate`, `order`, `page` | ✅ Keep — do not modify |
| Auth service | `api/services/auth/` | ✅ Keep — do not touch |
| Metrics service | `api/services/metrics/` | ✅ Keep — do not touch |
| User domain | `business/domain/userbus/` | ✅ Keep |
| Audit domain | `business/domain/auditbus/` | ✅ Keep — used by all new domains |
| Product domain | `business/domain/productbus/` | 🔁 Reference only — will be deleted in Phase 4 |
| Home domain | `business/domain/homebus/` | 🔁 Template for new domains |
| Virtual product | `business/domain/vproductbus/` | 🔁 Template for `vinvestmentbus` |
| Sales service | `api/services/sales/` | 🔁 Rename to `fingo` |
| Migrations | `business/sdk/migrate/sql/migrate.sql` | ⚠️ Extend — append new versions |
| Docker Compose | `zarf/compose/docker_compose.yaml` | ⚠️ Extend — add fingo service |

### Current Dependency Injection Pattern (in `main.go`)

The boilerplate already demonstrates the correct DI pattern we will follow for every new domain. Understanding this flow is critical before writing any new code:

```go
// Step 1: Create the delegate (cross-domain event bus, no direct imports)
delegate := delegate.New(log)

// Step 2: Build stores (DB layer), then wrap with cache if needed
userStorage := usercache.NewStore(log, userdb.NewStore(log, db), time.Minute)

// Step 3: Build extensions (OTel, Audit) — each is a func(ExtBusiness) ExtBusiness
userOtelExt  := userotel.NewExtension()
userAuditExt := useraudit.NewExtension(auditBus) // receives auditBus as dependency

// Step 4: Compose the Business with extensions applied in reverse order
userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

// Step 5: Pass composed buses into mux.BusConfig
cfgMux := mux.Config{
    BusConfig: mux.BusConfig{
        UserBus: userBus,
        // ...
    },
}
```

Every new FinGo domain will follow this exact five-step pattern.

### Current Migration File

Migrations are versioned SQL blocks inside a single file using the darwin library comment convention:

```sql
-- Version: 1.01
-- Description: Create table users
CREATE TABLE users ( ... );

-- Version: 1.02
-- Description: Create table products
CREATE TABLE products ( ... );
```

**Important:** Darwin checksums each version block. Never edit an already-applied version. Always append new version blocks to extend the schema.

---

## 3. Target Architecture

### Final Domain Map

```
business/domain/
├── auditbus/          ✅ unchanged
├── userbus/           ✅ unchanged
├── creditcardbus/     🆕 Phase 2
├── investmentbus/     🆕 Phase 2
├── budgetbus/         🆕 Phase 2
├── transactionbus/    🆕 Phase 2
├── dashboardbus/      🆕 Phase 3  (read-only, no mutations)
└── vinvestmentbus/    🆕 Phase 3  (read-only portfolio view)
```

### Final Service Map

```
api/services/
├── auth/             ✅ unchanged
├── metrics/          ✅ unchanged
└── fingo/            🔁 renamed from sales/
```

---

## 4. Migration Phases

---

### Phase 1 — Identity and Rebranding

**Goal:** Replace every reference to the `ardanlabs/service` sample project with the `github.com/garnizeh/fingo` identity. This phase is purely mechanical — no logic changes.

**Why this matters:** The module path `github.com/garnizeh/fingo` is embedded in every import statement, Docker image tag, environment variable prefix, and Kubernetes label. Establishing the correct identity now prevents cascading rename operations later.

---

#### 1.1 — Module Path

The module path is already correct in `go.mod`:

```
module github.com/garnizeh/fingo
```

Verify that no file in the repository still imports `github.com/ardanlabs/service`:

```bash
grep -r "ardanlabs/service" --include="*.go" .
```

If any matches are found, replace them with `github.com/garnizeh/fingo`.

---

#### 1.2 — Rename `api/services/sales/` to `api/services/fingo/`

```bash
mv api/services/sales api/services/fingo
```

Inside `api/services/fingo/main.go`, update the three places that reference the old service identity:

```go
// BEFORE
log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", otel.GetTraceID, events)

cfg := struct { ... }{ Version: conf.Version{ Build: tag, Desc: "Sales" } }

const prefix = "SALES"
```

```go
// AFTER
log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "FINGO", otel.GetTraceID, events)

cfg := struct { ... }{ Version: conf.Version{ Build: tag, Desc: "FinGo" } }

const prefix = "FINGO"
```

The `conf` prefix change means all environment variables shift from `SALES_*` to `FINGO_*` (e.g., `FINGO_DB_HOST`, `FINGO_WEB_API_HOST`). Update every environment definition in Docker Compose and Kubernetes manifests in the same commit.

---

#### 1.3 — Rename the Docker service name

In `api/services/fingo/build/` (formerly `sales/build/`), no logic changes are needed — only the service name in log lines and the `Tempo.ServiceName` config default:

```go
// In the cfg struct inside main.go:
Tempo struct {
    Host        string  `conf:"default:tempo:4317"`
    ServiceName string  `conf:"default:fingo"`   // was "sales"
    Probability float64 `conf:"default:0.05"`
}
```

---

#### 1.4 — Infrastructure files

```bash
mv zarf/docker/dockerfile.sales zarf/docker/dockerfile.fingo
```

Update `zarf/compose/docker_compose.yaml`:
- Rename service `sales` → `fingo`.
- Rename `sales-system-network` → `fingo-network` (update all `networks:` references in the file).
- Update image references: `localhost/ardanlabs/sales:0.0.1` → `localhost/garnizeh/fingo:0.0.1`.
- Update `init-migrate-seed` env vars from `SALES_DB_*` to `FINGO_DB_*`.

Update the `makefile`:

```makefile
# BEFORE
SERVICE_IMAGE_NAME := ardanlabs/sales
```

```makefile
# AFTER
SERVICE_IMAGE_NAME := garnizeh/fingo
```

---

#### 1.5 — Validation

After all renames, the project must still compile and all existing tests must pass:

```bash
go build ./...
go test ./...
```

If either command fails, fix the compilation error before proceeding to Phase 2.

---

### Phase 2 — Domain Implementation

**Goal:** Implement the three core FinGo business domains — Credit Cards, Investments, and Budget/Transactions — using the Ardan Labs layered pattern. Each domain is self-contained and follows the 19-step checklist in `copilot-instructions.md` Section 8.

The recommended implementation order is: **creditcardbus → budgetbus → transactionbus → investmentbus**. This order reflects dependency complexity: credit cards are the simplest domain with no cross-domain dependencies; investments are the most complex because they require the virtual view pattern.

---

#### 2.1 — Credit Card Domain

##### Why credit cards first?

Credit cards are the most straightforward domain: a card belongs to a user, and invoices are grouped by billing month. There are no external API calls and no cross-domain dependencies beyond `userbus`. This makes it the ideal domain to validate the full stack (migration → bus → store → app → route → DI) before tackling more complex domains.

##### Migration

Append to `business/sdk/migrate/sql/migrate.sql`:

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

-- Version: 2.02
-- Description: Create table invoices
CREATE TABLE IF NOT EXISTS invoices (
    invoice_id        UUID           NOT NULL,
    credit_card_id    UUID           NOT NULL REFERENCES credit_cards(credit_card_id),
    reference_month   DATE           NOT NULL,  -- stored as first day of month
    total_amount      NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status            TEXT           NOT NULL DEFAULT 'open',
    due_date          TIMESTAMPTZ    NOT NULL,
    date_created      TIMESTAMPTZ    NOT NULL,
    date_updated      TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (invoice_id),
    CONSTRAINT uq_invoice_card_month UNIQUE (credit_card_id, reference_month)
);
```

##### Business layer file structure

```
business/domain/creditcardbus/
├── model.go          -- CreditCard, NewCreditCard, UpdateCreditCard, Invoice, ...
├── filter.go         -- QueryFilter (all pointer fields)
├── order.go          -- DefaultOrderBy, OrderByXxx constants
├── event.go          -- DomainName, ActionCreated, ActionUpdated, ActionDeleted
├── creditcardbus.go  -- Storer interface, ExtBusiness interface, Business struct
├── creditcardbus_test.go
├── testutil.go
├── stores/
│   └── creditcarddb/
│       ├── creditcarddb.go  -- SQL implementation
│       └── model.go         -- dbCreditCard, toDBCreditCard, toCreditCard
└── extensions/
    ├── creditcardotel/
    │   └── creditcardotel.go
    └── creditcardaudit/
        └── creditcardaudit.go
```

##### `creditcardbus.go` — key interfaces

```go
package creditcardbus

var (
    ErrNotFound  = errors.New("credit card not found")
    ErrCardLimit = errors.New("credit card limit must be positive")
)

// Storer is the persistence contract. Only this package and its tests know it exists.
type Storer interface {
    NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
    Create(ctx context.Context, cc CreditCard) error
    Update(ctx context.Context, cc CreditCard) error
    Delete(ctx context.Context, cc CreditCard) error
    Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]CreditCard, error)
    Count(ctx context.Context, filter QueryFilter) (int, error)
    QueryByID(ctx context.Context, creditCardID uuid.UUID) (CreditCard, error)
}

// ExtBusiness is the public contract consumed by the app layer and extensions.
type ExtBusiness interface {
    NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
    Create(ctx context.Context, actorID uuid.UUID, nc NewCreditCard) (CreditCard, error)
    Update(ctx context.Context, actorID uuid.UUID, cc CreditCard, uc UpdateCreditCard) (CreditCard, error)
    Delete(ctx context.Context, actorID uuid.UUID, cc CreditCard) error
    Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]CreditCard, error)
    Count(ctx context.Context, filter QueryFilter) (int, error)
    QueryByID(ctx context.Context, creditCardID uuid.UUID) (CreditCard, error)
}
```

##### Audit extension pattern

The audit extension follows the exact same pattern as `useraudit`. It wraps `ExtBusiness` and calls `auditBus.Create` after each mutating operation:

```go
// business/domain/creditcardbus/extensions/creditcardaudit/creditcardaudit.go
package creditcardaudit

type Extension struct {
    bus      creditcardbus.ExtBusiness
    auditBus auditbus.ExtBusiness
}

func NewExtension(auditBus auditbus.ExtBusiness) creditcardbus.Extension {
    return func(bus creditcardbus.ExtBusiness) creditcardbus.ExtBusiness {
        return &Extension{bus: bus, auditBus: auditBus}
    }
}

func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, nc creditcardbus.NewCreditCard) (creditcardbus.CreditCard, error) {
    cc, err := ext.bus.Create(ctx, actorID, nc)
    if err != nil {
        return creditcardbus.CreditCard{}, err
    }

    if _, err := ext.auditBus.Create(ctx, auditbus.NewAudit{
        ObjID:     cc.ID,
        ObjDomain: domain.Domain("creditcard"),
        ObjName:   cc.Name,
        ActorID:   actorID,
        Action:    creditcardbus.ActionCreated,
        Data:      nc,
        Message:   "credit card created",
    }); err != nil {
        return creditcardbus.CreditCard{}, err
    }

    return cc, nil
}
```

##### App layer — `creditcardapp`

The app layer translates HTTP ↔ business types:

```go
// app/domain/creditcardapp/creditcardapp.go
func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
    var app NewCreditCard
    if err := web.Decode(r, &app); err != nil {
        return errs.New(errs.InvalidArgument, err)
    }

    nc := creditcardbus.NewCreditCard{
        UserID:         mid.GetSubjectID(ctx),
        Name:           mustParseName(app.Name),
        Limit:          mustParseMoney(app.Limit),
        ClosingDay:     app.ClosingDay,
        DueDay:         app.DueDay,
        LastFourDigits: app.LastFourDigits,
    }

    cc, err := a.creditCardBus.Create(ctx, mid.GetSubjectID(ctx), nc)
    if err != nil {
        if errors.Is(err, creditcardbus.ErrCardLimit) {
            return errs.New(errs.InvalidArgument, creditcardbus.ErrCardLimit)
        }
        return errs.Errorf(errs.Internal, "create: cc[%+v]: %s", cc, err)
    }

    return toAppCreditCard(cc)
}
```

##### Mux wiring

Add `CreditCardBus` to `BusConfig` in `app/sdk/mux/mux.go`:

```go
type BusConfig struct {
    AuditBus      auditbus.ExtBusiness
    UserBus       userbus.ExtBusiness
    CreditCardBus creditcardbus.ExtBusiness  // ← ADD
    // ... HomeBus, ProductBus (still present until Phase 4)
}
```

##### Dependency injection in `main.go`

```go
// api/services/fingo/main.go — inside run()

creditCardOtelExt  := creditcardotel.NewExtension()
creditCardAuditExt := creditcardaudit.NewExtension(auditBus)
creditCardStorage  := creditcarddb.NewStore(log, db)
creditCardBus      := creditcardbus.NewBusiness(log, userBus, delegate, creditCardStorage,
                          creditCardOtelExt, creditCardAuditExt)
```

Then pass `creditCardBus` into the `BusConfig`.

---

#### 2.2 — Budget and Transaction Domains

Follow the identical 19-step checklist for `budgetbus` and `transactionbus`. Key differences from `creditcardbus`:

- `Transaction` has a nullable `BudgetID *uuid.UUID` — not all transactions belong to a budget.
- `Transaction.Direction` is a value object (`debit` | `credit`) — create it in `business/types/direction/`.
- When a transaction is created, if a `BudgetID` is present, the `budgetbus` must be called to update the budget's running total. This cross-domain call goes through the **delegate**, not a direct import.

```go
// business/domain/transactionbus/event.go
const (
    DomainName    = "transaction"
    ActionCreated = "created"
)

// business/domain/budgetbus/budgetbus.go — registers for the event:
func (b *Business) registerDelegateFunctions() {
    if b.delegate != nil {
        b.delegate.Register(transactionbus.DomainName, transactionbus.ActionCreated,
            b.actionTransactionCreated)
    }
}
```

---

#### 2.3 — Investment Domain

The investment domain has two layers:
1. **`investmentbus`** — manages the user's holdings (ticker, quantity, average cost). This is a standard CRUD domain.
2. **`vinvestmentbus`** — a read-only virtual domain that joins `investments` with live price data to produce a portfolio view with mark-to-market valuations. Use `vproductbus` as the structural template.

New value objects required before implementing `investmentbus`:

| Type | Location | Description |
|---|---|---|
| `currency.Currency` | `business/types/currency/` | `BRL`, `USD`, `EUR` — with `Parse` and `String` |
| `ticker.Ticker` | `business/types/ticker/` | Validated uppercase string, 1–10 chars |
| `assetclass.AssetClass` | `business/types/assetclass/` | `stock`, `crypto`, `fixed_income` |

Each type must have a `Parse(s string) (Type, error)` constructor and a `String() string` method.

---

### Phase 3 — Real-time Dashboard

**Goal:** Expose a `/v1/dashboard/stream` endpoint that pushes aggregated financial snapshots to connected clients using Server-Sent Events (SSE).

---

#### 3.1 — Dashboard Bus (read-only)

`dashboardbus` is a **read-only** business domain. It has no `Create`, `Update`, or `Delete` methods. It contains one `Storer` interface and no `ExtBusiness` interface — there is nothing to decorate.

```go
// business/domain/dashboardbus/dashboardbus.go
package dashboardbus

type Storer interface {
    NetWorthSnapshot(ctx context.Context, userID uuid.UUID) (NetWorthSnapshot, error)
    CashFlowSummary(ctx context.Context, userID uuid.UUID, month time.Time) (CashFlowSummary, error)
    InvestmentPortfolio(ctx context.Context, userID uuid.UUID) ([]PositionSummary, error)
}

type Business struct {
    log    *logger.Logger
    storer Storer
}

func NewBusiness(log *logger.Logger, storer Storer) *Business {
    return &Business{log: log, storer: storer}
}

func (b *Business) NetWorthSnapshot(ctx context.Context, userID uuid.UUID) (NetWorthSnapshot, error) {
    ctx, span := otel.AddSpan(ctx, "business.dashboardbus.networth")
    defer span.End()
    return b.storer.NetWorthSnapshot(ctx, userID)
}
```

The SQL backing `NetWorthSnapshot` should be a CTE that aggregates across `credit_cards`, `investments`, and `transactions`:

```sql
-- business/domain/dashboardbus/stores/dashboarddb/dashboarddb.go (embedded query)
WITH
    credit_card_balance AS (
        SELECT user_id, COALESCE(SUM(total_amount), 0) AS total
        FROM invoices
        JOIN credit_cards USING (credit_card_id)
        WHERE user_id = :user_id AND status = 'open'
        GROUP BY user_id
    ),
    investment_value AS (
        SELECT user_id, COALESCE(SUM(quantity * avg_cost), 0) AS total
        FROM investments
        WHERE user_id = :user_id
        GROUP BY user_id
    ),
    cash_balance AS (
        SELECT user_id,
               COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END), 0) AS total
        FROM transactions
        WHERE user_id = :user_id
        GROUP BY user_id
    )
SELECT
    cb.total  AS credit_card_debt,
    iv.total  AS investment_value,
    ca.total  AS cash_balance,
    iv.total + ca.total - cb.total AS net_worth
FROM credit_card_balance cb, investment_value iv, cash_balance ca;
```

---

#### 3.2 — SSE streaming handler

SSE handlers cannot use the standard `web.Encoder` return pattern because they take streaming control of the `http.ResponseWriter`. Register them via a raw `http.HandlerFunc` adaptor in `route.go`:

```go
// app/domain/dashboardapp/route.go
func Routes(app *web.App, cfg Config) {
    const version = "v1"

    authen := mid.Authenticate(cfg.AuthClient)

    a := newApp(cfg.DashBus)

    // Note: HandlerFunc is the standard Encoder-based handler.
    // For SSE we attach the raw http.HandlerFunc wrapper directly on the mux.
    app.HandleFunc(http.MethodGet, version, "/dashboard/stream",
        http.HandlerFunc(a.stream), authen)
}
```

The handler itself:

```go
// app/domain/dashboardapp/dashboardapp.go
func (a *app) stream(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming not supported", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no")

    ctx := r.Context()
    userID := mid.GetSubjectID(ctx)

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            snapshot, err := a.dashBus.NetWorthSnapshot(ctx, userID)
            if err != nil {
                fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
                flusher.Flush()
                return
            }
            data, _ := json.Marshal(toAppSnapshot(snapshot))
            fmt.Fprintf(w, "data: %s\n\n", data)
            flusher.Flush()
        }
    }
}
```

---

#### 3.3 — Virtual Investment View (`vinvestmentbus`)

For mark-to-market portfolio valuation, use the `vproductbus` package as a direct structural template. The virtual bus:
- Has a read-only `Storer` interface.
- Joins the `investments` table with live prices (from a price-feed adapter injected at construction time).
- Never persists the calculated current value — it is always computed at query time.

---

### Phase 4 — Legacy Cleanup

**Goal:** Remove all boilerplate artifacts that were never part of FinGo. **This phase must only begin after all Phase 2 and Phase 3 domain tests pass on CI.**

---

#### 4.1 — Decommission Product and Home domains

Delete the following packages entirely:

```bash
rm -rf business/domain/productbus/
rm -rf business/domain/homebus/
rm -rf business/domain/vproductbus/
rm -rf app/domain/productapp/
rm -rf app/domain/homeapp/
```

Remove the corresponding imports from `app/sdk/mux/mux.go` and `api/services/fingo/main.go`. Remove the route registrations from `api/services/fingo/build/crud.go` (or whichever build file registers product/home routes).

---

#### 4.2 — Remove product and home migrations

> **Do NOT delete the existing migration versions.** Darwin has already applied versions 1.02 (`products`), 1.03 (`view_products`), and 1.04 (`homes`) to any environment that ran the boilerplate migrations. Deleting them from the SQL file will cause a checksum mismatch.

Instead, append **drop** migrations:

```sql
-- Version: 4.01
-- Description: Drop boilerplate tables no longer used by FinGo
DROP VIEW  IF EXISTS view_products;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS homes;
```

---

#### 4.3 — Final module tidy

```bash
go build ./...
go vet ./...
go test ./...
go mod tidy
```

All four commands must exit zero before this phase is marked complete.

---

## 5. Dependency Map

The diagram below shows the construction order required in `main.go`. A bus that appears lower in the list depends on a bus that appears higher:

```
delegate        (no dependencies)
    │
auditBus        (delegate, auditdb)
    │
userBus         (delegate, usercache → userdb, otel ext, audit ext)
    │
creditCardBus   (delegate, userBus, creditcarddb, otel ext, audit ext)
    │
budgetBus       (delegate, userBus, budgetdb, otel ext, audit ext)
    │
transactionBus  (delegate, userBus, budgetBus, transactiondb, otel ext, audit ext)
    │
investmentBus   (delegate, userBus, investmentdb, otel ext, audit ext)
    │
dashboardBus    (dashboarddb — read-only aggregate, no mutations)
vinvestmentBus  (vinvestmentdb — read-only view, no mutations)
```

**Rule:** Never pass a bus as a dependency to another bus that appears above it in this map. If you find yourself needing to do so, the correct solution is the **delegate pattern** — not a direct import.

---

## 6. Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Darwin checksum mismatch | Low | High | Never edit applied migration versions. Only append. |
| SSE connection scaling (many concurrent clients) | Medium | Medium | Use a fanout channel from a single DB poller rather than one ticker per connection when client count exceeds 100. |
| Cross-domain import cycle | Low | Medium | Use the delegate pattern. If a direct import feels necessary, stop and ask. |
| Dropping `products`/`homes` tables before data is migrated | Low | High | Drop migrations (Phase 4) only after confirming no application data exists in those tables. |

---

## 7. Definition of Done

A phase is complete only when **all** of the following are true:

- [ ] `go fmt ./...` exits zero.
- [ ] `go fix ./...` exits zero.
- [ ] `go build ./...` exits zero.
- [ ] `go vet ./...` exits zero.
- [ ] `go test ./...` exits zero.
- [ ] The migration runs against a fresh PostgreSQL instance without errors.
- [ ] The migration is **idempotent** — running it twice against the same database produces no error.
- [ ] Every new domain handler is covered by at least one integration test in `app/sdk/apitest`.
- [ ] The Docker Compose stack starts cleanly: `database` → `init-migrate-seed` → `fingo`.
- [ ] All environment variable names use the `FINGO_` prefix.

---

> **Reference:** Consult `.docs/copilot-instructions.md` Section 13 for the complete naming-conventions table during every step of this transition.
