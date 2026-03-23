# ── config ────────────────────────────────────────────────────────────────
APP_NAME  = go-rbac-starter
MAIN      = ./cmd/web/main.go
DB_URL    = postgres://postgres:postgres@localhost:5433/go?sslmode=disable

# ── run ───────────────────────────────────────────────────────────────────
.PHONY: run
run:
	go run $(MAIN)

## run with live reload (requires: go install github.com/air-verse/air@latest)
.PHONY: dev
dev:
	air

# ── build ─────────────────────────────────────────────────────────────────
.PHONY: build
build:
	go build -o bin/$(APP_NAME) $(MAIN)

.PHONY: clean
clean:
	rm -rf bin/

# ── test ──────────────────────────────────────────────────────────────────
.PHONY: test
test:
	go test ./... -v

.PHONY: test-cover
test-cover:
	go test -count=1 -coverprofile=coverage.out -coverpkg=./internal/... ./test/
	go tool cover -html=coverage.out

.PHONY: test-cover-check
test-cover-check:
	@./scripts/check-coverage.sh

.PHONY: test-race
test-race:
	go test ./... -race

# ── code quality ──────────────────────────────────────────────────────────
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run ./...

# ── database ──────────────────────────────────────────────────────────────
.PHONY: migrate-up
migrate-up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path db/migrations \
		-database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path db/migrations \
		-database "$(DB_URL)" down 1

.PHONY: migrate-drop
migrate-drop:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path db/migrations \
		-database "$(DB_URL)" drop -f

.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		create -ext sql -dir db/migrations -seq $$name

# ── all-in-one ────────────────────────────────────────────────────────────
.PHONY: check
check: fmt vet tidy lint test test-cover-check

.PHONY: help
help:
	@echo ""
	@echo "  make run            start the server"
	@echo "  make dev            start with live reload (requires air)"
	@echo "  make build          compile to bin/$(APP_NAME)"
	@echo "  make clean          remove bin/"
	@echo ""
	@echo "  make test           run all tests"
	@echo "  make test-cover       run tests + open coverage"
	@echo "  make test-cover-check enforce coverage thresholds"
	@echo "  make test-race        run tests with race detector"
	@echo ""
	@echo "  make fmt            format code"
	@echo "  make vet            run go vet"
	@echo "  make tidy           tidy modules"
	@echo "  make lint           run golangci-lint"
	@echo ""
	@echo "  make migrate-up     apply migrations"
	@echo "  make migrate-down   rollback last migration"
	@echo "  make migrate-drop   drop everything"
	@echo "  make migrate-create create a new migration"
	@echo ""
	@echo "  make check            fmt + vet + tidy + lint + test + coverage"
	@echo ""
