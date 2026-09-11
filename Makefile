.PHONY: help build run dev lint format \
        new-migration migrate-up migrate-down migrate-status migrate-reset \
        migrate-create migrate-version \
        docker-up docker-down test

help:
	@echo "Makefile commands:"
	@echo "  build             Build the Go application"
	@echo "  run               Run the Go application"
	@echo "  dev               Run with hot reload"
	@echo "  lint              Run golangci-lint"
	@echo "  format            Format code + imports"
	@echo "  new-migration     Create a goose migration. Usage: make new-migration name=create_users_table"
	@echo "  migrate-up        Apply all pending migrations"
	@echo "  migrate-down      Roll back the last migration"
	@echo "  migrate-status    Show migration status"
	@echo "  migrate-version   Show current DB version"
	@echo "  migrate-reset     Roll back ALL migrations"
	@echo "  docker-up         Start Docker containers"
	@echo "  docker-down       Stop Docker containers"

DB_URL        := postgres://abubakar-sadiq:mypassword@localhost:5432/ecommerce_shop?sslmode=disable
MIGRATION_DIR := db/migrations
GOOSE         := goose -env=none -dir $(MIGRATION_DIR) postgres "$(DB_URL)"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

dev:
	air   # or: go run ./cmd/api

lint:
	golangci-lint run ./...

format:
	gofmt -s -w .
	goimports -w .

# ---------- Goose migrations ----------
new-migration:
	$(GOOSE) -s create $(name) sql

migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

migrate-version:
	$(GOOSE) version


docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down

test:
	@echo "Running tests with ARGS=$(ARGS)"
	go test ./... $(ARGS)