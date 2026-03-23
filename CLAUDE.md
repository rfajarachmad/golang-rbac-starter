# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make run              # Start server (go run ./cmd/web/main.go)
make dev              # Live reload via air (watches .go and .json files)
make build            # Compile to bin/go-rbac-starter
make test             # Run all tests: go test ./... -v
make test-cover       # Tests with HTML coverage report
make test-cover-check # Enforce coverage thresholds (overall≥70%, usecase≥60%, delivery≥70%)
make test-race        # Tests with race detector
make lint             # golangci-lint run ./...
make check            # fmt + vet + tidy + lint + test + coverage enforcement
```

Run a single test:
```bash
go test -v -count=1 -run TestCreateContact ./test/
```

Database migrations (PostgreSQL):
```bash
make migrate-up       # Apply all migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration (prompts for name)
```

## Architecture

This is a Go Clean Architecture project following layered dependency flow:

**Entity → Repository → UseCase → Controller → Route**

All application code lives under `internal/`. Dependency injection is wired manually in `internal/config/app.go` via `Bootstrap()`.

### Layer responsibilities

- **`internal/entity/`** — GORM database models (User, Contact, Address). Define table names and relationships.
- **`internal/repository/`** — Data access. Generic `Repository[T]` base provides CRUD; domain repos add specific queries (e.g., `FindByEmail`, `Search` with pagination).
- **`internal/usecase/`** — Business logic. Each method opens a transaction (`tx.Begin()`/`tx.Commit()`), validates input, calls repository, returns model responses. Returns `fiber.ErrNotFound`/`fiber.ErrConflict` etc. for HTTP error mapping.
- **`internal/delivery/http/`** — Fiber HTTP handlers. Parse request, extract auth from `middleware.GetUser(ctx)`, delegate to usecase, return `WebResponse[T]`.
- **`internal/delivery/http/route/`** — Route registration split into `SetupGuestRoute()` (register, login) and `SetupAuthRoute()` (all protected endpoints).
- **`internal/delivery/http/middleware/`** — Auth middleware validates `Authorization` header token via `UserUseCase.Verify()`, stores `model.Auth` in `ctx.Locals("auth")`.
- **`internal/model/`** — Request/response DTOs with validation tags. `converter/` subpackage maps entity → response.
- **`internal/config/`** — Infrastructure initialization (Fiber, GORM/Postgres, Logrus, Viper, Validator) and the `Bootstrap()` DI wiring function.

### Key patterns

- **Generic repository**: `Repository[T]` in `repository/repository.go` provides `Create`, `Update`, `Delete`, `FindById`, `FindAll`, `CountById` for any entity type.
- **Transaction-per-request**: Every usecase method wraps operations in `tx := c.DB.WithContext(ctx).Begin()` with `defer tx.Rollback()`.
- **Token auth**: Users login with email/password → receive UUID token → pass token in `Authorization` header. No JWT.
- **Scoped data**: Contacts belong to a user, addresses belong to a contact. All queries enforce ownership via `user_id`/`contact_id`.
- **Config via `config.json`**: Viper loads from `config.json` (searches `./`, `../`, `../../`). Database, web port, log level all configured here.

## Testing

Tests are **integration tests** in the `test/` package. They bootstrap the full Fiber app with a real PostgreSQL database via `test/init.go`'s `init()` function, then use `app.Test(httptest.NewRequest(...))` to hit endpoints.

- `test/init.go` — Non-test file with `init()` that wires up the full app (same as production bootstrap).
- `test/helper_test.go` — `ClearAll()`, `SeedUser()`, `SeedContact()`, `SeedAddress()` helpers for test data setup.
- `test/user_test.go`, `test/contact_test.go`, `test/address_test.go` — Test files per domain.

Tests require a running PostgreSQL instance with migrations applied. The database connection is configured in `config.json`.

## API Endpoints

OpenAPI 3.0.3 spec is at `api/api-spec.json`.

- **Guest**: `POST /api/users` (register), `POST /api/users/_login` (login)
- **Auth required**: User profile (`GET/PATCH /api/users/_current`, `DELETE /api/users`), Contacts CRUD (`/api/contacts`), Addresses CRUD (`/api/contacts/:contactId/addresses`)

## Database

PostgreSQL with GORM. Migrations in `db/migrations/` using golang-migrate (sequential numbering: `000001_`, `000002_`, etc.).

Three tables: `users` (serial PK), `contacts` (UUID PK, FK to users), `addresses` (UUID PK, FK to contacts).
