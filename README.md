# FinGo — Personal Finance and Investment Manager

FinGo is a **high-performance Personal Finance and Investment Manager** designed to run on a single VPS. It provides tools for managing credit cards, investments, budgets, and real-time cash flow monitoring.

The project is built using the **Ardan Labs `service` pattern**, a layered Go monorepo architecture that isolates business domains from transport, persistence, and infrastructure concerns. This project is derived from the [Ardan Labs Service](https://github.com/ardanlabs/service) template.

---

## 🚀 Core Features

- **Credit Cards Management:** Track invoice cycles, statement closing dates, minimum payments, and spending limits.
- **Investment Portfolio:** Manage Stocks, Cryptocurrencies, and Fixed Income instruments with real-time valuation and gain/loss calculations.
- **Cash Flow & Budgeting:** Accounts Payable/Receivable, recurring transactions, and envelope-based budget management.
- **Real-time Dashboard:** Aggregated net-worth, cash position, and P&L monitoring via WebSocket/SSE.

---

## 🛠 Tech Stack

- **Language:** Go 1.26+ (Standard Library focused)
- **Architecture:** Layered Monorepo (Ardan Labs Service Pattern)
- **Database:** PostgreSQL (using `sqlx` and `pgx`)
- **Real-time:** SSE (Server-Sent Events) and WebSockets
- **Infrastructure:** Docker, Kubernetes (Kind), Helm
- **Observability:** OpenTelemetry (Tracing), Prometheus (Metrics)
- **Frontend:** React (Vite-based)

---

## 📂 Project Structure

```text
├── api/             # Entry points for services (auth, metrics, fingo)
├── app/             # Application layer (HTTP handlers, DTOs, routing)
├── business/        # Core business logic and domain entities
├── foundation/      # Framework primitives (logger, web, otel)
├── zarf/            # Infrastructure manifests (Docker, K8s, Helm)
└── .docs/           # Project documentation and design docs
```

---

## 🏁 Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)
- [Brew](https://brew.sh/) (recommended for tooling)

### Installation

1. Install required tooling and dependencies:
   ```bash
   make dev-brew
   make dev-docker
   make dev-gotooling
   ```

2. Verify the installation by running tests:
   ```bash
   make test
   ```

### Running the Project

Start the local development environment using KIND (Kubernetes in Docker):

```bash
make dev-up
make dev-update-apply
```

Generate an authentication token:
```bash
make token
export TOKEN=<generated_token>
```

Check cluster status:
```bash
make dev-status
```

---

## 📖 Documentation

Detailed documentation, design documents, and task lists are located in the [/.docs](/.docs) directory.

- [Copilot Instructions](.github/copilot-instructions.md): Guidelines for developing within this repository.
- [Design Docs](.docs/design-docs/): Architectural decisions and future roadmap.

---

## 📄 License

Distributed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for more information.
