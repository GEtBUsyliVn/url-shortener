# ===== config =====
GOOSE=goose
MIGRATIONS_DIR=./migrations
DB_DRIVER=postgres
DB_DSN=postgres://postgres:postgres@localhost:5432/url?sslmode=disable
# ==================

export GOOSE_DRIVER=$(DB_DRIVER)
export GOOSE_DBSTRING=$(DB_DSN)

up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) up

down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) down

status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) status

create:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

run:
	go run cmd/booking/main.go
