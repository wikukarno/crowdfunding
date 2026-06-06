# Backend Crowdfunding

[![CI](https://github.com/wikukarno/crowdfunding/actions/workflows/ci.yml/badge.svg)](https://github.com/wikukarno/crowdfunding/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wikukarno/crowdfunding)](https://goreportcard.com/report/github.com/wikukarno/crowdfunding)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A REST API for crowdfunding campaigns and donations, written in Go. Users can
register, create campaigns, and donate through a payment gateway. Built with
Gin, GORM, and PostgreSQL, with JWT authentication and Midtrans for payments.

## Features

- **Users** — register, login, email availability check, avatar upload
- **Campaigns** — list, detail, create, update, image upload (with ownership checks)
- **Transactions** — create donations, list per campaign and per user
- **Payments** — Midtrans Snap integration with a webhook for status updates
- **Auth** — stateless JWT, enforced via middleware

## Highlights

- **UUID identifiers** — every resource uses a non-sequential UUID, so IDs can't be guessed or enumerated from the API.
- **Compile-time dependency injection** with google/wire — no runtime reflection or service locator.
- **Atomic payments** — transaction status and campaign totals are updated inside a single database transaction.
- **Cloud storage** — uploads go to Cloudflare R2 (S3-compatible) with a local-disk fallback for development.
- **Operability** — graceful shutdown, request-id tracing, structured JSON logs, and database connection retry on startup.

## Architecture

The code is organised by domain. Each domain owns its entity, repository
(database access), and service (business logic). HTTP handlers depend only on
the service interfaces, and dependencies are wired together with
[google/wire](https://github.com/google/wire).

```
config/        Environment configuration
database/      GORM connection setup
auth/          JWT generation and validation
user/          User domain (entity, repository, service)
campaign/      Campaign domain
transaction/   Transaction domain
payment/       Midtrans gateway
storage/       File uploads (Cloudflare R2 or local disk)
handler/       Gin HTTP handlers, auth middleware, request logging
helper/        Shared API response helpers
migrations/    SQL schema migrations
docs/          Generated Swagger/OpenAPI spec
wire.go        Dependency graph (compile-time DI)
```

Request flow: `handler → service (interface) → repository (interface) → PostgreSQL`.
Because services and repositories are interfaces, the business logic is unit
tested with in-memory mocks and no database.

## Getting Started

### Run with Docker (recommended)

```bash
cp .env.example .env          # set JWT_SECRET_KEY at minimum
docker compose up --build
```

This starts PostgreSQL and the API. The service listens on `http://localhost:8080`.
If `8080` or `5432` are already in use, override `HOST_APP_PORT` / `HOST_DB_PORT`
in your `.env`.

### Run locally

Requires Go 1.25+ and a running PostgreSQL instance.

```bash
cp .env.example .env          # adjust DB credentials and JWT_SECRET_KEY
make migrate-up DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/crowdfunding?sslmode=disable"
make run
```

Migrations use [golang-migrate](https://github.com/golang-migrate/migrate). If
you don't have the CLI:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Configuration

All configuration comes from environment variables (a local `.env` is loaded
when present).

| Variable              | Description                       | Default       |
| --------------------- | --------------------------------- | ------------- |
| `APP_PORT`            | HTTP port                         | `8080`        |
| `DB_HOST`             | PostgreSQL host                   | `127.0.0.1`   |
| `DB_PORT`             | PostgreSQL port                   | `5432`        |
| `DB_USER`             | PostgreSQL user                   | `postgres`    |
| `DB_PASSWORD`         | PostgreSQL password               | _(empty)_     |
| `DB_NAME`             | Database name                     | `crowdfunding`|
| `JWT_SECRET_KEY`      | Signing key for JWT (**required**)| —             |
| `MIDTRANS_SERVER_KEY` | Midtrans server key               | —             |
| `MIDTRANS_CLIENT_KEY` | Midtrans client key               | —             |
| `R2_ACCOUNT_ID`       | Cloudflare R2 account id          | —             |
| `R2_ACCESS_KEY_ID`    | R2 access key                     | —             |
| `R2_SECRET_ACCESS_KEY`| R2 secret key                     | —             |
| `R2_BUCKET`           | R2 bucket name                    | —             |
| `R2_PUBLIC_URL`       | Public base URL for uploaded files| —             |

### File storage

Avatar and campaign images are uploaded to **Cloudflare R2** (S3-compatible)
through the AWS SDK. When the `R2_*` variables are empty the app falls back to
local disk under `./images` so development works without cloud credentials.

## API

Interactive Swagger UI is available at `http://localhost:8080/swagger/index.html`
once the server is running. Regenerate the spec after changing annotations with
`make docs`.

Base path: `/api/v1`. Protected routes need an `Authorization: Bearer <token>`
header.

| Method | Endpoint                       | Auth | Description                       |
| ------ | ------------------------------ | :--: | --------------------------------- |
| GET    | `/health`                      |      | Health check                      |
| POST   | `/users`                       |      | Register a user                   |
| POST   | `/sessions`                    |      | Login                             |
| POST   | `/email_checkers`              |      | Check email availability          |
| POST   | `/avatars`                     |  ✓   | Upload avatar                     |
| GET    | `/users/fetch`                 |  ✓   | Fetch current user                |
| GET    | `/campaigns`                   |      | List campaigns (`?user_id=` opt.) |
| GET    | `/campaigns/:id`               |      | Campaign detail                   |
| POST   | `/campaigns`                   |  ✓   | Create campaign                   |
| PUT    | `/campaigns/:id`               |  ✓   | Update campaign                   |
| POST   | `/campaign-images`             |  ✓   | Upload campaign image             |
| GET    | `/campaigns/:id/transactions`  |  ✓   | Transactions of a campaign        |
| GET    | `/transactions`                |  ✓   | Current user's transactions       |
| POST   | `/transactions`                |  ✓   | Create a donation                 |
| POST   | `/transactions/notification`   |      | Midtrans payment webhook          |

## Testing

```bash
make test     # go test -race ./...
make cover    # coverage report per function
```

## Tech Stack

Go · Gin · GORM · PostgreSQL · JWT · google/wire · Swagger (swaggo) · Cloudflare R2 (S3) · Midtrans · Docker

## License

[MIT](LICENSE) © Wiku Karno
