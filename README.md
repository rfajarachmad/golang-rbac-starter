# Go RBAC Starter

A RESTful API built with Go following Clean Architecture principles. Features role-based access control (RBAC) with token-based auth, contact management, and address management with scoped data ownership.

## Tech Stack

- **Go 1.26** with [Fiber](https://github.com/gofiber/fiber) HTTP framework
- **PostgreSQL** with [GORM](https://gorm.io) ORM
- **golang-migrate** for database migrations
- **golangci-lint** with 18 linters enabled
- **Viper** for configuration, **Logrus** for logging

## Architecture

Clean Architecture with layered dependency flow:

```
Entity → Repository → UseCase → Controller → Route
```

| Layer | Path | Responsibility |
|-------|------|----------------|
| Entity | `internal/entity/` | GORM database models |
| Repository | `internal/repository/` | Data access with generic `Repository[T]` base |
| UseCase | `internal/usecase/` | Business logic, transactions, validation |
| Controller | `internal/delivery/http/` | HTTP handlers (Fiber) |
| Route | `internal/delivery/http/route/` | Route registration (guest & auth) |
| Middleware | `internal/delivery/http/middleware/` | Auth + RBAC permission middleware |
| Model | `internal/model/` | Request/response DTOs with validation tags |
| Config | `internal/config/` | Infrastructure init & DI wiring via `Bootstrap()` |

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [air](https://github.com/air-verse/air) (optional, for live reload)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/rfajarachmad/golang-rbac-starter.git
   cd golang-rbac-starter
   ```

2. Configure the database connection in `config.json`:
   ```json
   {
     "database": {
       "host": "localhost",
       "port": 5432,
       "username": "your_user",
       "password": "your_password",
       "name": "your_db",
       "sslmode": "disable"
     }
   }
   ```

3. Run database migrations:
   ```bash
   make migrate-up
   ```

4. Start the server:
   ```bash
   make run
   ```

   Or with live reload:
   ```bash
   make dev
   ```

The server starts on `http://localhost:8080`.

A default admin user is seeded by the migration:
- Email: `admin@example.com`
- Password: `admin123`

## Authorization (RBAC)

The system enforces two layers of access control:

1. **Permission-based access** — Each route is protected by `RequirePermission()` middleware that checks the user's role permissions
2. **Data ownership scoping** — Users can only access their own contacts and addresses (enforced at the repository level)

### Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| `admin` | Full system access | All 16 permissions including admin endpoints |
| `user` | Standard user | CRUD on own data (11 permissions) |
| `viewer` | Read-only access | Read own profile, contacts, addresses (3 permissions) |

New registrations default to the `user` role. Admins can assign roles via `PATCH /api/admin/users/:userId/role`.

## API Endpoints

Full OpenAPI 3.0.3 spec available at `api/api-spec.json`.

### Guest Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Register new user |
| POST | `/api/users/_login` | Login (returns token) |

### Authenticated Routes (user + admin)

All require `Authorization: <token>` header. Permission required shown in parentheses.

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/users/_current` | Get current user profile | `user:read` |
| PATCH | `/api/users/_current` | Update current user | `user:update` |
| DELETE | `/api/users` | Logout | `user:delete` |
| POST | `/api/contacts` | Create contact | `contact:create` |
| GET | `/api/contacts/:contactId` | Get contact | `contact:read` |
| PUT | `/api/contacts/:contactId` | Update contact | `contact:update` |
| DELETE | `/api/contacts/:contactId` | Delete contact | `contact:delete` |
| GET | `/api/contacts` | Search contacts (with pagination) | `contact:read` |
| POST | `/api/contacts/:contactId/addresses` | Create address | `address:create` |
| GET | `/api/contacts/:contactId/addresses` | List addresses | `address:read` |
| GET | `/api/contacts/:contactId/addresses/:addressId` | Get address | `address:read` |
| PUT | `/api/contacts/:contactId/addresses/:addressId` | Update address | `address:update` |
| DELETE | `/api/contacts/:contactId/addresses/:addressId` | Delete address | `address:delete` |

### Admin Routes

Require `admin` role.

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/admin/users` | List all users (paginated) | `admin:user:list` |
| GET | `/api/admin/users/:userId` | Get any user | `admin:user:read` |
| PATCH | `/api/admin/users/:userId/role` | Assign role to user | `admin:user:update` |
| DELETE | `/api/admin/users/:userId` | Delete any user | `admin:user:delete` |
| GET | `/api/admin/roles` | List all roles | `admin:role:manage` |
| GET | `/api/admin/roles/:roleId` | Get role with permissions | `admin:role:manage` |

## Database

Six tables managed via sequential migrations in `db/migrations/`:

- **users** — serial PK, email/password/token auth, FK to roles
- **roles** — serial PK, name (admin/user/viewer)
- **permissions** — serial PK, name (e.g., `contact:create`)
- **role_permissions** — join table (role_id, permission_id)
- **contacts** — UUID PK, belongs to user
- **addresses** — UUID PK, belongs to contact

## Development

```bash
make help             # Show all available commands
make run              # Start server
make dev              # Live reload via air
make build            # Compile to bin/go-rbac-starter
make test             # Run all tests
make test-cover       # Tests with HTML coverage report
make test-cover-check # Enforce coverage thresholds
make test-race        # Tests with race detector
make lint             # Run golangci-lint (18 linters)
make check            # fmt + vet + tidy + lint + test + coverage
make migrate-up       # Apply all migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration
```

## Testing

Integration tests in `test/` bootstrap the full Fiber app against a real PostgreSQL database. Tests require a running PostgreSQL instance with migrations applied.

```bash
# Run all tests
make test

# Run a single test
go test -v -count=1 -run TestCreateContact ./test/

# Coverage thresholds: overall >= 70%, usecase >= 60%, delivery >= 70%
make test-cover-check
```

## License

MIT
