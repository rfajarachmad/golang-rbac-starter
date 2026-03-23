---
name: add-domain
description: Scaffold a new domain entity with full clean architecture layers (entity, model, converter, repository, usecase, controller, routes, migration, tests)
user_invocable: true
---

# Add Domain Entity Skill

You are scaffolding a new domain entity for the go-rbac-starter project. This project follows Go Clean Architecture with layers: Entity → Repository → UseCase → Controller → Route.

## Input

The user will provide:
- **Entity name** (e.g., "Product", "Order", "Note")
- **Fields** with types (e.g., "title string, price int, description string")
- **Ownership**: whether it belongs to a user (scoped by `user_id`) or another entity
- **ID type**: UUID string (default) or auto-increment int

If the user doesn't specify fields, ask before proceeding.

## Files to Generate

For an entity named `{Entity}` (e.g., Product) with table `{entities}` (e.g., products):

### 1. Migration: `db/migrations/{next_seq}_create_{entities}_table.up.sql`

Find the highest existing migration sequence number in `db/migrations/` and increment by 1. Use zero-padded 6-digit format (e.g., `000004`).

```sql
CREATE TABLE IF NOT EXISTS {entities} (
    id         VARCHAR(100) PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users(id),
    -- domain fields here --
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

Also create the matching `.down.sql`:
```sql
DROP TABLE IF EXISTS {entities};
```

### 2. Entity: `internal/entity/{entity}_entity.go`

```go
package entity

import "time"

type {Entity} struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	// domain fields with gorm tags
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreignKey:user_id;references:id"`
}

func (e *{Entity}) TableName() string {
	return "{entities}"
}
```

### 3. Model: `internal/model/{entity}_model.go`

Define these structs with `json` and `validate` tags:
- `{Entity}Response` — all fields with `json:",omitempty"`, includes CreatedAt/UpdatedAt
- `Create{Entity}Request` — UserId as `json:"-" validate:"required"`, domain fields with validation
- `Update{Entity}Request` — ID and UserId as `json:"-"`, domain fields with validation
- `Get{Entity}Request` — ID and UserId as `json:"-"`
- `Delete{Entity}Request` — ID and UserId as `json:"-"`
- `Search{Entity}Request` — UserId, filter fields, Page (min=1), Size (min=1,max=100)

### 4. Converter: `internal/model/converter/{entity}_converter.go`

```go
func {Entity}ToResponse(e *entity.{Entity}) *model.{Entity}Response {
	// map all fields from entity to response
}
```

### 5. Repository: `internal/repository/{entity}_repository.go`

```go
type {Entity}Repository struct {
	Repository[entity.{Entity}]
	Log *logrus.Logger
}

func New{Entity}Repository(log *logrus.Logger) *{Entity}Repository {
	return &{Entity}Repository{Log: log}
}

func (r *{Entity}Repository) FindByIdAndUserId(db *gorm.DB, e *entity.{Entity}, id string, userId int) error {
	return db.Where("id = ? AND user_id = ?", id, userId).Take(e).Error
}
```

Add a `Search` method with pagination and `ILIKE` filtering if SearchRequest has filter fields.

### 6. UseCase: `internal/usecase/{entity}_usecase.go`

Implement these methods following the transaction-per-request pattern:
- `Create` — validate, generate UUID, create, commit
- `Update` — validate, find by ID+UserId, update fields, commit
- `Get` — validate, find by ID+UserId, commit
- `Delete` — validate, find by ID+UserId, delete, commit
- `Search` — validate, call repository Search, build PageMetadata, commit

Every method must:
- Start with `tx := c.DB.WithContext(ctx).Begin()` and `defer tx.Rollback()`
- Validate with `c.Validate.Struct(request)`
- Log warnings with `c.Log.Warnf()`
- Return `fiber.ErrBadRequest`, `fiber.ErrNotFound`, `fiber.ErrInternalServerError` for errors
- Pre-allocate response slices: `make([]model.{Entity}Response, 0, len(items))`
- Use index-based range: `for i := range items`

### 7. Controller: `internal/delivery/http/{entity}_controller.go`

```go
type {Entity}Controller struct {
	Log     *logrus.Logger
	UseCase *usecase.{Entity}UseCase
}
```

Implement handlers:
- `Create` — BodyParser, set UserId from `middleware.GetUser(ctx).ID`
- `List` — build SearchRequest from query params (page default 1, size default 10)
- `Get` — read `:entityId` param, set UserId from auth
- `Update` — BodyParser, set ID from param, UserId from auth
- `Delete` — read `:entityId` param, set UserId from auth

All handlers return `model.WebResponse[T]`. List returns with `Paging` field.

### 8. Wire into existing files

**`internal/delivery/http/route/route.go`**:
- Add `{Entity}Controller *http.{Entity}Controller` to `RouteConfig`
- Add CRUD routes in `SetupAuthRoute()`:
```go
c.App.Get("/api/{entities}", c.{Entity}Controller.List)
c.App.Post("/api/{entities}", c.{Entity}Controller.Create)
c.App.Put("/api/{entities}/:{entityId}", c.{Entity}Controller.Update)
c.App.Get("/api/{entities}/:{entityId}", c.{Entity}Controller.Get)
c.App.Delete("/api/{entities}/:{entityId}", c.{Entity}Controller.Delete)
```

**`internal/config/app.go`**:
- Add repository: `{entity}Repository := repository.New{Entity}Repository(config.Log)`
- Add usecase: `{entity}UseCase := usecase.New{Entity}UseCase(config.DB, config.Log, config.Validate, {entity}Repository)`
- Add controller: `{entity}Controller := http.New{Entity}Controller({entity}UseCase, config.Log)`
- Add to RouteConfig: `{Entity}Controller: {entity}Controller,`

### 9. Tests

**`test/helper_test.go`**:
- Add `Clear{Entities}()` function — `db.Where("id is not null").Delete(&entity.{Entity}{})`
- Add it to `ClearAll()` — clear child tables before parent tables (FK order)
- Add `Seed{Entity}(t *testing.T, userId int) *entity.{Entity}` function

**`test/{entity}_test.go`**:
Write integration tests using `app.Test(httptest.NewRequest(...))`:
- `TestCreate{Entity}` — create with valid data, assert 200 + response fields
- `TestCreate{Entity}Unauthorized` — wrong token, assert 401
- `TestGet{Entity}` — seed then get, assert fields match
- `TestGet{Entity}NotFound` — assert 404
- `TestUpdate{Entity}` — seed then update, assert new values
- `TestUpdate{Entity}NotFound` — assert 404
- `TestDelete{Entity}` — seed then delete, assert true
- `TestDelete{Entity}NotFound` — assert 404
- `TestSearch{Entity}` — seed multiple, assert list length + paging
- `TestSearch{Entity}WithPagination` — seed 20, page=2 size=5, assert paging metadata

Every test: `ClearAll()` first, seed a user with `SeedUser(t, ...)`, use `user.Token` in Authorization header.

## After Scaffolding

1. Run `go mod tidy`
2. Run `golangci-lint run ./...` and fix any issues
3. Apply migration: `migrate -path db/migrations -database "$DB_URL" up`
4. Run tests: `go test -v -count=1 ./test/ | cat`
5. Report results to user
